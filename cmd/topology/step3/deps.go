package main

import (
	"github.com/ligato/cn-infra/core"
	"github.com/ligato/cn-infra/datasync"
	"github.com/ligato/cn-infra/datasync/kvdbsync"
	"github.com/ligato/cn-infra/datasync/resync"
	"github.com/ligato/cn-infra/db/keyval/etcdv3"
	"github.com/ligato/cn-infra/flavors/connectors"
	"github.com/ligato/cn-infra/flavors/local"
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

// Inject sets inter-plugin references.
func (tf *TopologyFlavor) Inject() (allReadyInjected bool) {
	// Init local flavor
	if tf.FlavorLocal == nil {
		tf.FlavorLocal = &local.FlavorLocal{}
	}
	tf.FlavorLocal.Inject()

	// Init Resync, ETCD + ETCD sync
	tf.ResyncOrch.Deps.PluginLogDeps = *tf.FlavorLocal.LogDeps("resync-orch")
	tf.ETCD.Deps.PluginInfraDeps = *tf.InfraDeps("etcdv3")
	connectors.InjectKVDBSync(&tf.ETCDDataSync, &tf.ETCD, tf.ETCD.PluginName, tf.FlavorLocal, &tf.ResyncOrch)

	// Inject infra + transport (publisher, watcher) to example plugin
	tf.TopologyExample.PluginInfraDeps = *tf.FlavorLocal.InfraDeps("topology-plugin")
	tf.TopologyExample.Publisher = &tf.ETCDDataSync
	tf.TopologyExample.Watcher = &tf.ETCDDataSync
	tf.TopologyExample.closeChannel = tf.closeChan

	return true
}

// Plugins combines all plugins in the flavor into a slice.
func (tf *TopologyFlavor) Plugins() []*core.NamedPlugin {
	tf.Inject()
	return core.ListPluginsInFlavor(tf)
}
