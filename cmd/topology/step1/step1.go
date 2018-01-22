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

type TopologyPlugin struct {
	local.PluginInfraDeps
}

func main() {

	// Create new agent using local flavor
	agent := local.NewAgent(local.WithPlugins(func(local *local.FlavorLocal) []*core.NamedPlugin {

		// create an instance of the plugin
		topoPlugin := &TopologyPlugin{}
		topoPlugin.PluginInfraDeps = *local.InfraDeps("topology-plugin")

		return []*core.NamedPlugin {
			{topoPlugin.PluginName, topoPlugin},
		}

	}))

	core.EventLoopWithInterrupt(agent, nil)
}

// Init is called on plugin startup. New logger is instantiated.
func (plugin *TopologyPlugin) Init() (err error) {
	plugin.Log.Info("Topology Plugin initialized.")
	return nil
}

// AfterInit logs a sample message.
func (plugin *TopologyPlugin) AfterInit() error {
	plugin.Log.Info("Topology Plugin is running.")
	return nil
}

// Close is called to cleanup the plugin resources.
func (plugin *TopologyPlugin) Close() error {
	return nil
}
