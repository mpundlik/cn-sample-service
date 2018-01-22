// Copyright (c) 2017 Cisco and/or its affiliates.
//
// Licensed under the Apache License, Version 2.0 (the "License");
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
	"errors"
	"github.com/ligato/cn-infra/core"
	"github.com/ligato/cn-infra/db/sql"
	"os"
	"fmt"
	"github.com/ligato/cn-infra/flavors/rpc"
	"github.com/ligato/cn-infra/db/sql/cassandra"
	"github.com/ligato/cn-infra/flavors/local"
)

// CassandraRestAPIPlugin is a plugin that showcase the extensibility of vpp agent.
type CassandraRestAPIPlugin struct {
	Deps

	// broker stores the cassandra data broker
	broker          sql.Broker
	pluginCompleted chan struct{}
	tweetTable      *tweet
	userTable       *user
}

//main entry point for the sample service
func main() {
	pluginCompleted := make(chan struct{}, 1)

	// Create new agent
	agent := rpc.NewAgent(rpc.WithPlugins(func(rpc *rpc.FlavorRPC) []*core.NamedPlugin {

		cassandraPlugin := &cassandra.Plugin{}
		cassandraPlugin.Deps.PluginInfraDeps = *rpc.InfraDeps("cassandra-plugin", local.WithConf())

		cassandraRestAPIPlugin := &CassandraRestAPIPlugin{}
		cassandraRestAPIPlugin.pluginCompleted = pluginCompleted
		cassandraRestAPIPlugin.Deps.PluginInfraDeps = *rpc.InfraDeps("cassandra-rest-api-plugin")
		cassandraRestAPIPlugin.Deps.HTTPHandlers = &rpc.HTTP
		cassandraRestAPIPlugin.Deps.BrokerPlugin = cassandraPlugin

		return []*core.NamedPlugin{
			{cassandraPlugin.PluginName, cassandraPlugin},
			{cassandraRestAPIPlugin.PluginName, cassandraRestAPIPlugin},
		}

	}))

	err := core.EventLoopWithInterrupt(agent, pluginCompleted)
	if err != nil {
		fmt.Printf("error in event loop %v", err)
		os.Exit(1)
	}
}

// Init is called on plugin startup. New logger is instantiated and required HTTP handlers are registered.
func (plugin *CassandraRestAPIPlugin) Init() (err error) {
	return nil
}

// AfterInit logs a sample message.
func (plugin *CassandraRestAPIPlugin) AfterInit() error {
	plugin.Log.Info("Cassandra REST API Plugin is up and running !!!")

	plugin.broker = plugin.BrokerPlugin.NewBroker()

	plugin.tweetTable = &tweet{}
	plugin.userTable = &user{}

	plugin.HTTPHandlers.RegisterHTTPHandler("/tweets", plugin.tweetsGetHandler, "GET")
	plugin.HTTPHandlers.RegisterHTTPHandler("/tweets/{id}", plugin.tweetsGetHandler, "GET")
	plugin.HTTPHandlers.RegisterHTTPHandler("/tweets", plugin.tweetsPostHandler, "POST")
	plugin.HTTPHandlers.RegisterHTTPHandler("/tweets/{id}", plugin.tweetsPutHandler, "PUT")
	plugin.HTTPHandlers.RegisterHTTPHandler("/tweets/{id}", plugin.tweetsDeleteHandler, "DELETE")
	plugin.HTTPHandlers.RegisterHTTPHandler("/users", plugin.usersGetHandler, "GET")
	plugin.HTTPHandlers.RegisterHTTPHandler("/users/{id}", plugin.usersGetHandler, "GET")
	plugin.HTTPHandlers.RegisterHTTPHandler("/users/{id}", plugin.usersPutHandler, "PUT")
	plugin.HTTPHandlers.RegisterHTTPHandler("/users", plugin.usersPostHandler, "POST")

	err := plugin.setup()
	if err != nil {
		return err
	}

	return nil
}

