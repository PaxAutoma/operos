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

check_partition_name() {
    local disk=$1
    local part_idx=$2
    local label=$3
    local fname=/dev/disk/by-partlabel/${label}

    return $([ -e "$fname" -a "$(readlink -f ${fname})" = "${disk}${part_idx}" ])
}

check_partitions() {
    local disk=$1
    local serial=$2
    check_partition_name $disk 1 "operos-eph-$serial" && \
    check_partition_name $disk 2 "operos-per-$serial"

    return $?
}

repartition() {
    local disk=$1
    local serial=$2
    local part_size=$3

    # clear all partition table signatures from disk
    wipefs -a $disk
    sgdisk -og $disk

    startsector=`sgdisk -F $disk`
    endsector=`sgdisk -E $disk`

    part1_size=$(printf %.0f $(echo "($endsector - $startsector) * (1.0 - 0.$part_size)" | bc))
    sgdisk -n 1:$startsector:+$part1_size -t 1:8300 -c 1:"operos-eph-$serial" $disk

    next_start=`sgdisk -F $disk`
    sgdisk -n 2:$next_start:$endsector -t 1:8300 -c 2:"operos-per-$serial" $disk
}

create_system_directories() {
    local disk=$1

    mkdir -p /tmp/repartition
    mount ${disk} /tmp/repartition
    mkdir -p /tmp/repartition/{docker,kubelet,journal}
    umount /tmp/repartition
    rm -rf /tmp/repartition
}

. /etc/paxautoma/settings

