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

name=installer-diag-$(date +"%Y%m%d-%H%M%S")
tmpdir=$(mktemp -d)
workdir=$tmpdir/$name
trap "rm -rf $tmpdir" EXIT

mkdir -p $workdir

echo "Collecting data"
cp -r /root/logs $workdir/logs
dmesg &> $workdir/dmesg
systemctl &> $workdir/systemctl
lshw -quiet -json &> $workdir/lshw
ifconfig -a &> $workdir/ifconfig
route -n &> $workdir/ifconfig
netstat -an &> $workdir/netstat
df -k &> $workdir/df
ps -aux &> $workdir/ps
journalctl -x &> $workdir/journal

# Failure is no longer an option
set -e

echo "Creating tarball"
cd $tmpdir
tar -jcf /root/$name.tar.bz2 $name

echo "Authorizing upload"
curlcmd=`/usr/bin/operosdiagauth $@ /root/$name.tar.bz2`

echo "Uploading"
eval $curlcmd
