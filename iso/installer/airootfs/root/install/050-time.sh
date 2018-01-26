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

echo \> Setting up NTP >&3

arch-chroot /mnt ln -sf /usr/share/zoneinfo/UTC /etc/localtime

cat > /mnt/etc/chrony.conf <<EOF
server 0.pool.ntp.org iburst
server 1.pool.ntp.org iburst
server 2.pool.ntp.org iburst
server 3.pool.ntp.org iburst
local

rtconutc
rtcsync

allow ${OPEROS_CONTROLLER_IP}${OPEROS_NODE_MASK}
deny
EOF

mkdir -p /mnt/etc/systemd/system/chronyd.service.d
cat > /mnt/etc/systemd/system/chronyd.service.d/10-wait-sync.conf <<EOF
[Unit]
Before=kubelet.service
After=network-online.target

[Service]
ExecStartPre=/usr/bin/chronyd -q -u chrony
TimeoutStartSec=infinity
EOF
