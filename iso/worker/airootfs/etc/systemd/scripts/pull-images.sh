#!/bin/bash -x
# Copyright 2018 Pax Automa Systems, Inc.
# 
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# 
#    http://www.apache.org/licenses/LICENSE-2.0
# 
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

. /etc/paxautoma/settings

load_images() {
    images=`curl --fail http://${OPEROS_CONTROLLER_IP}:5080/images/worker.txt`

    for img in $images; do
        curl --fail "http://${OPEROS_CONTROLLER_IP}:5080/images/$img" | docker load
    done
}

while true; do
    (set -e; load_images)
    [ $? -eq 0 ] && break
    sleep 10
done

systemctl disable pull-images.service
