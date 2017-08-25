# Cassandra REST API

Cassandra REST API is a REST based plugin which interacts with Cassandra database.
It uses rest/rpc flavor which includes logger and gorilla mux HTTP plugin, registers required handlers to handle HTTP requests.

1. Uses an existing rest/rpc flavour - collection of plugins
```go
	import "github.com/ligato/cn-infra/flavors/rpc"


	f := rpc.FlavorRPC{}
```
Alternatively, you can create a custom [flavour](https://github.com/ligato/cn-infra/tree/master/flavors).

2. Declare a structure for your plugin. Apart from the internal fields the structure must specify
the dependencies of the plugin.

```go
// CassandraRestAPIPlugin is a plugin that showcase the extensibility of vpp agent.
type CassandraRestAPIPlugin struct {

	// httpmux is a dependency of the plugin that needs to be injected.
	HTTPHandlers rest.HTTPHandlers

	// session gockle.Session stores the session to Cassandra
	session gockle.Session

	// broker stores the cassandra data broker
	broker *cassandra.BrokerCassa
}
```

3. Implement plugin lifecycle methods. `Init()` and `Close` are mandatory `AfterInit` is optional.

```go
// Init is called on plugin startup. New logger is instantiated.
func (plugin *CassandraRestAPIPlugin) Init() (err error) {
    // use injected dependencies in this case to create a logger
    plugin.Logger, err = plugin.LogFactory.NewLogger(string(PluginID))
	
	//...
	return err
}

// AfterInit logs a sample message.
func (plugin *CassandraRestAPIPlugin) AfterInit() error {
	//...
}

// Close is called to cleanup the plugin resources.
func (plugin *CassandraRestAPIPlugin) Close() error {
	//...
}
```
 
4. Create plugin instance and inject the dependencies.

```go
    	// create an instance of the plugin
    	cassRestAPIPlugin := CassandraRestAPIPlugin{}
    
    	// wire the dependencies, f refers to a flavour declared in the first step
    	cassRestAPIPlugin.HTTPHandlers = &f.HTTP
```

5. Pass all the plugins to the constructor and start the agent.

```go
	// Create new agent
	agent := core.NewAgent(logroot.Logger(), 15*time.Second, append(f.Plugins(), &core.NamedPlugin{PluginName: PluginID, Plugin: &hwPlugin})...)

	core.EventLoopWithInterrupt(agent, nil)
```

6. Register HTTP Handler

```go
    plugin.HTTPHandlers.RegisterHTTPHandler("/tweets", plugin.tweetsHandler, "GET")
    plugin.HTTPHandlers.RegisterHTTPHandler("/tweets/{id}", plugin.tweetsHandler, "GET")
    plugin.HTTPHandlers.RegisterHTTPHandler("/tweets", plugin.tweetsHandler, "POST")
    plugin.HTTPHandlers.RegisterHTTPHandler("/tweets/{id}", plugin.tweetsHandler, "PUT")
    plugin.HTTPHandlers.RegisterHTTPHandler("/tweets/{id}", plugin.tweetsHandler, "DELETE")
```

7. Create Configuration either using ClientConfig structure or by creating a client-config.yaml

```go
    clientConfig, configErr := createConfig()

    OR

    clientConfig, configErr := loadConfig("<client-config.yaml file path>")
```

8. Create session

```go
    session1, sessionErr := createSession(config)
```

9. Create data broker

```go
    db := cassandra.NewBrokerUsingSession(session1)
```

10. Use data broker to run DDL/DML statements in Cassandra database

```go
    err1 := db.Put(sql.FieldEQ(&insertTweet.ID), insertTweet)
```
Supported broker [API](https://github.com/ligato/cn-infra/blob/master/db/sql/cassandra/cassa_broker_impl.go)

## How to run micro-service

1. Navigate to the cn-sample-service/cmd/cassandra directory
2. Run go build
3. Run go install
4. Set environment variables
```
    export CASSANDRA_HOST=127.0.0.1
    export CASSANDRA_PORT=9042
    export CASSANDRA_CONFIG=<configuration file path>
```
5. Start the micro-service
```
    ./cassandra --cassandra-config=<cassandra.conf.yaml file path>
```
6. Start Cassandra locally using docker or your local installation
```
    sudo docker run -p 9042:9042 -it --name cassandra01 --rm cassandra:latest
```
7. Using curl post HTTP requests
```
curl -X POST http://localhost:9191/tweets
curl -X GET http://localhost:9191/tweets
curl -X GET http://localhost:9191/tweets/{id}
curl -X PUT http://localhost:9191/tweets/{id}
curl -X DELETE http://localhost:9191/tweets/{id}
curl -X POST http://localhost:9191/users
curl -X GET http://localhost:9191/users
curl -X GET http://localhost:9191/users/{id}
```



