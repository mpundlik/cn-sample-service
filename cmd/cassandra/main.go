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
	"github.com/ligato/cn-infra/core"
	"github.com/ligato/cn-infra/db/sql/cassandra"
	"github.com/ligato/cn-infra/flavors/rpc"
	"github.com/ligato/cn-infra/logging/logroot"
	"github.com/ligato/cn-infra/rpc/rest"
	"github.com/unrolled/render"
	"github.com/willfaught/gockle"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// PluginID of the Cassandra REST API Plugin
const PluginID core.PluginName = "cassandra-rest-api-plugin"

// CassandraRestAPIPlugin is a plugin that showcase the extensibility of vpp agent.
type CassandraRestAPIPlugin struct {

	// httpmux is a dependency of the plugin that needs to be injected.
	HTTPHandlers rest.HTTPHandlers

	// session gockle.Session stores the session to Cassandra
	session gockle.Session

	// broker stores the cassandra data broker
	broker *cassandra.BrokerCassa
}

//connectivityHandler defining route handler which performs basic connectivity test by reading/writing data to Cassandra
func (plugin *CassandraRestAPIPlugin) connectivityHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		logroot.StandardLogger().Info("Testing connectivity by getting all tweets from tweets table.")

		err := connectivity(plugin.broker)

		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
		} else {
			formatter.JSON(w, http.StatusOK, "Connectivity successful")
		}
	}
}

//alterTableHandler defining route handler which performs table alteration by adding a column to person table
func (plugin *CassandraRestAPIPlugin) alterTableHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		logroot.StandardLogger().Info("Testing Alter Table by adding a new column to tweets table.")

		err := alterTable(plugin.broker)

		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
		} else {
			formatter.JSON(w, http.StatusOK, "Alter table successful")
		}
	}
}

//keyspaceIfNotExistHandler defining route handler which indicates use of IF NOT EXISTS clause while creating a keyspace
func (plugin *CassandraRestAPIPlugin) keyspaceIfNotExistHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		logroot.StandardLogger().Info("Testing use of IF NOT EXISTS clause.")

		err := createKeySpaceIfNotExist(plugin.broker)

		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
		} else {
			formatter.JSON(w, http.StatusOK, "Keyspace successful")
		}
	}
}

//customDataStructureHandler defining route handler which indicates use of custom data structure and types
//used to return map of addresses as HTTP response
func (plugin *CassandraRestAPIPlugin) customDataStructureHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		logroot.StandardLogger().Info("Testing use of user-defined data types.")

		addresses, err := insertCustomizedDataStructure(plugin.broker)

		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
		} else {
			formatter.JSON(w, http.StatusOK, addresses)
		}
	}
}

//reconnectIntervalHandler defining route handler which allows configuring redial_interval for a session
func (plugin *CassandraRestAPIPlugin) reconnectIntervalHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		logroot.StandardLogger().Info("Testing gocql reconnect interval behaviour.")

		err := reconnectInterval(plugin.broker)

		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
		} else {
			formatter.JSON(w, http.StatusOK, "Reconnect Interval successful")
		}
	}
}

//queryTimeoutHandler defining route handler which allows configuring op_timeout for a session
func (plugin *CassandraRestAPIPlugin) queryTimeoutHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		logroot.StandardLogger().Info("Testing gocql timeout behaviour.")

		err := queryTimeout(plugin.broker)

		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
		} else {
			formatter.JSON(w, http.StatusOK, "Query Timeout successful")
		}
	}
}

//connectTimeoutHandler defining route handler which allows configuring dial_timeout for a session
func (plugin *CassandraRestAPIPlugin) connectTimeoutHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		logroot.StandardLogger().Info("Testing gocql connect timeout behaviour.")

		err := connectTimeout(plugin.broker)

		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
		} else {
			formatter.JSON(w, http.StatusOK, "Connect Timeout successful")
		}
	}
}

//main entry point for the sample service
func main() {
	// leverage an existing flavour - set of plugins
	f := rpc.FlavorRPC{}

	// create an instance of the plugin
	cassPlugin := CassandraRestAPIPlugin{}

	// wire the dependencies
	cassPlugin.HTTPHandlers = &f.HTTP

	// creating an instance of the named plugin
	cassNamedPlugin := &core.NamedPlugin{PluginName: PluginID, Plugin: &cassPlugin}

	// Create new agent
	agent := core.NewAgent(logroot.StandardLogger(), 15*time.Second, append(f.Plugins(), cassNamedPlugin)...)

	err := core.EventLoopWithInterrupt(agent, nil)
	if err != nil {
		logroot.StandardLogger().Errorf("Error in event loop %v", err)
	}
}

// Init is called on plugin startup. New logger is instantiated and required HTTP handlers are registered.
func (plugin *CassandraRestAPIPlugin) Init() (err error) {

	return err
}

