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

package defaultplugins

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/ligato/cn-infra/logging/logrus"
	intf "github.com/ligato/vpp-agent/plugins/defaultplugins/ifplugin/model/interfaces"
	"github.com/ligato/vpp-agent/plugins/defaultplugins/l2plugin/model/l2"
)

const kafkaIfStateTopic = "if_state" // Kafka topic where interface state changes are published.

// Resync deletes obsolete operation status of network interfaces in DB
// Obsolete state is one that is not part of SwIfIndex
func (plugin *Plugin) resyncIfStateEvents(keys []string) error {
	for _, key := range keys {
		ifaceName, err := intf.ParseNameFromKey(key)
		if err != nil {
			return err
		}

		_, _, found := plugin.swIfIndexes.LookupIdx(ifaceName)
		if !found {
			log.Debug("deleting obsolete status begin ", key)
			err := plugin.Transport.PublishData(key, nil /*means delete*/)
			log.Debug("deleting obsolete status end ", key, err)
		} else {
			log.WithField("ifaceName", ifaceName).Debug("interface status is needed")
		}
	}

	return nil
}

// publishIfState goroutine is used to watch interface state notifications that are propagated to Kafka topic
func (plugin *Plugin) publishIfStateEvents(ctx context.Context) {
	plugin.wg.Add(1)
	defer plugin.wg.Done()

	for {
		select {
		case ifState := <-plugin.ifStateChan:
			plugin.Transport.PublishData(intf.InterfaceStateKey(ifState.State.Name), ifState.State)

			// marshall data into JSON & send kafka message
			if plugin.kafkaConn != nil && ifState.Type == intf.UPDOWN {
				json, err := json.Marshal(ifState.State)
				if err != nil {
					log.Error(err)
				} else {

					// send kafka message
					_, err = plugin.kafkaConn.SendSyncString(kafkaIfStateTopic,
						fmt.Sprintf("%s", ifState.State.Name), string(json))
					if err != nil {
						log.Error(err)
					} else {
						log.Debug("Sending Kafka notification")
					}
				}
			}

		case <-ctx.Done():
			// stop watching for state data updates
			return
		}
	}
}

// Resync deletes old operation status of bridge domains in ETCD
func (plugin *Plugin) resyncBdStateEvents(keys []string) error {
	for _, key := range keys {
		bdName, err := intf.ParseNameFromKey(key)
		if err != nil {
			return err
		}
		_, _, found := plugin.bdIndexes.LookupIdx(bdName)
		if !found {
			log.Debug("deleting obsolete status begin ", key)
			err := plugin.Transport.PublishData(key, nil)
			log.Debug("deleting obsolete status end ", key, err)
		} else {
			log.WithField("bdName", bdName).Debug("bridge domain status required")
		}
	}

	return nil
}

// PublishBdState is used to watch bridge domain state notifications
func (plugin *Plugin) publishBdStateEvents(ctx context.Context) {
	plugin.wg.Add(1)
	defer plugin.wg.Done()

	for {
		select {
		case bdState := <-plugin.bdStateChan:
			if bdState != nil && bdState.State != nil {
				key := l2.BridgeDomainStateKey(bdState.State.InternalName)
				// Remove BD state
				if bdState.State.Index == 0 && bdState.State.InternalName != "" {
					plugin.Transport.PublishData(key, nil)
					log.Debugf("Bridge domain %v: state removed from ETCD", bdState.State.InternalName)
					// Write/Update BD state
				} else if bdState.State.Index != 0 {
					plugin.Transport.PublishData(key, bdState.State)
					log.Debugf("Bridge domain %v: state stored in ETCD", bdState.State.InternalName)
				} else {
					log.Warnf("Unable to process bridge domain state with Idx %v and Name %v",
						bdState.State.Index, bdState.State.InternalName)
				}
			}
		case <-ctx.Done():
			// Stop watching for state data updates
			return
		}
	}
}
