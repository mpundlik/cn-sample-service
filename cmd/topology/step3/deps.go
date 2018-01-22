package main

import (
	"github.com/ligato/cn-infra/core"
	"github.com/ligato/cn-infra/datasync"
	"github.com/ligato/cn-infra/datasync/kvdbsync"
	"github.com/ligato/cn-infra/datasync/resync"
	"github.com/ligato/cn-infra/db/keyval/etcdv3"
	"github.com/ligato/cn-infra/flavors/local"
	"github.com/namsral/flag"
)

var (
	etcdv3Config string
)

// Deps lists dependencies of TopologyPlugin.
type Deps struct {
	local.PluginInfraDeps                             // injected
	Publisher             datasync.KeyProtoValWriter  // injected - To write ETCD data
	Watcher               datasync.KeyValProtoWatcher // injected - To watch ETCD data
}

// TopologyFlavor is a set of plugins required for the topology example. It could be used AllConnectionFlavors
// but don't need to pull all the dependencies plugin doesn't need.
type TopologyFlavor struct {
	// Local flavor to access the Infra (logger, service label, status check)
	*local.FlavorLocal
	// Resync orchestrator
	ResyncOrch resync.Plugin
	// Etcd plugin
	ETCD etcdv3.Plugin
	// Etcd sync which manages and injects connection
	ETCDDataSync kvdbsync.Plugin
	// Topology plugin
	TopologyExample TopologyPlugin
	// Use channel when the topology plugin is finished
	closeChan *chan struct{}
}

// Plugins combines all plugins in the flavor into a slice.
func (tf *TopologyFlavor) Plugins() []*core.NamedPlugin {
	tf.Inject()
	return core.ListPluginsInFlavor(tf)
}


//init defines cassandra flags // TODO switch to viper to avoid global configuration
func init() {
	flag.StringVar(&etcdv3Config, "etcdv3-config", "etcd.conf",
		"Location of the Etcd Client configuration file; also set via 'ETCD_CONFIG' env variable.")
}