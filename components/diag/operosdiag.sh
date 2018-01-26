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

name=operos-diag-$(date +"%Y%m%d-%H%M%S")-$(hostname)
tmpdir=$(mktemp -d)
workdir=$tmpdir/$name
trap "rm -rf $tmpdir" EXIT

. /etc/paxautoma/settings

collect() {
    set -x

    echo " > System information" >&3
    hostname &> $workdir/hostname
    dmesg &> $workdir/dmesg
    systemctl &> $workdir/systemctl
    lshw -quiet -json &> $workdir/lshw
    ifconfig -a &> $workdir/ifconfig
    route -n &> $workdir/ifconfig
    netstat -an &> $workdir/netstat
    df -k &> $workdir/df
    ps -aux &> $workdir/ps
    cp /etc/paxautoma/settings $workdir/paxautoma-settings

    echo " > System logs" >&3
    journalctl -x &> $workdir/journal

    echo " > Docker image and container list" >&3
    docker images &> $workdir/docker-images
    docker ps -a &> $workdir/docker-ps

    if [ -e /etc/kubernetes/worker-kubeconfig.yml ]; then
        collect_kube --kubeconfig=/etc/kubernetes/worker-kubeconfig.yml -s "https://${OPEROS_CONTROLLER_IP}:8443"
    else
        collect_kube
    fi
}

collect_kube() {
    kubeargs=$@

    echo " > Kubernetes node and Operos pod information" >&3
    kubectl $kubeargs cluster-info dump --namespaces kube-system,operos --output-directory=$workdir/kube-cluster-info
}

# Descriptor 3 always goes to stdout even if 1 & 2 are being redirected
exec 3>&1

echo "Collecting data"
mkdir -p $workdir

# Run in subshell to avoid set -x bleeding out
(collect &> $workdir/collection-log)

# Failure no longer an option
set -e

echo "Creating tarball"
cd $tmpdir
tar -jcf /root/$name.tar.bz2 $name

echo "Authorizing upload"
curlcmd=`operosdiagauth $@ --install=$OPEROS_INSTALL_ID /root/$name.tar.bz2`

echo "Uploading"
eval $curlcmd