// AfterInit logs a sample message.
func (plugin *CassandraRestAPIPlugin) AfterInit() error {
	logroot.StandardLogger().Info("Cassandra REST API Plugin is up and running !!!")

	//create configuration using config structure
	clientConfig, configErr := createConfig()
	if configErr != nil {
		logroot.StandardLogger().Errorf("Config err = %v", configErr)
		return configErr
	}

	//OR create configuration from a client configuration file
	/*clientConfig, configErr := loadConfig("/Users/mpundlik/go/src/github.com/ligato/cn-sample-service/cmd/cassandra/client-config.yaml")
	if configErr != nil {
		logroot.StandardLogger().Errorf("Config err = %v", configErr)
		return configErr
	}*/

	session1, setupErr := setup(clientConfig)
	if setupErr != nil {
		logroot.StandardLogger().Errorf("Setup error = %v", setupErr)
		return setupErr
	}

	plugin.session = session1

	db := cassandra.NewBrokerUsingSession(session1)
	plugin.broker = db

	plugin.HTTPHandlers.RegisterHTTPHandler("/connectivity", plugin.connectivityHandler, "GET")
	plugin.HTTPHandlers.RegisterHTTPHandler("/altertable", plugin.alterTableHandler, "GET")
	plugin.HTTPHandlers.RegisterHTTPHandler("/keyspaceifnotexists", plugin.keyspaceIfNotExistHandler, "GET")
	plugin.HTTPHandlers.RegisterHTTPHandler("/customdatastructure", plugin.customDataStructureHandler, "GET")
	plugin.HTTPHandlers.RegisterHTTPHandler("/reconnectinterval", plugin.reconnectIntervalHandler, "GET")
	plugin.HTTPHandlers.RegisterHTTPHandler("/querytimeout", plugin.queryTimeoutHandler, "GET")
	plugin.HTTPHandlers.RegisterHTTPHandler("/connecttimeout", plugin.connectTimeoutHandler, "GET")

	return nil
}

// Close is called to cleanup the plugin resources.
func (plugin *CassandraRestAPIPlugin) Close() error {

	err := closeConnection(plugin.session)
	if err != nil {
		logroot.StandardLogger().Errorf("Error closing connection %v", err)
	}

	return nil
}

//setup used to setup Cassandra before running each request
func setup(config *cassandra.ClientConfig) (session gockle.Session, err error) {
	session1, sessionErr := createSession(config)
	if sessionErr != nil {
		logroot.StandardLogger().Errorf("Error creating session %v", sessionErr)
		return nil, sessionErr
	}

	db := cassandra.NewBrokerUsingSession(session1)

	err1 := db.Exec(`CREATE KEYSPACE IF NOT EXISTS example with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 }`)
	if err1 != nil {
		logroot.StandardLogger().Errorf("Error creating keyspace %v", err1)
		return nil, err1
	}

	err2 := db.Exec(`CREATE TABLE IF NOT EXISTS example.tweet(timeline text, id text, text text, user text, PRIMARY KEY(id))`)
	if err2 != nil {
		logroot.StandardLogger().Errorf("Error creating table %v", err2)
		return nil, err2
	}

	err3 := db.Exec(`CREATE TABLE IF NOT EXISTS example.person(id text, name text, PRIMARY KEY(id))`)
	if err3 != nil {
		logroot.StandardLogger().Errorf("Error creating table %v", err3)
		return nil, err3
	}

	err4 := db.Exec(`CREATE INDEX IF NOT EXISTS ON example.tweet(timeline)`)
	if err4 != nil {
		logroot.StandardLogger().Errorf("Error creating index %v", err4)
		return nil, err4
	}

	return session1, err
}

//closeConnection used to clean up and close connection to cassandra
func closeConnection(session gockle.Session) (err error) {

	defer session.Close()

	db := cassandra.NewBrokerUsingSession(session)

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

//createSession used to create a session/connection with the given cassandra client configuration
func createSession(config *cassandra.ClientConfig) (session gockle.Session, err error) {

	session1, err2 := cassandra.CreateSessionFromConfig(config)

	if err2 != nil {
		logroot.StandardLogger().Errorf("Error creating session %v", err2)
		return nil, err2
	}

	session2 := gockle.NewSession(session1)

	return session2, nil
}

//createConfig depicts use of creating a configuration structure
func createConfig() (config *cassandra.ClientConfig, err error) {
	// connect to the cluster
	cassandraHost := os.Getenv("CASSANDRA_HOST")
	cassandraPort := os.Getenv("CASSANDRA_PORT")
	logroot.StandardLogger().Infof("Using cassandra host from environment variable %v", cassandraHost)
	logroot.StandardLogger().Infof("Using cassandra port from environment variable %v", cassandraPort)

	endpoints := strings.Split(cassandraHost, ",")

	if cassandraPort == "" {
		logroot.StandardLogger().Infof("Using default port, since CASSANDRA_PORT environment variable is not set")
		cassandraPort = "9042"
	}

	port, portErr := strconv.Atoi(cassandraPort)
	if portErr != nil {
		logroot.StandardLogger().Errorf("Error getting cassandra port %v", portErr)
		return nil, portErr
	}

	config1 := &cassandra.Config{
		Endpoints:      endpoints,
		Port:           port,
		DialTimeout:    600,
		OpTimeout:      60,
		RedialInterval: 60,
	}

	clientConfig, err2 := cassandra.ConfigToClientConfig(config1)
	if err != nil {
		logroot.StandardLogger().Errorf("Error in converting from config to ClientConfig")
		return nil, err2
	}

	return clientConfig, nil
}

/* DO NOT DELETE - kept as an example
//loadConfig used to create configuration structure from configuration file
func loadConfig(configFileName string) (*cassandra.ClientConfig, error) {
	var cfg cassandra.Config

	err := config.ParseConfigFromYamlFile(configFileName, &cfg)
	if err != nil {
		logroot.StandardLogger().Errorf("Error parsing the yaml client configuration file")
		return nil, err
	}

	clientConfig, err2 := cassandra.ConfigToClientConfig(&cfg)
	if err != nil {
		logroot.StandardLogger().Errorf("Error in converting from config to ClientConfig")
		return nil, err2
	}

	return clientConfig, nil
}*/
