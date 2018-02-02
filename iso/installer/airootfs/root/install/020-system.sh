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

echo \> Copying system files >&3

# Mount
mkdir -p /run/archiso/sfs/airootfs-node
mount -o loop /run/archiso/bootmnt/operos/x86_64/airootfs-node.sfs /run/archiso/sfs/airootfs-node

mkdir -p /run/archiso/sfs/airootfs-controller
mount -o loop /run/archiso/bootmnt/operos/x86_64/airootfs-controller.sfs /run/archiso/sfs/airootfs-controller

mkdir -p /run/archiso/cowspace/controller
mount ${CONTROLLER_DISK}6 /run/archiso/cowspace/controller/
mkdir /run/archiso/cowspace/controller/{upperdir,workdir}

mount -t overlay overlay \
    -o lowerdir=/run/archiso/sfs/airootfs-controller:/run/archiso/sfs/airootfs-node:/run/archiso/sfs/airootfs-base,upperdir=/run/archiso/cowspace/controller/upperdir,workdir=/run/archiso/cowspace/controller/workdir \
    /mnt

mkdir -p /mnt/efi /mnt/boot /mnt/run/archiso/bootmnt
mount ${CONTROLLER_DISK}2 /mnt/efi
mount ${CONTROLLER_DISK}3 /mnt/boot

# Kube recommends running without swap. We'll keep the partition around just in
# case they decide to enable this later, but will not enable it in fstab for
# now.
#swapon ${CONTROLLER_DISK}4

mount ${CONTROLLER_DISK}5 /mnt/run/archiso/bootmnt

# Copy files
mkdir -p /mnt/run/archiso/bootmnt/operos-${OPEROS_VERSION}/
cp -af /run/archiso/bootmnt/operos/boot/* /mnt/boot/
cp -af /run/archiso/bootmnt/operos/x86_64 /mnt/run/archiso/bootmnt/operos-${OPEROS_VERSION}/

# Unmount the directory with the SquashFS images to prevent it from being added
# to fstab later. It will be mounted by the cpio hooks (as archisodevice).
umount /mnt/run/archiso/bootmnt
