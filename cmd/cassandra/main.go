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
	"github.com/ligato/cn-infra/flavors/rpc"
	"github.com/ligato/cn-infra/logging/logroot"
	"github.com/ligato/cn-infra/rpc/rest"
	"github.com/unrolled/render"
	"net/http"
	"time"
)

// PluginID of the custom govpp_call plugin
const PluginID core.PluginName = "cassandra-plugin"

// CassandraPlugin is a plugin that showcase the extensibility of vpp agent.
type CassandraPlugin struct {
	// LogFactory is a dependency of the plugin that needs to be injected.
	//LogFactory logging.LogFactory
	//logging.Logger
	//
	// httpmux is a dependency of the plugin that needs to be injected.
	HTTPHandlers rest.HTTPHandlers
}

//connectivityHandler defining route handler which performs basic connectivity test by reading/writing data to Cassandra
func connectivityHandler(formatter *render.Render) http.HandlerFunc {

	// An example HTTP request handler which prints out attributes of a trivial Go structure in JSON format.
	return func(w http.ResponseWriter, req *http.Request) {
		logroot.StandardLogger().Info("Testing connectivity by getting all tweets from tweets table.")

		err := connectivity()

		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
		} else {
			formatter.JSON(w, http.StatusOK, "Connectivity successful")
		}
	}
}

//alterTableHandler defining route handler which performs table alteration by adding a column to person table
func alterTableHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		logroot.StandardLogger().Info("Testing Alter Table by adding a new column to tweets table.")

		err := alterTable()

		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
		} else {
			formatter.JSON(w, http.StatusOK, "Alter table successful")
		}
	}
}

//keyspaceIfNotExistHandler defining route handler which indicates use of IF NOT EXISTS clause while creating a keyspace
func keyspaceIfNotExistHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		logroot.StandardLogger().Info("Testing use of IF NOT EXISTS clause.")

		err := createKeySpaceIfNotExist()

		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
		} else {
			formatter.JSON(w, http.StatusOK, "Keyspace successful")
		}
	}
}

//customDataStructureHandler defining route handler which indicates use of custom data structure and types
//used to return map of addresses as HTTP response
func customDataStructureHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		logroot.StandardLogger().Info("Testing use of user-defined data types.")

		addresses, err := insertCustomizedDataStructure()

		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
		} else {
			formatter.JSON(w, http.StatusOK, addresses)
		}
	}
}

//reconnectIntervalHandler defining route handler which allows configuring redial_interval for a session
func reconnectIntervalHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		logroot.StandardLogger().Info("Testing gocql reconnect interval behaviour.")

		err := reconnectInterval()

		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
		} else {
			formatter.JSON(w, http.StatusOK, "Reconnect Interval successful")
		}
	}
}

//queryTimeoutHandler defining route handler which allows configuring op_timeout for a session
func queryTimeoutHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		logroot.StandardLogger().Info("Testing gocql timeout behaviour.")

		err := queryTimeout()

		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
		} else {
			formatter.JSON(w, http.StatusOK, "Query Timeout successful")
		}
	}
}

//connectTimeoutHandler defining route handler which allows configuring dial_timeout for a session
func connectTimeoutHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		logroot.StandardLogger().Info("Testing gocql connect timeout behaviour.")

		err := connectTimeout()

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
	cassPlugin := CassandraPlugin{}

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
func (plugin *CassandraPlugin) Init() (err error) {
	//plugin.Logger = plugin.LogFactory.NewLogger(string(PluginID))
	plugin.HTTPHandlers.RegisterHTTPHandler("/connectivity", connectivityHandler, "GET")
	plugin.HTTPHandlers.RegisterHTTPHandler("/altertable", alterTableHandler, "GET")
	plugin.HTTPHandlers.RegisterHTTPHandler("/keyspaceifnotexists", keyspaceIfNotExistHandler, "GET")
	plugin.HTTPHandlers.RegisterHTTPHandler("/customdatastructure", customDataStructureHandler, "GET")
	plugin.HTTPHandlers.RegisterHTTPHandler("/reconnectinterval", reconnectIntervalHandler, "GET")
	plugin.HTTPHandlers.RegisterHTTPHandler("/querytimeout", queryTimeoutHandler, "GET")
	plugin.HTTPHandlers.RegisterHTTPHandler("/connecttimeout", connectTimeoutHandler, "GET")
	return err
}

// AfterInit logs a sample message.
func (plugin *CassandraPlugin) AfterInit() error {
	logroot.StandardLogger().Info("Cassandra Plugin is up and running !!!")
	return nil
}

// Close is called to cleanup the plugin resources.
func (plugin *CassandraPlugin) Close() error {
	return nil
}
