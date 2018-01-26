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

ln -sf /usr/share/zoneinfo/UTC /etc/localtime
ln -s /dev/null /etc/systemd/network/99-default.link

#usermod -s /usr/bin/zsh root
#cp -aT /etc/skel/ /root/
chmod 700 /root
chown 0:0 /root

sed -i 's/#\(PermitRootLogin \).\+/\1yes/' /etc/ssh/sshd_config
sed -i "s/#Server/Server/g" /etc/pacman.d/mirrorlist
sed -i 's/#\(Storage=\)auto/\1volatile/' /etc/systemd/journald.conf

# Networking
systemctl enable systemd-resolved.service

# Pacman
systemctl enable pacman-init.service

# Operos services
sed -i 's/#ShowStatus=.*/ShowStatus=no/' /etc/systemd/system.conf
systemctl disable getty@tty1.service
systemctl enable statustty.service

# Pax Automa worker setup
systemctl enable make-partitions.service
systemctl enable apply-settings.service
systemctl enable pull-images.service
