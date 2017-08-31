# HelloWorld

HelloWorld as the name suggest is a simple extension of cn-infra. It adds a single plugin that creates
a custom logger an logs a message "Hello World!!!". The aim of this example is to list basic steps that
are necessary in order to integrate a custom plugin to cn-infra.

1. Use an existing flavour - collection of plugins
```go
	import "github.com/ligato/cn-infra/flavors/local"


	f := local.FlavorLocal{}
```
Alternatively, you can create a custom [flavour](https://github.com/ligato/cn-infra/tree/master/flavors).

2. Declare a structure for your plugin. Apart from the internal fields the structure must specify
the dependencies of the plugin.

```go
// HelloWorldPlugin is a plugin that showcase the extensibility of cn-infra.
type HelloWorldPlugin struct {
	// LogFactory is a dependency of the plugin that needs to be injected.
	// This dependency provides an API to create a logger instance
	LogFactory logging.LogFactory

    //..
}
```

3. Implement plugin lifecycle methods. `Init()` and `Close` are mandatory `AfterInit` is optional.

```go
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
 
4. Create plugin instance and inject the dependencies.

```go
    	// create an instance of the plugin
    	hwPlugin := HelloWorldPlugin{}
    
    	// wire the dependencies, f refers to a flavour declared in the first step
    	hwPlugin.LogFactory = &f.Logrus
```

5. Pass all the plugins to the constructor and start the agent.

```go
	// Create new agent
	agent := core.NewAgent(logroot.Logger(), 15*time.Second, append(f.Plugins(), &core.NamedPlugin{PluginName: PluginID, Plugin: &hwPlugin})...)

	core.EventLoopWithInterrupt(agent, nil)
```