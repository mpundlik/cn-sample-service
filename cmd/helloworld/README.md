# HelloWorld

HelloWorld as the name suggest is a simple extension of vpp-agent. It add a single plugin that creates
a custom logger an logs a message "Hello World!!!". The aim of this example is to list basic steps that
are necessary in order to add a custom plugin to vpp-agent.

1. Use an existing flavour - collection of plugins
```
	import "github.com/ligato/vpp-agent/flavours/vpp"


	f := vpp.Flavour{}
```
Alternatively, you can create a custom [flavour](https://github.com/ligato/vpp-agent/tree/master/flavours).

2. Declare a structure for your plugin. Apart from the internal fields the structure must specifify
the dependencies of the plugin.

```go
// HelloWorldPlugin is a plugin that showcase the extensibility of vpp agent.
type HelloWorldPlugin struct {
	// LogFactory is a dependency of the plugin that needs to be injected.
	// This dependency provides an API to create a logger instance
	LogFactory logging.LogFactory

    //..
}
```

3. Implement plugin lifecycle methods. `Init()` and `Close` are mandatory `AfterInit` is optional.

```
// Init is called on plugin startup. New logger is instantiated.
func (plugin *HelloWorldPlugin) Init() (err error) {
    // use injected dependencies in this case to create a logger
    plugin.Logger, err = plugin.LogFactory.NewLogger(string(PluginID))
	
	//...
	return err
}

// AfterInit logs a sample message.
func (plugin *HelloWorldPlugin) AfterInit() error {
	//...
}

// Close is called to cleanup the plugin resources.
func (plugin *HelloWorldPlugin) Close() error {
	//...
}
```
 
4. Create plugin instance and inject the dependencies

```
    	// create an instance of the plugin
    	hwPlugin := HelloWorldPlugin{}
    
    	// wire the dependencies, f refers to a flavour declared in the first step
    	hwPlugin.LogFactory = &f.Logrus
```

5. Pass all the plugins to the constructor and start the agent.

```
	// Create new agent
	agent := core.NewAgent(logroot.Logger(), 15*time.Second, append(f.Plugins(), &core.NamedPlugin{PluginName: PluginID, Plugin: &hwPlugin})...)

	core.EventLoopWithInterrupt(agent, nil)
```