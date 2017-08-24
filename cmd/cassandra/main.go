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
	"github.com/gorilla/mux"
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
func (plugin *CassandraRestAPIPlugin) tweetsHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		logroot.StandardLogger().Info("Received tweets request")

		pathParams := mux.Vars(req)
		logroot.StandardLogger().Infof("pathParams = %v", pathParams)

		switch req.Method {
		case "POST":
			err := insertTweets(plugin.broker)

			if err != nil {
				formatter.JSON(w, http.StatusInternalServerError, err.Error())
			} else {
				formatter.JSON(w, http.StatusOK, "Tweets inserted successfully")
			}
		case "GET":
			if pathParams != nil && len(pathParams) > 0 {
				id := pathParams["id"]
				if id != "" {
					result, err := getTweetByID(plugin.broker, id)

					if err != nil {
						formatter.JSON(w, http.StatusInternalServerError, err.Error())
					} else {
						formatter.JSON(w, http.StatusOK, result)
					}
				} else {
					formatter.JSON(w, http.StatusBadRequest, errors.New("id is nil"))
				}
			} else {
				result, err := getAllTweets(plugin.broker)

				if err != nil {
					formatter.JSON(w, http.StatusInternalServerError, err.Error())
				} else {
					formatter.JSON(w, http.StatusOK, result)
				}
			}
		case "PUT":
			if pathParams != nil && len(pathParams) > 0 {
				id := pathParams["id"]
				if id != "" {
					err := insertTweet(plugin.broker, id)

					if err != nil {
						formatter.JSON(w, http.StatusInternalServerError, err.Error())
					} else {
						formatter.JSON(w, http.StatusCreated, "Tweet inserted successfully")
					}
				} else {
					formatter.JSON(w, http.StatusBadRequest, errors.New("id is nil"))
				}
			} else {
				formatter.JSON(w, http.StatusBadRequest, errors.New("Request not supported"))
			}
		case "DELETE":
			if pathParams != nil && len(pathParams) > 0 {
				id := pathParams["id"]
				if id != "" {
					err := deleteTweetByID(plugin.broker, id)

					if err != nil {
						formatter.JSON(w, http.StatusInternalServerError, err.Error())
					} else {
						formatter.JSON(w, http.StatusOK, "Tweet deleted successfully")
					}
				} else {
					formatter.JSON(w, http.StatusBadRequest, errors.New("id is nil"))
				}
			} else {
				formatter.JSON(w, http.StatusBadRequest, errors.New("Request not supported"))
			}

		default:
			formatter.JSON(w, http.StatusMethodNotAllowed, nil)

		}
	}
}

//usersHandler defining route handler which indicates use of user defined types
//used to return map of addresses as HTTP response
func (plugin *CassandraRestAPIPlugin) usersHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {

		logroot.StandardLogger().Info("Received users request")

		pathParams := mux.Vars(req)
		logroot.StandardLogger().Infof("pathParams = %v", pathParams)

		switch req.Method {
		case "POST":
			err := insertUsers(plugin.broker)

			if err != nil {
				formatter.JSON(w, http.StatusInternalServerError, err.Error())
			} else {
				formatter.JSON(w, http.StatusOK, "Inserted users successfully")
			}
		case "GET":
			if pathParams != nil && len(pathParams) > 0 {
				id := pathParams["id"]
				if id != "" {
					result, err := getUserByID(plugin.broker, id)

					if err != nil {
						formatter.JSON(w, http.StatusInternalServerError, err.Error())
					} else {
						formatter.JSON(w, http.StatusOK, result)
					}
				} else {
					formatter.JSON(w, http.StatusBadRequest, errors.New("id is nil"))
				}
			} else {
				result, err := getAllUsers(plugin.broker)

				if err != nil {
					formatter.JSON(w, http.StatusInternalServerError, err.Error())
				} else {
					formatter.JSON(w, http.StatusOK, result)
				}
			}
		default:
			formatter.JSON(w, http.StatusMethodNotAllowed, nil)
		}
	}
}

