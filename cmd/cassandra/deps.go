package main

import (
	"github.com/ligato/cn-infra/rpc/rest"
	"github.com/ligato/cn-infra/flavors/local"
	"github.com/ligato/cn-infra/db/sql"
	"github.com/ligato/cn-infra/flavors/rpc"
	"github.com/ligato/cn-infra/db/sql/cassandra"
	"github.com/ligato/cn-infra/core"
	"github.com/namsral/flag"
)

var (
	cassandraConfig string
)

type Deps struct {
	// httpmux is a dependency of the plugin that needs to be injected.
	local.PluginLogDeps
	HTTPHandlers rest.HTTPHandlers
	BrokerPlugin sql.BrokerPlugin
}

type CassandraRestFlavor struct {
	rpc.FlavorRPC
	CASSANDRA cassandra.Plugin
	CassandraRestAPIPlugin
}

// Inject sets object references
func (f *CassandraRestFlavor) Inject() (isInjected bool) {
	if !f.FlavorRPC.Inject() {
		return false
	}

	f.CassandraRestAPIPlugin.Deps.HTTPHandlers = &f.HTTP
	f.CassandraRestAPIPlugin.Deps.BrokerPlugin = &f.CASSANDRA
	f.CASSANDRA.Deps.PluginInfraDeps = *f.InfraDeps("cassandra")
	f.CassandraRestAPIPlugin.Deps.PluginLogDeps = *f.LogDeps("cassandra-rest-api-plugin")

	return true
}

// Plugins combines all Plugins in flavor to the list
func (f *CassandraRestFlavor) Plugins() []*core.NamedPlugin {
	f.Inject()
	return core.ListPluginsInFlavor(f)
}

//init defines cassandra flags // TODO switch to viper to avoid global configuration
func init() {
	flag.StringVar(&cassandraConfig, "cassandra-config", "cassandra.conf.yaml",
		"Location of the Cassandra Client configuration file; also set via 'CASSANDRA_CONFIG' env variable.")
}