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

echo \> Setting up container networking >&3

mkdir -p /mnt/etc/cni/net.d
cat > /mnt/etc/cni/net.d/10-calico.conflist <<EOF
{
    "name": "k8s-pod-network",
    "cniVersion": "0.3.0",
    "plugins": [
        {
            "type": "calico",
            "log_level": "info",
            "datastore_type": "kubernetes",
            "kubernetes": {
                "k8s_api_root": "http://127.0.0.1:${OPEROS_KUBE_API_INSECURE_PORT}"
            },
            "policy": {
                "type": "k8s"
            },
            "ipam": {
                "type": "host-local",
                "subnet": "usePodCidr"
            }
        },
        {
            "type": "portmap",
            "capabilities": {"portMappings": true},
            "snat": true
        }
    ]
}
EOF
