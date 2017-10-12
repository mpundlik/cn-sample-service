# Topology example plugin

**Notice**
_This example working with simple topology data. This is not meant to be use in production code. 
The topology data is simplified and to the fully working topology need to be changed.
This task were good practice for reader. Write TAPs, Interfaces and/or Bridges separately
into DB and change the final topology plugin example to react for the changes as well_

This example and this document is about to show you step by step how to create a simple plugin using ligato framework. 
At the final stage the plugin should be able create a simple topology written in [ETCD](https://github.com/ligato/cn-infra/tree/master/db/keyval/etcdv3) 
used with [protocol buffers](https://github.com/golang/protobuf)

After reading documentation of [cn-infra](https://github.com/ligato/cn-infra/blob/master/README.md) 
the reader should be able to create a simple plugin lets say we name it a topology.

Every step is in separate go file so you don't need to copy paste but just run example.

#### Step 1
 
Create a simple plugin is pretty straightforward
This is inspired by [HelloWorld](https://github.com/ligato/cn-sample-service/tree/master/cmd/helloworld) example.
Important is to look at the tree methods Init, AfterInit and Close. Init is called when the plugin started. 
This is the plugin main initialization phase and it is mandatory. 
AfterInit is called after plugin initialized and this part of the plugin is optional. 
To clean up everything can be used Close method.

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

Second thing is how to run plugin itself. It is pretty simply.
[Flavors](https://github.com/ligato/cn-infra/tree/master/flavors), instantiate and run the plugin.

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

Now is time to make small changes. 

1. Need to change [flavor](https://github.com/ligato/cn-infra/tree/master/flavors) from local to connection.
As the plugin is intended to use [ETCD](https://github.com/ligato/cn-infra/tree/master/db/keyval/etcdv3) 
and [protocol buffers](https://github.com/golang/protobuf)

```go
	flavor := connectors.AllConnectorsFlavor{}	
```

Also it is about a time to add some data structure. The data structure is in model directory not in step2 because it will be used by another steps as well.

```proto
syntax = "proto3";

package model;

message IPAdress {
    string ip = 1;
}

message Interface {
    string name = 1;
}

message Tap {
    string name = 1;
    string mac = 2;
    repeated IPAdress ip_adresses = 3;
}

message Bridge {
    string name = 1;
    repeated Interface interfaces = 2;
}

message Topology {

    Bridge bridge = 1;
    repeated Tap taps = 2;

}
```

After creating structure and generating go file with [(see tutorial)](https://developers.google.com/protocol-buffers/docs/gotutorial)

```bash
protoc --go_out=. *.proto
```

the file topology.pb.go should be created. Now some functions needs to be added. 
- Function for creating sample data.
- Function to get ETCD path _key to store data_
- Function put() to store data into ETDC

```go
func (plugin *TopologyPlugin) buildData() *model.Topology {
	var ip1 = model.IPAdress{Ip:"127.0.0.1"}
	var ip2 = model.IPAdress{Ip:"127.0.0.2"}
	var ip3 = model.IPAdress{Ip:"127.0.0.3"}
	var ip4 = model.IPAdress{Ip:"127.0.0.4"}
	var if1 = model.Interface{Name:"interface1"}
	var if2 = model.Interface{Name:"interface2"}
	var tap1 = model.Tap{Name:"tap1", Mac:"00:00:00:00:00:00", IpAdresses:[]*model.IPAdress{&ip1, &ip2}}
	var tap2 = model.Tap{Name:"tap2", Mac:"00:00:ff:00:00:00", IpAdresses:[]*model.IPAdress{&ip4, &ip3}}
	return &model.Topology{
		Bridge:&model.Bridge{Name:"bridge1", Interfaces:[]*model.Interface{&if1, &if2}},
		Taps:[]*model.Tap{&tap1, &tap2},
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

Also the Init() function now initializes ETDC with etcd.conf file and AfterInit() stores data. 
Close() is closing decorator.

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

To be able test it as write sample data ETCD need to be running. 

```bash
sudo docker run -p 2379:2379 --name etcd --rm \
    quay.io/coreos/etcd:latest /usr/local/bin/etcd \
    -advertise-client-urls http://0.0.0.0:2379 \
    -listen-client-urls http://0.0.0.0:2379
```

When the step2 is running it should be clearly visible in logs.

#### Step 3

This step add a watcher to react for ETDC changes to be able work if data are changed (or deleted).
Also creating dependencies and if the plugin doesn't need all connectors flavors will be changed to only needed plugins.

- First of all a Deps _like dependencies_ need to be created

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

All needed dependencies are injected through deps file and initialization of ETCD can be romoved from Init() function.

- Now it needed a watcher to be registered to obtain any changes.

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

- And creating consumer who is react for all the changes

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

How to test it locally:
- Build it direct in ste3 directory
```bash
go build step3.go deps.go
```
- Run etcd docker (see in this document [Step2](#step-2))
- Run step3 with configuration
```bash
./step3 -etcdv3-config="etcd.conf"
```

Now it will be seen in logs twice that plugin saves data and receives event with change 
type Put and first with update:false second with update:true




  
