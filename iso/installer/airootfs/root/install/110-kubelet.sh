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

echo \> Setting up Kubernetes >&3

install_manifests() {
    local src_path=$1
    local tgt_path=$2

    mkdir -p $src_path
    for manifest in $src_path/*; do
        envsubst < $manifest > $tgt_path/$(basename $manifest)
    done
}

mkdir -p /mnt/etc/kubernetes/{manifests,addons}
install_manifests /root/manifests/kubelet /mnt/etc/kubernetes/manifests
install_manifests /root/manifests/addons /mnt/etc/kubernetes/addons

cat > /mnt/etc/kubernetes/config <<EOF
KUBE_ALLOW_PRIV="--allow-privileged=true"
KUBE_MASTER="--master=http://127.0.0.1:8080"
EOF

cat > /mnt/etc/kubernetes/kubelet <<EOF
KUBELET_ARGS="--container-runtime=docker \\
    --kubeconfig=/etc/kubernetes/kubeconfig.yml \\
    --cluster-dns=${OPEROS_DNS_SERVICE_IP} \\
    --cluster-domain=${OPEROS_DNS_DOMAIN} \\
    --pod-manifest-path=/etc/kubernetes/manifests \\
    --node-labels=node-role.kubernetes.io/master=true \\
    --register-with-taints=node-role.kubernetes.io/master=true:NoSchedule \\
    --network-plugin=cni"
EOF

cat > /mnt/etc/kubernetes/kubeconfig.yml <<EOF
apiVersion: v1
kind: Config
clusters:
- name: local
  cluster:
    server: http://127.0.0.1:8080
contexts:
- context:
    cluster: local
  name: kubelet-context
current-context: kubelet-context
EOF
