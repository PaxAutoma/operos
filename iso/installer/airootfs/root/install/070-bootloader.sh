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

echo \> Installing bootloader >&3

install_bios() {
    # install MBR
    dd bs=440 conv=notrunc count=1 if=/usr/lib/syslinux/bios/gptmbr.bin of=${CONTROLLER_DISK}

    cat > /mnt/boot/syslinux.cfg <<EOF
PATH /syslinux-controller/
DEFAULT loadconfig

LABEL loadconfig
  CONFIG /syslinux-controller/syslinux.cfg
EOF

    # set up syslinux
    arch-chroot /mnt extlinux --install /boot
    sed -i "s/%CONTROLLER_DISK%/${CONTROLLER_DISK//\//\\\/}/g;
            s/%OPEROS_VERSION%/${OPEROS_VERSION}/g;" /mnt/boot/syslinux-controller/*.cfg
}

install_efi() {
    mkdir -p /mnt/efi/EFI/operos-${OPEROS_VERSION}
    cp /run/archiso/bootmnt/operos/boot/intel_ucode.img /mnt/efi/EFI/operos-${OPEROS_VERSION}/intel_ucode.img
    cp /run/archiso/bootmnt/operos/boot/intel_ucode.LICENSE /mnt/efi/EFI/operos-${OPEROS_VERSION}/intel_ucode.LICENSE
    cp /run/archiso/bootmnt/operos/boot/x86_64/archiso.img /mnt/efi/EFI/operos-${OPEROS_VERSION}/archiso.img
    cp /run/archiso/bootmnt/operos/boot/x86_64/vmlinuz /mnt/efi/EFI/operos-${OPEROS_VERSION}/vmlinuz.efi

    # copy loader config
    mkdir /mnt/efi/loader /mnt/efi/loader/entries
    cp efiboot/loader/loader.conf /mnt/efi/loader/loader.conf
    cp efiboot/loader/entries/operos.conf /mnt/efi/loader/entries/operos-${OPEROS_VERSION}.conf 

    find /mnt/efi/loader -name "*.conf" -exec \
        sed -i "s/%CONTROLLER_DISK%/${CONTROLLER_DISK//\//\\\/}/g;
                s/%OPEROS_VERSION%/${OPEROS_VERSION}/g;" {} \;

    # install the systemd-boot efi binaries
    bootctl --path=/mnt/efi install
}

if [ -d /sys/firmware/efi/efivars ]; then
    install_efi
else
    install_bios
fi