case $1 in
"start")
    declare -a VOLGROUP_MEMBERS
    VOLGROUP_NMEMBERS=0

    for dname in $( lsblk -o KNAME,TYPE,TRAN | grep disk | grep -v usb | cut -d ' ' -f 1) ; do
        disk="/dev/$dname"
        serial=$(lsblk --nodeps -o name,serial | grep ^$dname | awk '{ print $2 }' )
        if [[ -z "$serial" ]]; then
            serial=$dname
        fi

        if ! check_partitions "$disk" "$serial"; then
            echo "Re-partitioning disk $disk $serial"

            set -x
            repartition "$disk" "$serial" "$OPEROS_WORKER_STORAGE_PERCENTAGE"
            result=$?
            set +x

            if [[ $result -ne 0 ]]; then
                echo "Warning: Could not repartition ${disk} -- skipping."
                continue
            fi

            # give the kernel a few seconds to pick up the new partitions
            sleep 3
        else
            echo "Partitions ok, moving on"
        fi

        VOLGROUP_MEMBERS[$VOLGROUP_NMEMBERS]=${disk}1
        VOLGROUP_NMEMBERS=$((VOLGROUP_NMEMBERS+1))

        # osd

        BDEVUUID=$(/usr/bin/prospector --blk-device-uuid ${disk} | grep ^$dname | cut -d ',' -f 2)
        OSD_ID=$(cat /etc/paxautoma/osd-loadout | grep "^${BDEVUUID}" | cut -d ':' -f 2)
        OSD_KEY=$(cat /etc/paxautoma/osd-loadout | grep "${BDEVUUID}" | cut -d ':' -f 3)

        if [ -n "$OSD_ID" ]; then
            ptype=$(blkid ${disk}2 | grep -Po 'TYPE=\"\K[^ \"]+' || /bin/true)
            if [ "$ptype" != "xfs" ] ; then
                mkfs.xfs -f ${disk}2
            fi

            mkdir -p /var/lib/ceph/osd/ceph-${OSD_ID}
            mount ${disk}2 /var/lib/ceph/osd/ceph-${OSD_ID}

            init_osd=0
            if [ ! -f /var/lib/ceph/osd/ceph-${OSD_ID}/fsid ] ; then
                init_osd=1
            elif [ "$(cat /var/lib/ceph/osd/ceph-${OSD_ID}/fsid)" != "$BDEVUUID" ]  ; then
                init_osd=1
            elif [ "$(cat /var/lib/ceph/osd/ceph-${OSD_ID}/ceph_fsid)" != "$(cat /etc/ceph/ceph.conf | grep fsid | cut -d '=' -f 2 | sed -e 's/^ //')" ] ; then
                init_osd=1
            fi

            if [ $init_osd -eq 1 ] ; then
                rm -rf /var/lib/ceph/osd/ceph-${OSD_ID}/*
                printf "[osd.$OSD_ID]\n\tkey = $OSD_KEY\n" > /var/lib/ceph/osd/ceph-${OSD_ID}/keyring
                ceph-osd -i ${OSD_ID} --mkfs --osd-uuid $BDEVUUID
                chown -R ceph:ceph /var/lib/ceph/osd/ceph-${OSD_ID}
            fi

            systemctl enable ceph-osd@${OSD_ID}
            systemctl start ceph-osd@${OSD_ID}
        else
            echo "Warning: Controller did not see fit to add ${disk} ${BDEVUUID} to the ceph cluster"
        fi
    done

    if [ $VOLGROUP_NMEMBERS -eq 0 ] ; then
        echo "No disks to add to system storage found"
        exit 1
    fi

    #volume group
    init_vg=0
    destroy_vg=0
    vols=$(vgs -o vg_name,vg_missing_pv_count --separator ',' --noheadings | sed -e 's/^  //' | grep operos_system_vg0 || /bin/true)

    if [ -z "$vols" ] ; then
        init_vg=1
    elif [ $(echo $vols | cut -d ',' -f 2) != "0" ] ; then
        # we lost a disk
        init_vg=1
        destroy_vg=1
    else
        # check if all the disks we found are part of the volume group
        init_vg=$VOLGROUP_NMEMBERS

        readarray -t vgmembers <<<"$(vgs --separator ',' -o pv_name,vg_name --noheadings | sed -e 's/^  //' | grep operos_system_vg0 | cut -d ',' -f 1)"

        for cvgm in ${VOLGROUP_MEMBERS[*]} ; do 
            for vgm in ${vgmembers[*]} ; do
                if [ "$vgm" = "$cvgm" ] ; then
                    init_vg=$((init_vg-1))
                    break
                fi
            done
        done

        if (( $init_vg  != 0 )) ; then
            destroy_vg=1
        fi
    fi

    if (( $destroy_vg != 0 )) ; then
        readarray -t old_pd <<< "$(vgs --separator ',' -o pv_name,vg_name --noheadings | sed -e 's/^  //' | grep operos_system_vg0 | cut -d ',' -f 1)"
        lvremove -f operos_system_vg0/operos_system_volume 
        vgremove -f operos_system_vg0
        for pd in ${old_pd[*]} ; do
            if [ "$pd" != "[unknown]" ] ; then
                pvremove $pd
            fi
        done

    fi

    if (( $init_vg != 0 )) ; then
        for pd in ${VOLGROUP_MEMBERS[*]} ; do
	    # Linux raid metadata blocks may prevent pvcreate from initializing
	    # the disk with 'Device (/dev/XXXX) excluded by filter)'
	    dd if=/dev/zero bs=1k count=10 of=$pd
            pvcreate -f $pd
        done
        vgcreate operos_system_vg0 ${VOLGROUP_MEMBERS[*]}
        lvcreate -l 100%FREE operos_system_vg0 -n operos_system_volume

        mkfs.ext4 -F -F /dev/mapper/operos_system_vg0-operos_system_volume
        create_system_directories /dev/mapper/operos_system_vg0-operos_system_volume
    fi

    mkdir -p /storage/system /var/lib/{docker,kubelet,journal}
    mount /dev/mapper/operos_system_vg0-operos_system_volume /storage/system
    mount --bind /storage/system/docker /var/lib/docker
    mount --bind /storage/system/kubelet /var/lib/kubelet
    mount --bind /storage/system/journal /var/lib/journal

    ;;

"stop")
    umount /var/lib/kubelet
    umount /var/lib/docker
    umount /var/lib/journal
    umount /storage/system
    ;;

*)
    echo "Invalid command: $1"
    exit 1
    ;;
esac
