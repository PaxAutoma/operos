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

echo \> Configuring networking >&3

if [ "$CONTROLLER_PUBLIC_IF_MODE" = "dhcp" ]; then
    cat > /mnt/etc/systemd/network/10-public.network <<EOF
[Match]
Name=${CONTROLLER_PUBLIC_IF}

[Network]
DHCP=ipv4
EOF
elif [ "$CONTROLLER_PUBLIC_IF_MODE" = "static" ]; then
    cat > /mnt/etc/systemd/network/10-public.network <<EOF
[Match]
Name=${CONTROLLER_PUBLIC_IF}

[Network]
Address=${CONTROLLER_PUBLIC_IF_IPNET}
Gateway=${CONTROLLER_PUBLIC_IF_GW}
EOF
fi

cat > /mnt/etc/systemd/network/20-private.network <<EOF
[Match]
Name=${CONTROLLER_PRIVATE_IF}

[Network]
Address=${OPEROS_CONTROLLER_IP}${OPEROS_NODE_MASK}

[Route]
Destination=${OPEROS_SERVICE_CIDR}
EOF

if [ "$CONTROLLER_PUBLIC_IF_MODE" = "disabled" ]; then
    cat >> /mnt/etc/systemd/network/20-private.network <<EOF
[Route]
Gateway=${OPEROS_PRIVATE_GW}
EOF
else
    arch-chroot /mnt systemctl enable nat.service
fi

ln -sf /run/systemd/resolve/resolv.conf /mnt/etc/resolv.conf
