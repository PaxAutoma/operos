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
cp /mnt/boot/syslinux-worker/ldlinux.c32 /mnt/etc/paxautoma/iso/

cat > /mnt/etc/paxautoma/iso/operos.ipxe <<EOF
#!ipxe
:retry_dhcp
dhcp && isset \${next-server} || goto retry_dhcp
echo Operos worker booting from http://\${next-server}:5080
set 210:string http://\${next-server}:5080/
set 209:string worker.cfg
chain http://\${next-server}:5080/boot/syslinux-worker/lpxelinux.0
EOF


cat > /mnt/etc/paxautoma/iso/worker.cfg <<EOF
PATH syslinux-worker/
DEFAULT loadconfig

LABEL loadconfig
  CONFIG boot/syslinux-worker/syslinux.cfg
  APPEND boot/
EOF

# tftp
cat > /mnt/etc/conf.d/tftpd <<EOF
TFTPD_ARGS="--verbose --address ${OPEROS_CONTROLLER_IP} -m /tftpboot/mapfile -u ftp --secure /tftpboot"
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

mkdir -p /mnt/etc/dhcpd

cat > /mnt/etc/dhcpd.conf <<EOF
ddns-update-style none;
deny bootp;     #default
authoritative;

include "/etc/dhcpd/ipxe-option-space.conf";

subnet ${CLUSTER_SN_START} netmask ${CLUSTER_SN_MASK} {
    range ${CLUSTER_NODE_START} ${CLUSTER_NODE_END};
    option subnet-mask ${CLUSTER_SN_MASK};
    option routers ${OPEROS_CONTROLLER_IP};

    option domain-name "${OPEROS_DNS_DOMAIN}";
    option domain-search "${OPEROS_DNS_DOMAIN}";
    option domain-name-servers ${OPEROS_CONTROLLER_IP};

    include "/etc/dhcpd/ipxe-bootp.conf";
}

use-host-decl-names on;
EOF

cat > /mnt/etc/dhcpd/ipxe-option-space.conf <<EOF
# Declare the iPXE/gPXE/Etherboot option space
option space ipxe;
option ipxe-encap-opts code 175 = encapsulate ipxe;

# iPXE options, can be set in DHCP response packet
option ipxe.priority         code   1 = signed integer 8;
option ipxe.keep-san         code   8 = unsigned integer 8;
option ipxe.skip-san-boot    code   9 = unsigned integer 8;
option ipxe.syslogs          code  85 = string;
option ipxe.cert             code  91 = string;
option ipxe.privkey          code  92 = string;
option ipxe.crosscert        code  93 = string;
option ipxe.no-pxedhcp       code 176 = unsigned integer 8;
option ipxe.bus-id           code 177 = string;
option ipxe.bios-drive       code 189 = unsigned integer 8;
option ipxe.username         code 190 = string;
option ipxe.password         code 191 = string;
option ipxe.reverse-username code 192 = string;
option ipxe.reverse-password code 193 = string;
option ipxe.version          code 235 = string;
option iscsi-initiator-iqn   code 203 = string;

# iPXE feature flags, set in DHCP request packet
option ipxe.pxeext    code 16 = unsigned integer 8;
option ipxe.iscsi     code 17 = unsigned integer 8;
option ipxe.aoe       code 18 = unsigned integer 8;
option ipxe.http      code 19 = unsigned integer 8;
option ipxe.https     code 20 = unsigned integer 8;
option ipxe.tftp      code 21 = unsigned integer 8;
option ipxe.ftp       code 22 = unsigned integer 8;
option ipxe.dns       code 23 = unsigned integer 8;
option ipxe.bzimage   code 24 = unsigned integer 8;
option ipxe.multiboot code 25 = unsigned integer 8;
option ipxe.slam      code 26 = unsigned integer 8;
option ipxe.srp       code 27 = unsigned integer 8;
option ipxe.nbi       code 32 = unsigned integer 8;
option ipxe.pxe       code 33 = unsigned integer 8;
option ipxe.elf       code 34 = unsigned integer 8;
option ipxe.comboot   code 35 = unsigned integer 8;
option ipxe.efi       code 36 = unsigned integer 8;
option ipxe.fcoe      code 37 = unsigned integer 8;
option ipxe.vlan      code 38 = unsigned integer 8;
option ipxe.menu      code 39 = unsigned integer 8;
option ipxe.sdi       code 40 = unsigned integer 8;
option ipxe.nfs       code 41 = unsigned integer 8;

# Other useful general options
# http://www.ietf.org/assignments/dhcpv6-parameters/dhcpv6-parameters.txt
option arch code 93 = unsigned integer 16;
EOF

cat > /mnt/etc/dhcpd/ipxe-bootp.conf <<EOF
allow bootp;
allow booting;
next-server ${OPEROS_CONTROLLER_IP};

# Disable ProxyDHCP, we're in control of the primary DHCP server
#option ipxe.no-pxedhcp 1;

# Make sure the iPXE we're loading supports what we need,
# if not load a full-featured version
if    exists ipxe.http
  and exists ipxe.menu
  and ( ( exists ipxe.pxe
      and exists ipxe.bzimage
      and exists ipxe.elf
  ) or (
      exists ipxe.efi
  ) ) {
    filename "http://${OPEROS_CONTROLLER_IP}:5080/operos.ipxe";
}
elsif exists user-class and option user-class = "iPXE" {
    if option arch = 00:06 {
        filename "ipxe-x86.efi";
    } elsif option arch = 00:07 {
        filename "ipxe-x64.efi";
    } elsif option arch = 00:00 {
        filename "ipxe.pxe";
    }
}
elsif exists user-class and option user-class = "gPXE" {
    filename "ipxe.pxe";
}
elsif option arch = 00:06 {
    # EFI 32-bit
    filename "ipxe-x86.efi";
}
elsif option arch = 00:07 {
    # EFI 64-bit
    filename "ipxe-x64.efi";
}
elsif option arch = 00:00 {
    # Legacy BIOS x86 mode
    filename "ipxe.pxe";
}
else {
    # Unsupported client architecture type, so do nothing
}
EOF


arch-chroot /mnt systemctl enable dhcpd4@${CONTROLLER_PRIVATE_IF}.service
