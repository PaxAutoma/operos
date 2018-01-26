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
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func SubscribeHostname(closer <-chan struct{}) <-chan *string {
	ch := make(chan *string)

	go func() {
		ch <- GetHostname()
		ticker := time.NewTicker(3 * time.Second)

		for {
			select {
			case _, ok := <-closer:
				if !ok {
					break
				}
			case <-ticker.C:
				ch <- GetHostname()
			}
		}
	}()

	return ch
}

func GetHostname() *string {
	name, err := os.Hostname()
	if err != nil {
		log.Errorf("cannot get hostname: %v", err)
		return nil
	}
	return &name
}
