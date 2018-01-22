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
)

func main() {

	// Create new agent from local flavor
	agent := local.NewAgent(local.WithPlugins(func(local *local.FlavorLocal) []*core.NamedPlugin {

		// create an instance of the plugin
		hwPlugin := &HelloWorldPlugin{}
		hwPlugin.Deps.PluginLogDeps = *local.LogDeps("hello-world-example")

		return []*core.NamedPlugin{
			{hwPlugin.PluginName, hwPlugin}}
	}))

	core.EventLoopWithInterrupt(agent, nil)
}

// PluginID of the custom govpp_call plugin
const PluginID core.PluginName = "helloworld-plugin"

type Deps struct {
	local.PluginLogDeps
}

// HelloWorldPlugin is a plugin that showcase the extensibility of vpp agent.
type HelloWorldPlugin struct {
	Deps
}

// Init is called on plugin startup. New logger is instantiated.
func (plugin *HelloWorldPlugin) Init() (err error) {
	plugin.Log.Info("HelloWorldPlugin initialized.")
	return nil
}

// AfterInit logs a sample message.
func (plugin *HelloWorldPlugin) AfterInit() error {
	plugin.Log.Info("Hello World!!!")
	return nil
}

// Close is called to cleanup the plugin resources.
func (plugin *HelloWorldPlugin) Close() error {
	return nil
}
