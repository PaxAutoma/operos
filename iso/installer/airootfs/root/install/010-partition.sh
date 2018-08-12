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

set -x

echo \> Partitioning disk >&3

. install/_common || exit 1

# Disable LVM volumes to let go of disk devices
vgchange -an

# Clear out any existing partition tables
sgdisk -Z $CONTROLLER_DISK
ret=$?
# sgdisk -Z returns 2 if the partition tables were broken even if it
# successfully cleared them
if [ $ret -ne 0 ] && [ $ret -ne 2 ]; then
    echo "Failed to clear partition tables" >&2
    exit 1
fi

set -e

# Partition
sgdisk -n 1:1M:+1M -t 1:ef02 -c 1:'BIOS Boot' \
       -n 2:2M:+512M -t 2:ef00 -c 2:'EFI System' \
       -n 3:514M:+1G -t 3:8300 -A 3:set:2 -c 3:'Boot' \
       -n 4:1538M:+4G -t 4:8200 -c 4:'Swap' \
       -n 5:5634M:+10G -t 5:8304 -c 5:'OS' \
       $CONTROLLER_DISK
ENDSECTOR=`sgdisk -E $CONTROLLER_DISK`
sgdisk -n 6:15874M:$ENDSECTOR -t 6:8304 -c 6:'Root overlay' $CONTROLLER_DISK

sleep 3

partprobe -s $CONTROLLER_DISK

# Make file systems
yes | mkfs.fat -F32 $(partition_dev $CONTROLLER_DISK 2)
yes | mkfs.ext4 -O \^64bit $(partition_dev $CONTROLLER_DISK 3)
yes | mkswap $(partition_dev $CONTROLLER_DISK 4)
yes | mkfs.ext4 $(partition_dev $CONTROLLER_DISK 5)
yes | mkfs.ext4 $(partition_dev $CONTROLLER_DISK 6)