//main entry point for the sample service
func main() {
	// leverage an existing flavour - set of plugins
	f := rpc.FlavorRPC{}

	// create an instance of the plugin
	cassRestAPIPlugin := CassandraRestAPIPlugin{}

	// wire the dependencies
	cassRestAPIPlugin.HTTPHandlers = &f.HTTP

	// creating an instance of the named plugin
	cassRestAPINamedPlugin := &core.NamedPlugin{PluginName: PluginID, Plugin: &cassRestAPIPlugin}

	// Create new agent
	agent := core.NewAgent(logroot.StandardLogger(), 15*time.Second, append(f.Plugins(), cassRestAPINamedPlugin)...)

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

	// create configuration from a client configuration file
	/*clientConfig, configErr := loadConfig("/Users/mpundlik/go/src/github.com/ligato/cn-sample-service/cmd/cassandra/client-config.yaml")
	if configErr != nil {
		logroot.StandardLogger().Errorf("Config err = %v", configErr)
		return configErr
	}*/

	//OR create configuration using config structure
	clientConfig, configErr := createConfig()
	if configErr != nil {
		logroot.StandardLogger().Errorf("Config err = %v", configErr)
		return configErr
	}

	logroot.StandardLogger().Infof("clientconfig = %v", clientConfig.ReconnectInterval)
	logroot.StandardLogger().Infof("clientconfig = %v", clientConfig.Timeout)
	logroot.StandardLogger().Infof("clientconfig = %v", clientConfig.ConnectTimeout)

	session1, setupErr := setup(clientConfig)
	if setupErr != nil {
		logroot.StandardLogger().Errorf("Setup error = %v", setupErr)
		return setupErr
	}

	plugin.session = session1

	db := cassandra.NewBrokerUsingSession(session1)
	plugin.broker = db

	plugin.HTTPHandlers.RegisterHTTPHandler("/tweets", plugin.tweetsHandler, "GET")
	plugin.HTTPHandlers.RegisterHTTPHandler("/tweets/{id}", plugin.tweetsHandler, "GET")
	plugin.HTTPHandlers.RegisterHTTPHandler("/tweets", plugin.tweetsHandler, "POST")
	plugin.HTTPHandlers.RegisterHTTPHandler("/tweets/{id}", plugin.tweetsHandler, "PUT")
	plugin.HTTPHandlers.RegisterHTTPHandler("/tweets/{id}", plugin.tweetsHandler, "DELETE")
	plugin.HTTPHandlers.RegisterHTTPHandler("/users", plugin.usersHandler, "GET")
	plugin.HTTPHandlers.RegisterHTTPHandler("/users/{id}", plugin.usersHandler, "GET")
	plugin.HTTPHandlers.RegisterHTTPHandler("/users", plugin.usersHandler, "POST")

	return nil
}

// Close is called to cleanup the plugin resources.
func (plugin *CassandraRestAPIPlugin) Close() error {

	err := closeConnection(plugin.session)
	if err != nil {
		logroot.StandardLogger().Errorf("Error closing connection %v", err)
		os.Exit(1)
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

	err4 := db.Exec(`CREATE INDEX IF NOT EXISTS ON example.tweet(timeline)`)
	if err4 != nil {
		logroot.StandardLogger().Errorf("Error creating index %v", err4)
		return nil, err4
	}

	err5 := db.Exec(`CREATE KEYSPACE IF NOT EXISTS example2 with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 }`)
	if err5 != nil {
		logroot.StandardLogger().Errorf("Error creating keyspace %v", err5)
		return nil, err5
	}

	err6 := db.Exec(`CREATE TYPE IF NOT EXISTS example2.phone (
			countryCode int,
			number text,
		)`)

	if err6 != nil {
		logroot.StandardLogger().Errorf("Error creating user-defined type phone %v", err6)
		return nil, err6
	}

	err7 := db.Exec(`CREATE TYPE IF NOT EXISTS example2.address (
			street text,
			city text,
			zip text,
			phones map<text, frozen<phone>>
		)`)

	if err7 != nil {
		logroot.StandardLogger().Errorf("Error creating user-defined type address %v", err7)
		return nil, err7
	}

	err8 := db.Exec(`CREATE TABLE IF NOT EXISTS example2.user (
			ID text PRIMARY KEY,
			addresses map<text, frozen<address>>
		)`)

	if err8 != nil {
		logroot.StandardLogger().Errorf("Error creating table user %v", err8)
		return nil, err8
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
		DialTimeout:    600000000,
		OpTimeout:      60,
		RedialInterval: 60000000000,
	}

	clientConfig, err2 := cassandra.ConfigToClientConfig(config1)
	if err != nil {
		logroot.StandardLogger().Errorf("Error in converting from config to ClientConfig")
		return nil, err2
	}

	return clientConfig, nil
}

//loadConfig used to create configuration structure from configuration file
/*func loadConfig(configFileName string) (*cassandra.ClientConfig, error) {
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
