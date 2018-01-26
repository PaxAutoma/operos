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

set -e

# Permissions
chmod 700 /root
chown 0:0 /root

# Networking
systemctl enable nat.service
systemctl enable systemd-networkd.service
systemctl enable systemd-resolved.service

# Images
systemctl enable load-images.service

# etcd
mkdir -p /var/operos/cfg
systemctl enable etcd.service
systemctl enable operos-cfg-store.service
systemctl enable operos-cfg-populate.service

# Operos services
systemctl enable teamster.service
systemctl enable operos-image.service
systemctl enable tftpd.service
sed -i 's/#ShowStatus=.*/ShowStatus=no/' /etc/systemd/system.conf
systemctl disable getty@tty1.service
systemctl enable statustty.service

# Addons to run inside Kubernetes
systemctl enable start-addons.service

# Ceph
systemctl enable operos-ceph-mon-init.service
