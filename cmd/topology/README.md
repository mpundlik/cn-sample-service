# Topology example plugin

**Notice**
_This example works with simple topology data. The data is not meant to be use in production code. 
The topology data is simplified. It needs to be changed to become fully working topology.
This task is a good practice for the reader. Write TAPs, Interfaces and/or Bridges separately
into DB and change the final topology plugin example to reflect the changes as well._

The example and the document provide step by step instruction on how to create a simple plugin using ligato framework. At the final stage, the plugin should be able to create this simple topology and write it into [ETCD](https://github.com/ligato/cn-infra/tree/master/db/keyval/etcdv3), using [protocol buffers](https://github.com/golang/protobuf).

After reading [cn-infra documentation](https://github.com/ligato/cn-infra/blob/master/README.md), the reader should be able to create a simple plugin. For the example purpose, the plugin's name is topology.

Every step is a separate go file. Thus the reader does not need to copy-paste, only to run the example.

#### Step 1
 
To create a simple plugin is straightforward.
This is inspired by [HelloWorld](https://github.com/ligato/cn-sample-service/tree/master/cmd/helloworld) example.
It is important to look at three methods: Init, AfterInit and Close.
- Init is called when the plugin is started. This is the plugin main initialization phase and it is mandatory. 
- AfterInit is called after plugin is initialized and this part of the plugin is optional. 
- Close method cleans channels and routines and it is mandatory in case channels and routines are used. Otherwise it does not need to be implemented.

```go
func (plugin *TopologyPlugin) Init() (err error) {
	logroot.StandardLogger().Info("Topology Plugin initialized.")
	return nil
}

// AfterInit logs a sample message.
func (plugin *TopologyPlugin) AfterInit() error {
	logroot.StandardLogger().Info("Topology Plugin is running.")
	return nil
}

// Close is called to cleanup the plugin resources.
func (plugin *TopologyPlugin) Close() error {
	return nil
}
```

Another important issue is running the plugin itself. [Flavors](https://github.com/ligato/cn-infra/tree/master/flavors) instantiate and run the plugin.

```go
func main() {
	// Some flavors from cn-infra
	flavor := local.FlavorLocal{}

	// create an instance of the plugin
	topoPlugin := TopologyPlugin{}

	// Create new agent
	agent := core.NewAgent(logroot.StandardLogger(), 15*time.Second, append(flavor.Plugins(), &core.NamedPlugin{PluginName: PluginID, Plugin: &topoPlugin})...)

	core.EventLoopWithInterrupt(agent, nil)
}
```

#### Step 2

In this step we focus on making few small changes. 

Change [flavor](https://github.com/ligato/cn-infra/tree/master/flavors) from local to connection, as the plugin is intended to use [ETCD](https://github.com/ligato/cn-infra/tree/master/db/keyval/etcdv3) and [protocol buffers](https://github.com/golang/protobuf).

```go
	flavor := connectors.AllConnectorsFlavor{}	
```

Add some data structure. The data structure is located in model directory not in step2, because it will be used by another steps as well.

```proto
syntax = "proto3";

package model;

message Topology {
    string name = 1;

    message Interface {
        string name = 1;

        message Tap {
            string mac = 1;

            message IPAdress {
                string ip = 1;
            }
            
            repeated IPAdress ip_adresses = 2;
        }
        
        Tap tap = 2;
    }
    
    repeated Interface interfaces = 2;
}
```

After creating structure and [generating go file](https://developers.google.com/protocol-buffers/docs/gotutorial),

```bash
protoc --go_out=. *.proto
```

the file topology.pb.go should be created. Now, add some functions: 
- Function for creating sample data.
- Function to get ETCD path _key to store data_
- Function put() to save data into ETDC

```go
func (plugin *TopologyPlugin) buildData() *model.Topology {
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

func (plugin *TopologyPlugin) put() {

	key := etcdPath()
	topo := plugin.buildData();

	logroot.StandardLogger().Info("Saving: ", key)
	logroot.StandardLogger().Info("Data: ", topo)

	// Insert the key-value pair.
	plugin.protoDB.Put(key, topo)

}

func etcdPath() string {
	return "/topology"
}
```

The Init() function initializes ETDC with etcd.conf file and AfterInit() stores the data. 
Close() is a closing decorator.

etcd.conf:
```bash
insecure-transport: true
endpoints:
  - "172.17.0.1:2379"
```

```go
func (plugin *TopologyPlugin) Init() (err error) {

	//Lets configure etcd
	fileConfig := &etcdv3.Config{}
	err = config.ParseConfigFromYamlFile("etcd.conf", fileConfig)
	.
	.
	cfg, err := etcdv3.ConfigToClientv3(fileConfig)
	.
	.
	db, err := etcdv3.NewEtcdConnectionWithBytes(*cfg, logroot.StandardLogger())
	
	// Initialize proto decorator.
	plugin.protoDB = kvproto.NewProtoWrapper(db)

	logroot.StandardLogger().Info("Topology plugin initialized properly.")
	return nil
}

// AfterInit logs a sample message.
func (plugin *TopologyPlugin) AfterInit() error {
	// Write some data.
	plugin.put()
	return nil
}

// Close is called to cleanup the plugin resources.
func (plugin *TopologyPlugin) Close() error {
	// Close proto decorator.
	plugin.protoDB.Close()
	return nil
}
```

ETCD needs to be running in order to test the plugin and write the data to ETCD. 

```bash
sudo docker run -p 2379:2379 --name etcd --rm \
    quay.io/coreos/etcd:latest /usr/local/bin/etcd \
    -advertise-client-urls http://0.0.0.0:2379 \
    -listen-client-urls http://0.0.0.0:2379
```

The logs content should reflect that step2 is running.

#### Step 3

Step3 adds a watcher into the plugin, that watches for changes in ETDC and reacts to them. 
In previous steps, predefined flavours from cn-infra were used. In this step, flavor specific to the plugin was created. 

- First of all, _Deps structure_ needs to be created.

```go
// Deps lists dependencies of TopologyPlugin.
type Deps struct {
	local.PluginInfraDeps                 // injected
	Publisher datasync.KeyProtoValWriter  // injected - To write ETCD data
	Watcher   datasync.KeyValProtoWatcher // injected - To watch ETCD data
}

// TopologyFlavor is a set of plugins required for the topology example. It could be used AllConnectionFlavors
// but don't need to pull all the dependencies plugin doesn't need.
type TopologyFlavor struct {
	// Local flavor to access the Infra (logger, service label, status check)
	*local.FlavorLocal
	// Resync orchestrator
	ResyncOrch resync.Plugin
	// Etcd plugin
	ETCD etcdv3.Plugin
	// Etcd sync which manages and injects connection
	ETCDDataSync kvdbsync.Plugin
	// Topology plugin
	TopologyExample TopologyPlugin
	// Use channel when the topology plugin is finished
	closeChan *chan struct{}
}
```

All necessary dependencies are injected through _Deps struct_.

- Now a watcher needs to be registered to receive any changes.

```go
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
```

- And consumer that reacts to all the changes needs to be created.

```go
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
```

To test the plugin locally:
- Build the plugin directly in step3 directory.
```bash
go build step3.go deps.go
```
- Run etcd docker (see in this document [Step2](#step-2)).
- Run step3 with configuration.
```bash
./step3 -etcdv3-config="etcd.conf"
```

After running the plugin, two events associated with writing into ETCD are present in the logs and two events associated with changes in ETCD are presnet in the logs.




  