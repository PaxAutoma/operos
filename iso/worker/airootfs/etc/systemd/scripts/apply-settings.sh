#!/bin/bash
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

set -xue

get_boot_if() {
    local BOOTIF=$(sed "s/.*BOOTIF=[\'\"]\?\([^[:space:]\'\"]\+\).*/\1/;tx;d;:x" /proc/cmdline)
    local bootif_mac
    local bootif_dev

    if [ -n "$BOOTIF" ]; then
        bootif_mac=${BOOTIF#01-}
        bootif_mac=${bootif_mac//-/:}
        for i in /sys/class/net/*/address; do
            read net_mac < ${i}
            if [[ "${bootif_mac}" == "${net_mac}" ]]; then
                bootif_dev=${i#/sys/class/net/}
                bootif_dev=${bootif_dev%/address}
                break
            fi
        done
    fi

    echo ${bootif_dev:-eth0}
}

get_boot_server() {
    local bootif_dev=${1:-eth0}
    local boot_server=$(dhcpcd -U ${bootif_dev} | grep dhcp_server | cut -d '=' -f 2)

    if [[ -z "$boot_server" ]] ; then
        echo "ERROR: could got get operos server from dhcpcd" >&2
        return 1
    fi

    echo $boot_server
}

with_backoff() {
    # from https://stackoverflow.com/questions/8350942/how-to-re-run-the-curl-command-automatically-when-the-error-occurs/8351489#8351489
    local max_attempts=${ATTEMPTS-0}
    local timeout=${TIMEOUT-1}
    local max_timeout=${MAXTIMEOUT-30}
    local attempt=0
    local exit_code=0

    while (( $max_attempts == 0 || $attempt < $max_attempts ))
    do
        "$@"
        exit_code=$?

        if [[ $exit_code == 0 ]]
        then
            break
        fi

        echo "Failure! Retrying in $timeout.." 1>&2
        sleep $timeout
        attempt=$(( attempt + 1 ))
        timeout=$(( timeout * 2 <= max_timeout ? timeout * 2 : max_timeout ))
    done

    if [[ $exit_code != 0 ]]
    then
        echo "Giving up ($@)" 1>&2
    fi

    return $exit_code
}

download_settings() {
    /usr/bin/prospector | curl -X POST -d @- http://${boot_server}:2680/whoami | tar -C / -zxv
}

boot_if=$(set +e; get_boot_if)
boot_server=$(set +e; get_boot_server $boot_if)
(set +e; with_backoff download_settings)

# Set root password
(set -x; echo "root:`cat /etc/rootpasshash`" | chpasswd -e)
rm /etc/rootpasshash

. /etc/paxautoma/settings

hostnamectl set-hostname $(cat /etc/hostname)

cat > /etc/kubernetes/config <<EOF
KUBE_ALLOW_PRIV="--allow-privileged=true"
KUBE_MASTER="--master=https://${OPEROS_CONTROLLER_IP}:${OPEROS_KUBE_API_SECURE_PORT}"
EOF

cat > /etc/kubernetes/kubelet <<EOF
KUBELET_ARGS="--container-runtime=docker \\
    --cluster-dns=${OPEROS_DNS_SERVICE_IP} \\
    --cluster-domain=${OPEROS_DNS_DOMAIN} \\
    --register-node=true \\
    --kubeconfig=/etc/kubernetes/worker-kubeconfig.yml \\
    --tls-cert-file=/etc/kubernetes/ssl/worker.pem \\
    --tls-private-key-file=/etc/kubernetes/ssl/worker-key.pem \\
    --network-plugin=cni"
EOF

cat > /etc/kubernetes/worker-kubeconfig.yml <<EOF
apiVersion: v1
kind: Config
clusters:
- name: local
  cluster:
    server: "https://${OPEROS_CONTROLLER_IP}:${OPEROS_KUBE_API_SECURE_PORT}"
    certificate-authority: /etc/kubernetes/ssl/ca.pem
users:
- name: kubelet
  user:
    client-certificate: /etc/kubernetes/ssl/worker.pem
    client-key: /etc/kubernetes/ssl/worker-key.pem
contexts:
- context:
    cluster: local
    user: kubelet
  name: kubelet-context
current-context: kubelet-context
EOF

cat > /etc/kubernetes/proxy <<EOF
KUBE_PROXY_ARGS="--kubeconfig=/etc/kubernetes/worker-kubeconfig.yml"
EOF

mkdir -p /etc/cni/net.d
cat > /etc/cni/net.d/10-calico.conflist <<EOF
{
    "name": "k8s-pod-network",
    "cniVersion": "0.3.0",
    "plugins": [
        {
            "type": "calico",
            "log_level": "info",
            "datastore_type": "kubernetes",
            "kubernetes": {
                "k8s_api_root": "https://${OPEROS_CONTROLLER_IP}:${OPEROS_KUBE_API_SECURE_PORT}",
                "kubeconfig": "/etc/kubernetes/worker-kubeconfig.yml"
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

cat > /etc/chrony.conf <<EOF
server ${OPEROS_CONTROLLER_IP} iburst
rtconutc
rtcsync
EOF

# make sure that the service IPs are routed over the internal network
ip route add ${OPEROS_SERVICE_CIDR} dev ${boot_if}
