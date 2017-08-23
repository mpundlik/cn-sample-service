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
	"github.com/ligato/cn-infra/flavors/local"
	"github.com/ligato/cn-infra/logging/logroot"
	"time"
)

func main() {
	// leverage an existing flavour - set of plugins
	f := local.FlavorLocal{}

	// create an instance of the plugin
	hwPlugin := HelloWorldPlugin{}

	// Create new agent
	agent := core.NewAgent(logroot.StandardLogger(), 15*time.Second, append(f.Plugins(), &core.NamedPlugin{PluginName: PluginID, Plugin: &hwPlugin})...)

	core.EventLoopWithInterrupt(agent, nil)
}

// PluginID of the custom govpp_call plugin
const PluginID core.PluginName = "helloworld-plugin"

// HelloWorldPlugin is a plugin that showcase the extensibility of vpp agent.
type HelloWorldPlugin struct {
}

// Init is called on plugin startup. New logger is instantiated.
func (plugin *HelloWorldPlugin) Init() (err error) {
	logroot.StandardLogger().Info("HelloWorldPlugin initialized.")
	return nil
}

// AfterInit logs a sample message.
func (plugin *HelloWorldPlugin) AfterInit() error {
	logroot.StandardLogger().Info("Hello World!!!")
	return nil
}

// Close is called to cleanup the plugin resources.
func (plugin *HelloWorldPlugin) Close() error {
	return nil
}
