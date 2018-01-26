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

script_path=$(readlink -f ${0%/*})

. ${script_path}/versions

# Now, grab the necessary Docker images and save them into the installer
img_path=${script_path}/cache/docker
controller_images="
    gcr.io/google_containers/hyperkube:v${OPEROS_KUBERNETES_VERSION}
    gcr.io/google_containers/pause-amd64:3.0
    gcr.io/google_containers/k8s-dns-kube-dns-amd64:${OPEROS_KUBEDNS_VERSION}
    gcr.io/google_containers/k8s-dns-dnsmasq-nanny-amd64:${OPEROS_KUBEDNS_VERSION}
    gcr.io/google_containers/k8s-dns-sidecar-amd64:${OPEROS_KUBEDNS_VERSION}
    gcr.io/google_containers/kubernetes-dashboard-amd64:v${OPEROS_KUBE_DASHBOARD_VERSION}
    quay.io/calico/node:v${OPEROS_CALICO_VERSION}
    quay.io/prometheus/prometheus:${OPEROS_PROMETHEUS_VERSION}
    prom/node-exporter:${OPEROS_NODE_EXPORTER_VERSION}
    nginx:${OPEROS_NGINX_VERSION}
    quay.io/external_storage/rbd-provisioner:${OPEROS_RBD_PROVISIONER_VERSION}
"

# These images will be sent to the worker. Versions from above will be used.
worker_images="
    quay.io/calico/node:v${OPEROS_CALICO_VERSION}
    prom/node-exporter:${OPEROS_NODE_EXPORTER_VERSION}
    gcr.io/google_containers/pause-amd64:3.0
"

all_images=`
    comm --output-delimiter="" \
        <(for i in $controller_images; do echo $i; done | sort) \
        <(for i in $worker_images; do echo $i; done | sort)
`

pull() {
    local img=$1

    # .io domains seem to intermittently fail to resolve
    # so let's repeat this a few times
    for i in {1..3}; do
        docker pull $img
        [ $? -eq 0 ] && return
        sleep 3
    done

    return 1
}


mkdir -p $img_path
rm -rf $img_path/*
for img in $all_images; do
    (set +e; pull $img)

    mkdir -p $img_path/$(dirname $img)
    docker save $img | gzip > $img_path/$img.tar.gz
done

find $img_path -name "*.tar.gz" | sed s,^$img_path/,, > $img_path/index.txt
rm -f $img_path/worker.txt
for i in $worker_images; do
    echo $i.tar.gz >> $img_path/worker.txt
done
