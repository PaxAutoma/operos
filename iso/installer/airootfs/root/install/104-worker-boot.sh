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

echo \> Setting up worker boot infrastructure >&3

mkdir -p /mnt/etc/paxautoma/iso
cp /mnt/boot/syslinux-worker/lpxelinux.0 /mnt/boot/
cp /mnt/boot/syslinux-worker/ldlinux.c32 /mnt/boot/
mkdir -p /mnt/boot/pxelinux.cfg

cat > /mnt/boot/pxelinux.cfg/default <<EOF
PATH /syslinux-worker/
DEFAULT loadconfig

LABEL loadconfig
  CONFIG /syslinux-worker/syslinux.cfg
EOF

# tftp
cat > /mnt/etc/conf.d/tftpd <<EOF
TFTPD_ARGS="--verbose --address ${OPEROS_CONTROLLER_IP} -m /etc/tftpd.mapfile -u ftp --secure /boot"
EOF

# dhcp

function increment_ipv4 {
  local ip_base=$1
  local incr_amount=$2
  local -a ip_components
  local ip_regex="([0-9]+).([0-9]+).([0-9]+).([0-9]+)"
  [[ $ip_base =~ $ip_regex ]]
  ip_components=("${BASH_REMATCH[1]}" "${BASH_REMATCH[2]}" "${BASH_REMATCH[3]}" "${BASH_REMATCH[4]}")
  ip_dec=0
  local comp
  for comp in "${ip_components[@]}"; do
    ip_dec=$((ip_dec<<8))
    ip_dec=$((ip_dec + $comp))
  done

  ip_dec=$((ip_dec + $incr_amount))

  ip_components=()
  local i
  for ((i=0; i < 4; i++)); do
    comp=$((ip_dec & 0xFF))
    ip_components+=($comp)
    ip_dec=$((ip_dec>>8))
  done
  echo "${ip_components[3]}.${ip_components[2]}.${ip_components[1]}.${ip_components[0]}"
}

ip2int()
{
    local a b c d
    { IFS=. read a b c d; } <<< $1
    echo $(((((((a << 8) | b) << 8) | c) << 8) | d))
}

network()
{
    local addr=$(ip2int $1); shift
    local mask=$((0xffffffff << (32 -$1))); shift
    int2ip $((addr & mask))
}

netmask()
{
    local mask=$((0xffffffff << (32 - $1))); shift
    int2ip $mask
}

broadcast()
{
    local addr=$(ip2int $1); shift
    local mask=$((0xffffffff << (32 -$1))); shift
    int2ip $((addr | ~mask))
}

end_ip()
{
    local addr=$(ip2int $1); shift
    local mask=$((0xffffffff << (32 -$1))); shift
    int2ip $((addr | ~mask<<1))
}

int2ip()
{
    local ui32=$1; shift
    local ip n
    for n in 1 2 3 4; do
        ip=$((ui32 & 0xff))${ip:+.}$ip
        ui32=$((ui32 >> 8))
    done
    echo $ip
}

CLUSTER_SN_MASK=$(netmask ${OPEROS_NODE_MASK//\/})

CLUSTER_NODE_START=$(increment_ipv4 $OPEROS_CONTROLLER_IP 10)
CLUSTER_SN_START=$(network $OPEROS_CONTROLLER_IP ${OPEROS_NODE_MASK//\/})

CLUSTER_NODE_END=$(end_ip $CLUSTER_NODE_START ${OPEROS_NODE_MASK//\/})

cat > /mnt/etc/dhcpd.conf <<EOF
ddns-update-style none;
deny bootp;     #default
authoritative;

subnet ${CLUSTER_SN_START} netmask ${CLUSTER_SN_MASK} {
    allow bootp;
    allow booting;
    range ${CLUSTER_NODE_START} ${CLUSTER_NODE_END};
    option subnet-mask ${CLUSTER_SN_MASK};
    option routers ${OPEROS_CONTROLLER_IP};

    option domain-name "${OPEROS_DNS_DOMAIN}";
    option domain-search "${OPEROS_DNS_DOMAIN}";
    option domain-name-servers ${OPEROS_CONTROLLER_IP};

    next-server $OPEROS_CONTROLLER_IP;
    filename "lpxelinux.0";
}

use-host-decl-names on;
EOF

arch-chroot /mnt systemctl enable dhcpd4@${CONTROLLER_PRIVATE_IF}.service
