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

	"github.com/ligato/cn-infra/config"
	"github.com/ligato/cn-infra/core"
	"github.com/ligato/cn-infra/datasync"
	"github.com/ligato/cn-infra/db/keyval/etcdv3"
	"github.com/ligato/cn-infra/db/keyval/kvproto"
	"github.com/ligato/cn-infra/flavors/connectors"
	"github.com/ligato/cn-sample-service/cmd/topology/model"
	"github.com/ligato/cn-infra/flavors/local"
)

type TopologyPlugin struct {
	local.PluginInfraDeps
	data model.Topology
	db   datasync.KeyProtoValWriter
}

func main() {

	// Create new agent using connectors flavor
	agent := connectors.NewAgent(connectors.WithPlugins(func(allConnectors *connectors.AllConnectorsFlavor) []*core.NamedPlugin {
		// create an instance of the plugin
		topoPlugin := &TopologyPlugin{}
		topoPlugin.PluginInfraDeps = *allConnectors.InfraDeps("topology-plugin")

		return []*core.NamedPlugin{
			{topoPlugin.PluginName, topoPlugin},
		}
	}))


	core.EventLoopWithInterrupt(agent, nil)
}

// PluginID of the custom topology plugin
const PluginID core.PluginName = "topology-plugin"

// Plugin key to save topology
const TopologyKey string = "/topology"

func (plugin *TopologyPlugin) initializeEtcd() error {

	//Lets configure etcd
	fileConfig := &etcdv3.Config{}
	if parseError := config.ParseConfigFromYamlFile("etcd.conf", fileConfig); parseError == nil {
		if cfg, configError := etcdv3.ConfigToClientv3(fileConfig); configError == nil {
			if db, dbError := etcdv3.NewEtcdConnectionWithBytes(*cfg, plugin.Log); dbError == nil {
				plugin.db = kvproto.NewProtoWrapper(db)
				return nil
			} else {
				plugin.Log.Error("Cannot connect ETCD.")
				return dbError
			}
		} else {
			plugin.Log.Error("Wrong ETCD configure file")
			return configError
		}
	} else {
		plugin.Log.Error("Cannot find ETCD or corrupted ETCD configure file")
		return parseError
	}

}

// Init is called on plugin startup. New logger is instantiated.
func (plugin *TopologyPlugin) Init() (err error) {

	plugin.buildData()
	if err := plugin.initializeEtcd(); err == nil {
		plugin.Log.Info("Topology plugin initialized properly.")
	}
	return err

}

// AfterInit logs a sample message.
func (plugin *TopologyPlugin) AfterInit() error {
	// Write some data.
	return plugin.put()
}

// Close is called to cleanup the plugin resources.
func (plugin *TopologyPlugin) Close() error {
	return nil
}

func (plugin *TopologyPlugin) buildData() {
	var ip1 = model.Topology_Interface_Tap_IPAdress{Ip: "127.0.0.1"}
	var ip2 = model.Topology_Interface_Tap_IPAdress{Ip: "127.0.0.2"}
	var ip3 = model.Topology_Interface_Tap_IPAdress{Ip: "127.0.0.3"}
	var ip4 = model.Topology_Interface_Tap_IPAdress{Ip: "127.0.0.4"}
	var tap1 = model.Topology_Interface_Tap{Mac: "00:00:00:00:00:00", IpAdresses: []*model.Topology_Interface_Tap_IPAdress{&ip1, &ip2}}
	var tap2 = model.Topology_Interface_Tap{Mac: "00:00:ff:00:00:00", IpAdresses: []*model.Topology_Interface_Tap_IPAdress{&ip4, &ip3}}
	var if1 = model.Topology_Interface{Name: "interface1", Tap:&tap1}
	var if2 = model.Topology_Interface{Name: "interface2", Tap:&tap2}
	plugin.data = model.Topology{
		Name: "topology",
		Interfaces:[]*model.Topology_Interface{&if1, &if2},
	}
}

func (plugin *TopologyPlugin) put() error {

	plugin.Log.Info("Saving: ", TopologyKey)
	plugin.Log.Info("Data: ", plugin.data)

	// Insert the key-value pair.
	return plugin.db.Put(TopologyKey, &plugin.data)

}
