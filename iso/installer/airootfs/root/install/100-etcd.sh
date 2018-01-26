#!/bin/bash -xe
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

echo \> Setting up the Operos configuration store >&3

cat > /mnt/etc/conf.d/etcd <<EOF
ETCD_ADVERTISE_CLIENT_URLS="http://${OPEROS_CONTROLLER_IP}:2379"
ETCD_INITIAL_ADVERTISE_PEER_URLS="http://${OPEROS_CONTROLLER_IP}:2380"
ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:2379
ETCD_LISTEN_PEER_URLS=http://0.0.0.0:2380
ETCD_INITIAL_CLUSTER="$(cat /mnt/etc/machine-id)=http://${OPEROS_CONTROLLER_IP}:2380"
EOF

