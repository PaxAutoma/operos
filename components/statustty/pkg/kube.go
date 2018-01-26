/*
Copyright 2018 Pax Automa Systems, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package statustty

import (
	"net/http"
	"time"

	"github.com/paxautoma/operos/components/common"

	log "github.com/sirupsen/logrus"
)

type KubeStatus struct {
	Reachable bool
}

func SubscribeKubeStatus(kubeURL string, closer <-chan struct{}) <-chan *KubeStatus {
	ch := make(chan *KubeStatus)

	if kubeURL != "" {
		go func() {
			defer common.LogPanic()

			ch <- GetKubeStatus(kubeURL)

			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case _, ok := <-closer:
					if !ok {
						break
					}
				case <-ticker.C:
					ch <- GetKubeStatus(kubeURL)
				}

			}
		}()
	}

	return ch
}

func GetKubeStatus(kubeURL string) *KubeStatus {
	resp, err := http.Get(kubeURL)
	if err != nil {
		log.Debugf("could not reach kubernetes: %v", err)
		return &KubeStatus{Reachable: false}
	}
	defer resp.Body.Close()

	return &KubeStatus{Reachable: true}
}
