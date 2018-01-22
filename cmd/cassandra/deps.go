package main

import (
	"github.com/ligato/cn-infra/rpc/rest"
	"github.com/ligato/cn-infra/flavors/local"
	"github.com/ligato/cn-infra/db/sql"
	"github.com/namsral/flag"
)

var (
	cassandraConfig string
)

type Deps struct {
	// httpmux is a dependency of the plugin that needs to be injected.
	local.PluginInfraDeps
	HTTPHandlers rest.HTTPHandlers
	BrokerPlugin sql.BrokerPlugin
}

//init defines cassandra flags // TODO switch to viper to avoid global configuration
func init() {
	flag.StringVar(&cassandraConfig, "cassandra-config", "cassandra.conf.yaml",
		"Location of the Cassandra Client configuration file; also set via 'CASSANDRA_CONFIG' env variable.")
}