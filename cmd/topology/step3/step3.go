// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"github.com/ligato/cn-infra/core"
	"github.com/ligato/cn-infra/datasync"
	"github.com/ligato/cn-infra/examples/model"
	"github.com/ligato/cn-infra/utils/safeclose"
	"github.com/ligato/cn-sample-service/cmd/topology/model"
	"strings"
	"time"
	"github.com/ligato/cn-infra/flavors/local"
	"github.com/ligato/cn-infra/flavors/connectors"
)

type TopologyPlugin struct {
	Deps

	changeChannel chan datasync.ChangeEvent  // Channel used by the watcher for change events.
	resyncChannel chan datasync.ResyncEvent  // Channel used by the watcher for resync events.
	context       context.Context            // Used to cancel watching.
	watchDataReg  datasync.WatchRegistration // To subscribe on data change/resync events.
	closeChannel  *chan struct{}
}

func main() {
	// Init close channel used to stop the example.
	exampleFinished := make(chan struct{}, 1)

	// Start Agent with ExampleFlavor
	// (combination of ExamplePlugin & cn-infra plugins).

	tf := TopologyFlavor{closeChan: &exampleFinished}

	if tf.FlavorLocal == nil {
		tf.FlavorLocal = &local.FlavorLocal{}
	}

	tf.ResyncOrch.Deps.PluginLogDeps = *tf.FlavorLocal.LogDeps("resync-orch")
	tf.ETCD.Deps.PluginInfraDeps = *tf.InfraDeps("etcdv3")
	connectors.InjectKVDBSync(&tf.ETCDDataSync, &tf.ETCD, tf.ETCD.PluginName, tf.FlavorLocal, &tf.ResyncOrch)

	// Inject infra + transport (publisher, watcher) to example plugin
	tf.TopologyExample.PluginInfraDeps = *tf.FlavorLocal.InfraDeps("topology-plugin")
	tf.TopologyExample.Publisher = &tf.ETCDDataSync
	tf.TopologyExample.Watcher = &tf.ETCDDataSync
	tf.TopologyExample.closeChannel = tf.closeChan

	agent := core.NewAgent(core.Inject(&tf))
	core.EventLoopWithInterrupt(agent, exampleFinished)

}

// Init is called on plugin startup. New logger is instantiated.
func (plugin *TopologyPlugin) Init() error {

	// Initialize plugin fields.
	plugin.resyncChannel = make(chan datasync.ResyncEvent)
	plugin.changeChannel = make(chan datasync.ChangeEvent)
	plugin.context = context.Background()

	// Start the consumer (ETCD watcher).
	go plugin.consumer()
	// Subscribe watcher to be able to watch on data changes and resync events.
	err := plugin.subscribeWatcher()
	if err != nil {
		return err
	}

	plugin.Log.Info("Topology plugin initialized properly.")
	return nil
}

// AfterInit logs a sample message.
func (plugin *TopologyPlugin) AfterInit() error {
	// Write some data.
	go plugin.writeData()
	go plugin.closeExample()
	return nil
}

// Close is called to cleanup the plugin resources.
func (plugin *TopologyPlugin) Close() error {
	safeclose.CloseAll(plugin.Publisher, plugin.Watcher, plugin.resyncChannel, plugin.changeChannel, plugin.closeChannel)
	return nil
}

func (plugin *TopologyPlugin) writeData() {
	// Wait for the consumer to initialize
	time.Sleep(1 * time.Second)

	key := etcdPath()
	topo := plugin.buildData(1)

	plugin.Log.Infof("Saving data: %v into: %v", topo, key)
	// Insert the key-value pair.
	plugin.Publisher.Put(key, topo)

	time.Sleep(1 * time.Second)

	topo = plugin.buildData(2)

	plugin.Log.Infof("Saving data: %v into: %v", topo, key)
	// Insert the key-value pair.
	plugin.Publisher.Put(key, topo)

}

func (plugin *TopologyPlugin) buildData(topologyId int) *model.Topology {
	var ip1 = model.Topology_Interface_Tap_IPAdress{Ip: "127.0.0.1"}
	var ip2 = model.Topology_Interface_Tap_IPAdress{Ip: "127.0.0.2"}
	var ip3 = model.Topology_Interface_Tap_IPAdress{Ip: "127.0.0.3"}
	var ip4 = model.Topology_Interface_Tap_IPAdress{Ip: "127.0.0.4"}
	var tap1 = model.Topology_Interface_Tap{Mac: "00:00:00:00:00:00", IpAdresses: []*model.Topology_Interface_Tap_IPAdress{&ip1, &ip2}}
	var tap2 = model.Topology_Interface_Tap{Mac: "00:00:ff:00:00:00", IpAdresses: []*model.Topology_Interface_Tap_IPAdress{&ip4, &ip3}}
	var if1 = model.Topology_Interface{Name: "interface1", Tap:&tap1}
	var if2 = model.Topology_Interface{Name: "interface2", Tap:&tap2}
	if topologyId == 1 {
		return &model.Topology{
			Name: "topology",
			Interfaces:[]*model.Topology_Interface{&if1},
		}
	} else {
		return &model.Topology{
			Name: "topology2",
			Interfaces:[]*model.Topology_Interface{&if2},
		}
	}
}

func (plugin *TopologyPlugin) closeExample() {

	time.Sleep(8 * time.Second)

	plugin.context.Done()
	plugin.Log.Info("topology plugin finished, sending shutdown ...")
	// Close the example
	*plugin.closeChannel <- struct{}{}
}

// consumer (watcher) is subscribed to watch on data store changes.
// Changes arrive via data change channel, get identified based on the key
// and printed into the log.
func (plugin *TopologyPlugin) consumer() {
	plugin.Log.Info("KeyValProtoWatcher started")
	for {
		select {
		// WATCH: demonstrate how to receive data change events.
		case dataChng := <-plugin.changeChannel:
			plugin.Log.Printf("Received event: %v", dataChng)
			// If event arrives, the key is extracted and used together with
			// the expected prefix to identify item.
			key := dataChng.GetKey()
			if strings.HasPrefix(key, etcdPathPrefix()) {
				var value, previousValue etcdexample.EtcdExample
				// The first return value is diff - boolean flag whether previous value exists or not
				err := dataChng.GetValue(&value)
				if err != nil {
					plugin.Log.Error(err)
				}
				diff, err := dataChng.GetPrevValue(&previousValue)
				if err != nil {
					plugin.Log.Error(err)
				}
				plugin.Log.Infof("Event arrived to etcd eventHandler, key %v, update: %v, change type: %v,",
					dataChng.GetKey(), diff, dataChng.GetChangeType())
			}
		case <-plugin.context.Done():
			plugin.Log.Warnf("Stop watching events")
		}
	}
}

// subscribeWatcher subscribes for data change and data resync events.
// Events are delivered to the consumer via the selected channels.
// ETCD watcher adapter is used to perform the registration behind the scenes.
func (plugin *TopologyPlugin) subscribeWatcher() (err error) {
	prefix := etcdPathPrefix()
	plugin.Log.Infof("Prefix: %v", prefix)
	plugin.watchDataReg, err = plugin.Watcher.
		Watch("TopologyPlugin", plugin.changeChannel, plugin.resyncChannel, prefix)
	if err != nil {
		return err
	}

	plugin.Log.Info("KeyValProtoWatcher subscribed")

	return nil
}

func etcdPath() string {
	return etcdPathPrefix() + "index"
}

func etcdPathPrefix() string {
	return "/topology/"
}