// Close is called to cleanup the plugin resources.
func (plugin *CassandraRestAPIPlugin) Close() error {

	err := plugin.teardown()
	if err != nil {
		return err
	}

	plugin.pluginCompleted <- struct{}{}

	return nil
}

//setup used to setup Cassandra before running each request
func (plugin *CassandraRestAPIPlugin) setup() (err error) {

	err1 := plugin.broker.Exec(`CREATE KEYSPACE IF NOT EXISTS example with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 }`)
	if err1 != nil {
		plugin.Log.Errorf("Error creating keyspace %v", err1)
		return err1
	}

	err2 := plugin.broker.Exec(`CREATE TABLE IF NOT EXISTS example.tweet(timeline text, id text, text text, user text, PRIMARY KEY(id))`)
	if err2 != nil {
		plugin.Log.Errorf("Error creating table %v", err2)
		return err2
	}

	err4 := plugin.broker.Exec(`CREATE INDEX IF NOT EXISTS ON example.tweet(timeline)`)
	if err4 != nil {
		plugin.Log.Errorf("Error creating index %v", err4)
		return err4
	}

	err5 := plugin.broker.Exec(`CREATE KEYSPACE IF NOT EXISTS example2 with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 }`)
	if err5 != nil {
		plugin.Log.Errorf("Error creating keyspace %v", err5)
		return err5
	}

	err6 := plugin.broker.Exec(`CREATE TYPE IF NOT EXISTS example2.phone (
			countryCode int,
			number text,
		)`)

	if err6 != nil {
		plugin.Log.Errorf("error creating user-defined type phone %v", err6)
		return err6
	}

	err7 := plugin.broker.Exec(`CREATE TYPE IF NOT EXISTS example2.address (
			street text,
			city text,
			zip text,
			phones map<text, frozen<phone>>
		)`)

	if err7 != nil {
		plugin.Log.Errorf("error creating user-defined type address %#v", err7)
		return err7
	}

	err8 := plugin.broker.Exec(`CREATE TABLE IF NOT EXISTS example2.user (
			ID text PRIMARY KEY,
			addresses map<text, frozen<address>>
		)`)

	if err8 != nil {
		plugin.Log.Errorf("error creating table user %v", err8)
		return err8
	}

	return nil
}

//teardown used to clean up tables/schema from cassandra
func (plugin *CassandraRestAPIPlugin) teardown() (err error) {

	if plugin.broker == nil {
		plugin.Log.Errorf("cassandra broker is nil")
		return errors.New("cassandra broker is nil")
	}

	err1 := plugin.broker.Exec(`DROP TABLE IF EXISTS example.tweet`)
	if err1 != nil {
		plugin.Log.Errorf("error dropping table %v", err1)
		return err1
	}

	err2 := plugin.broker.Exec(`DROP TABLE IF EXISTS example2.user`)
	if err2 != nil {
		plugin.Log.Errorf("error dropping table %v", err2)
		return err2
	}

	err3 := plugin.broker.Exec(`DROP TYPE IF EXISTS example2.address`)
	if err3 != nil {
		plugin.Log.Errorf("error dropping type %v", err3)
		return err3
	}

	err4 := plugin.broker.Exec(`DROP TYPE IF EXISTS example2.phone`)
	if err4 != nil {
		plugin.Log.Errorf("error dropping type %v", err4)
		return err4
	}

	err5 := plugin.broker.Exec(`DROP KEYSPACE IF EXISTS example`)
	if err5 != nil {
		plugin.Log.Errorf("error dropping keyspace %v", err5)
		return err5
	}

	err6 := plugin.broker.Exec(`DROP KEYSPACE IF EXISTS example2`)
	if err6 != nil {
		plugin.Log.Errorf("error dropping keyspace %v", err6)
		return err6
	}

	return nil
}


