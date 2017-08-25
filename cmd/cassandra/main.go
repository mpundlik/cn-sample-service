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
	"github.com/ligato/cn-infra/config"
	"github.com/ligato/cn-infra/core"
	"github.com/ligato/cn-infra/db/sql"
	"github.com/ligato/cn-infra/db/sql/cassandra"
	"github.com/ligato/cn-infra/flavors/rpc"
	"github.com/ligato/cn-infra/logging/logroot"
	"github.com/ligato/cn-infra/rpc/rest"
	"os"
	"time"
)

type Deps struct {
	// httpmux is a dependency of the plugin that needs to be injected.
	HTTPHandlers rest.HTTPHandlers
	BrokerPlugin sql.BrokerPlugin
}

type CassandraRestFlavor struct {
	rpc.FlavorRPC
	CASSANDRA cassandra.Plugin
	CassandraRestAPIPlugin
	injected bool
}

// CassandraRestAPIPlugin is a plugin that showcase the extensibility of vpp agent.
type CassandraRestAPIPlugin struct {
	Deps
	// broker stores the cassandra data broker
	broker sql.Broker
}

//main entry point for the sample service
func main() {

	flavor := CassandraRestFlavor{}

	cassRESTPlugin := CassandraRestAPIPlugin{}
	cassRESTPlugin.Deps.HTTPHandlers = &flavor.HTTP
	cassRESTPlugin.Deps.BrokerPlugin = &flavor.CASSANDRA

	flavor.CassandraRestAPIPlugin = cassRESTPlugin

	// Create new agent
	agent := core.NewAgent(logroot.StandardLogger(), 15*time.Second, append(flavor.Plugins())...)

	err := core.EventLoopWithInterrupt(agent, nil)
	if err != nil {
		logroot.StandardLogger().Errorf("Error in event loop %v", err)
	}
}

// Init is called on plugin startup. New logger is instantiated and required HTTP handlers are registered.
func (plugin *CassandraRestAPIPlugin) Init() (err error) {
	return nil
}

// AfterInit logs a sample message.
func (plugin *CassandraRestAPIPlugin) AfterInit() error {
	logroot.StandardLogger().Info("Cassandra REST API Plugin is up and running !!!")

	plugin.HTTPHandlers.RegisterHTTPHandler("/tweets", plugin.tweetsHandler, "GET")
	plugin.HTTPHandlers.RegisterHTTPHandler("/tweets/{id}", plugin.tweetsHandler, "GET")
	plugin.HTTPHandlers.RegisterHTTPHandler("/tweets", plugin.tweetsHandler, "POST")
	plugin.HTTPHandlers.RegisterHTTPHandler("/tweets/{id}", plugin.tweetsHandler, "PUT")
	plugin.HTTPHandlers.RegisterHTTPHandler("/tweets/{id}", plugin.tweetsHandler, "DELETE")
	plugin.HTTPHandlers.RegisterHTTPHandler("/users", plugin.usersHandler, "GET")
	plugin.HTTPHandlers.RegisterHTTPHandler("/users/{id}", plugin.usersHandler, "GET")
	plugin.HTTPHandlers.RegisterHTTPHandler("/users", plugin.usersHandler, "POST")

	plugin.broker = plugin.BrokerPlugin.NewBroker()

	plugin.setup()

	return nil
}

// Close is called to cleanup the plugin resources.
func (plugin *CassandraRestAPIPlugin) Close() error {

	err := plugin.closeConnection()
	if err != nil {
		logroot.StandardLogger().Errorf("Error closing connection %v", err)
		os.Exit(1)
	}

	return nil
}

//setup used to setup Cassandra before running each request
func (plugin *CassandraRestAPIPlugin) setup() (err error) {
	db := plugin.broker

	err1 := db.Exec(`CREATE KEYSPACE IF NOT EXISTS example with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 }`)
	if err1 != nil {
		logroot.StandardLogger().Errorf("Error creating keyspace %v", err1)
		return err1
	}

	err2 := db.Exec(`CREATE TABLE IF NOT EXISTS example.tweet(timeline text, id text, text text, user text, PRIMARY KEY(id))`)
	if err2 != nil {
		logroot.StandardLogger().Errorf("Error creating table %v", err2)
		return err2
	}

	err4 := db.Exec(`CREATE INDEX IF NOT EXISTS ON example.tweet(timeline)`)
	if err4 != nil {
		logroot.StandardLogger().Errorf("Error creating index %v", err4)
		return err4
	}

	err5 := db.Exec(`CREATE KEYSPACE IF NOT EXISTS example2 with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 }`)
	if err5 != nil {
		logroot.StandardLogger().Errorf("Error creating keyspace %v", err5)
		return err5
	}

	err6 := db.Exec(`CREATE TYPE IF NOT EXISTS example2.phone (
			countryCode int,
			number text,
		)`)

	if err6 != nil {
		logroot.StandardLogger().Errorf("Error creating user-defined type phone %v", err6)
		return err6
	}

	err7 := db.Exec(`CREATE TYPE IF NOT EXISTS example2.address (
			street text,
			city text,
			zip text,
			phones map<text, frozen<phone>>
		)`)

	if err7 != nil {
		logroot.StandardLogger().Errorf("Error creating user-defined type address %v", err7)
		return err7
	}

	err8 := db.Exec(`CREATE TABLE IF NOT EXISTS example2.user (
			ID text PRIMARY KEY,
			addresses map<text, frozen<address>>
		)`)

	if err8 != nil {
		logroot.StandardLogger().Errorf("Error creating table user %v", err8)
		return err8
	}

	return nil
}

//closeConnection used to clean up and close connection to cassandra
func (plugin *CassandraRestAPIPlugin) closeConnection() (err error) {

	db := plugin.broker

	err1 := db.Exec(`DROP TABLE IF EXISTS example.tweet`)
	if err1 != nil {
		logroot.StandardLogger().Errorf("Error dropping table %v", err1)
		return err1
	}

	err2 := db.Exec(`DROP TABLE IF EXISTS example2.user`)
	if err2 != nil {
		logroot.StandardLogger().Errorf("Error dropping table %v", err2)
		return err2
	}

	err3 := db.Exec(`DROP TYPE IF EXISTS example2.address`)
	if err3 != nil {
		logroot.StandardLogger().Errorf("Error dropping type %v", err3)
		return err3
	}

	err4 := db.Exec(`DROP TYPE IF EXISTS example2.phone`)
	if err4 != nil {
		logroot.StandardLogger().Errorf("Error dropping type %v", err4)
		return err4
	}

	err5 := db.Exec(`DROP KEYSPACE IF EXISTS example`)
	if err5 != nil {
		logroot.StandardLogger().Errorf("Error dropping keyspace %v", err5)
		return err5
	}

	err6 := db.Exec(`DROP KEYSPACE IF EXISTS example2`)
	if err6 != nil {
		logroot.StandardLogger().Errorf("Error dropping keyspace %v", err6)
		return err6
	}

	return nil
}

// Inject sets object references
func (f *CassandraRestFlavor) Inject() (allReadyInjected bool) {
	if !f.FlavorRPC.Inject() {
		return false
	}

	f.CASSANDRA.Deps.PluginInfraDeps = *f.InfraDeps("cassandra")
	f.CASSANDRA.Deps.PluginInfraDeps.PluginConfig = config.ForPlugin("cassandra")

	return true
}

// Plugins combines all Plugins in flavor to the list
func (f *CassandraRestFlavor) Plugins() []*core.NamedPlugin {
	f.Inject()
	return core.ListPluginsInFlavor(f)
}
