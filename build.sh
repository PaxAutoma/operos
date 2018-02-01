#!/bin/bash

set -e -u -x

iso_name=operos
iso_label="OPEROS_$(date +%Y%m)"
iso_version=$(date +%Y.%m.%d)
install_dir=operos
arch=$(uname -m)
work_dir=work
out_dir=out
sfs_comp="xz"

verbose=""
devbuild=""
script_path=$(readlink -f ${0%/*})

_usage ()
{
    echo "usage ${0} [options]"
    echo
    echo " General options:"
    echo "    -N <iso_name>      Set an iso filename (prefix)"
    echo "                        Default: ${iso_name}"
    echo "    -V <iso_version>   Set an iso version (in filename)"
    echo "                        Default: ${iso_version}"
    echo "    -L <iso_label>     Set an iso label (disk label)"
    echo "                        Default: ${iso_label}"
    echo "    -D <install_dir>   Set an install_dir (directory inside iso)"
    echo "                        Default: ${install_dir}"
    echo "    -w <work_dir>      Set the working directory"
    echo "                        Default: ${work_dir}"
    echo "    -o <out_dir>       Set the output directory"
    echo "                        Default: ${out_dir}"
    echo "    -d                 Development build (includes test ssh key)"
    echo "    -v                 Enable verbose output"
    echo "    -h                 This help message"
    exit ${1}
}

_join() {
    local IFS="$1"
    shift
    echo "$*"
}

# Helper function to run make_*() only one time per layer/architecture.
run_once() {
    if [[ ! -e ${work_dir}/build.${1}_${layer}_${arch} ]]; then
        "$@"
        touch ${work_dir}/build.${1}_${layer}_${arch}
    fi
}

# Additional packages (airootfs)
make_packages() {
    mkarchiso ${verbose} -w "${work_dir}" -C "${work_dir}/pacman.conf" -D "${install_dir}" -p "$(grep -h -v ^# ${script_path}/iso/${layer}/packages)" install
}

# Copy mkinitcpio archiso hooks and build initramfs (airootfs)
make_setup_mkinitcpio() {
    mkdir -p ${work_dir}/airootfs/etc/initcpio/hooks
    mkdir -p ${work_dir}/airootfs/etc/initcpio/install
    for _hook in operosiso operosiso_pxe_common operosiso_pxe_http; do
        cp ${script_path}/iso/_cpiohooks/hooks/${_hook} ${work_dir}/airootfs/etc/initcpio/hooks
        cp ${script_path}/iso/_cpiohooks/install/${_hook} ${work_dir}/airootfs/etc/initcpio/install
    done
    cp ${script_path}/iso/${layer}/mkinitcpio.conf ${work_dir}/airootfs/etc/mkinitcpio-archiso.conf
    mkarchiso ${verbose} -w "${work_dir}" -C "${work_dir}/pacman.conf" -D "${install_dir}" -r 'mkinitcpio -c /etc/mkinitcpio-archiso.conf -k /boot/vmlinuz-linux -g /boot/archiso.img' run
}

make_customize_airootfs() {
    cp -af --no-preserve=ownership ${script_path}/iso/${layer}/airootfs ${work_dir}

    # repair permissions that may have been damaged by the above copy that may
    # generate warnings from pacman in subsequent layers
    for d in ${work_dir}/airootfs/usr ${work_dir}/airootfs/usr/bin; do
        [[ -d ${d} ]] && chmod 755 ${d}
    done

    mkdir -p ${work_dir}/airootfs/etc/paxautoma

    if [ -x ${work_dir}/airootfs/root/customize_airootfs.sh ]; then
        mkarchiso ${verbose} -w "${work_dir}" -C "${work_dir}/pacman.conf" -D "${install_dir}" -r '/root/customize_airootfs.sh' run
        rm ${work_dir}/airootfs/root/customize_airootfs.sh
    fi

    ln -sf /run/systemd/resolve/resolv.conf ${work_dir}/airootfs/etc/resolv.conf
}

# Prepare ${install_dir}/boot/
make_boot() {
    mkdir -p ${work_dir}/${layer}/iso/${install_dir}/boot/${arch}
    cp ${work_dir}/airootfs/boot/archiso.img ${work_dir}/${layer}/iso/${install_dir}/boot/${arch}/archiso.img
    cp ${work_dir}/airootfs/boot/vmlinuz-linux ${work_dir}/${layer}/iso/${install_dir}/boot/${arch}/vmlinuz

    cp ${work_dir}/airootfs/boot/memtest86+/memtest.bin ${work_dir}/${layer}/iso/${install_dir}/boot/memtest
    cp ${work_dir}/airootfs/usr/share/licenses/common/GPL2/license.txt ${work_dir}/${layer}/iso/${install_dir}/boot/memtest.COPYING
    cp ${work_dir}/airootfs/boot/intel-ucode.img ${work_dir}/${layer}/iso/${install_dir}/boot/intel_ucode.img
    cp ${work_dir}/airootfs/usr/share/licenses/intel-ucode/LICENSE ${work_dir}/${layer}/iso/${install_dir}/boot/intel_ucode.LICENSE
}

# Prepare /${install_dir}/boot/syslinux
make_syslinux() {
    local archiso_layers=$@
    mkdir -p ${work_dir}/${layer}/iso/${install_dir}/boot/syslinux-${layer}
    for _cfg in ${script_path}/iso/${layer}/syslinux/*.cfg; do
        sed "s|%ARCHISO_LABEL%|${iso_label}|g;
             s|%INSTALL_DIR%|${install_dir}|g;
             s|%ARCHISO_LAYER%|${layer}|g;
             s|%ARCHISO_LAYERS%|${archiso_layers}|g;
             s|%OPEROS_VERSION%|${iso_version}|g;" \
             ${_cfg} > ${work_dir}/${layer}/iso/${install_dir}/boot/syslinux-${layer}/${_cfg##*/}
    done
    cp ${script_path}/iso/${layer}/syslinux/splash.png ${work_dir}/${layer}/iso/${install_dir}/boot/syslinux-${layer}/
    cp ${work_dir}/airootfs/usr/lib/syslinux/bios/*.c32 ${work_dir}/${layer}/iso/${install_dir}/boot/syslinux-${layer}/
    cp ${work_dir}/airootfs/usr/lib/syslinux/bios/lpxelinux.0 ${work_dir}/${layer}/iso/${install_dir}/boot/syslinux-${layer}/
    cp ${work_dir}/airootfs/usr/lib/syslinux/bios/memdisk ${work_dir}/${layer}/iso/${install_dir}/boot/syslinux-${layer}/

    mkdir -p ${work_dir}/${layer}/iso/${install_dir}/boot/syslinux-${layer}/hdt
    gzip -c -9 ${work_dir}/airootfs/usr/share/hwdata/pci.ids > ${work_dir}/${layer}/iso/${install_dir}/boot/syslinux-${layer}/hdt/pciids.gz
    gzip -c -9 ${work_dir}/airootfs/usr/lib/modules/*-ARCH/modules.alias > ${work_dir}/${layer}/iso/${install_dir}/boot/syslinux-${layer}/hdt/modalias.gz
}

# Prepare /isolinux
make_isolinux() {
    mkdir -p ${work_dir}/${layer}/iso/isolinux
    sed "s|%INSTALL_DIR%|${install_dir}|g;
         s|%ARCHISO_LAYER%|${layer}|g;" ${script_path}/iso/_isolinux/isolinux.cfg > ${work_dir}/${layer}/iso/isolinux/isolinux.cfg
    cp ${work_dir}/airootfs/usr/lib/syslinux/bios/isolinux.bin ${work_dir}/${layer}/iso/isolinux/
    cp ${work_dir}/airootfs/usr/lib/syslinux/bios/isohdpfx.bin ${work_dir}/${layer}/iso/isolinux/
    cp ${work_dir}/airootfs/usr/lib/syslinux/bios/ldlinux.c32 ${work_dir}/${layer}/iso/isolinux/
}

# Prepare /EFI
make_efi() {
    local archiso_layers=$@
    mkdir -p ${work_dir}/${layer}/iso/EFI/boot
    cp ${work_dir}/airootfs/usr/share/efitools/efi/PreLoader.efi ${work_dir}/${layer}/iso/EFI/boot/bootx64.efi
    cp ${work_dir}/airootfs/usr/share/efitools/efi/HashTool.efi ${work_dir}/${layer}/iso/EFI/boot/

    cp ${work_dir}/airootfs/usr/lib/systemd/boot/efi/systemd-bootx64.efi ${work_dir}/${layer}/iso/EFI/boot/loader.efi

    mkdir -p ${work_dir}/${layer}/iso/loader/entries
    cp ${script_path}/iso/${layer}/efiboot/loader/loader.conf ${work_dir}/${layer}/iso/loader/
    cp ${script_path}/iso/${layer}/efiboot/loader/entries/uefi-shell-v2-x86_64.conf ${work_dir}/${layer}/iso/loader/entries/
    cp ${script_path}/iso/${layer}/efiboot/loader/entries/uefi-shell-v1-x86_64.conf ${work_dir}/${layer}/iso/loader/entries/

    sed "s|%ARCHISO_LABEL%|${iso_label}|g;
         s|%INSTALL_DIR%|${install_dir}|g
         s|%ARCHISO_LAYER%|${layer}|g;
         s|%ARCHISO_LAYERS%|${archiso_layers}|g;
         s|%OPEROS_VERSION%|${iso_version}|g;" \
        ${script_path}/iso/${layer}/efiboot/loader/entries/archiso-x86_64-usb.conf > ${work_dir}/${layer}/iso/loader/entries/archiso-x86_64.conf

    # EFI Shell 2.0 for UEFI 2.3+
    curl -o ${work_dir}/${layer}/iso/EFI/shellx64_v2.efi https://raw.githubusercontent.com/tianocore/edk2/master/ShellBinPkg/UefiShell/X64/Shell.efi
    # EFI Shell 1.0 for non UEFI 2.3+
    curl -o ${work_dir}/${layer}/iso/EFI/shellx64_v1.efi https://raw.githubusercontent.com/tianocore/edk2/master/EdkShellBinPkg/FullShell/X64/Shell_Full.efi
}

# Prepare efiboot.img::/EFI for "El Torito" EFI boot mode
make_efiboot() {
    local archiso_layers=$@
    mkdir -p ${work_dir}/${layer}/iso/EFI/archiso
    truncate -s 64M ${work_dir}/${layer}/iso/EFI/archiso/efiboot.img
    mkfs.fat -n ARCHISO_EFI ${work_dir}/${layer}/iso/EFI/archiso/efiboot.img

    mkdir -p ${work_dir}/${layer}/efiboot
    mount ${work_dir}/${layer}/iso/EFI/archiso/efiboot.img ${work_dir}/${layer}/efiboot

    mkdir -p ${work_dir}/${layer}/efiboot/EFI/archiso
    cp ${work_dir}/${layer}/iso/${install_dir}/boot/x86_64/vmlinuz ${work_dir}/${layer}/efiboot/EFI/archiso/vmlinuz.efi
    cp ${work_dir}/${layer}/iso/${install_dir}/boot/x86_64/archiso.img ${work_dir}/${layer}/efiboot/EFI/archiso/archiso.img

    cp ${work_dir}/${layer}/iso/${install_dir}/boot/intel_ucode.img ${work_dir}/${layer}/efiboot/EFI/archiso/intel_ucode.img

    mkdir -p ${work_dir}/${layer}/efiboot/EFI/boot
    cp ${work_dir}/airootfs/usr/share/efitools/efi/PreLoader.efi ${work_dir}/${layer}/efiboot/EFI/boot/bootx64.efi
    cp ${work_dir}/airootfs/usr/share/efitools/efi/HashTool.efi ${work_dir}/${layer}/efiboot/EFI/boot/

    cp ${work_dir}/airootfs/usr/lib/systemd/boot/efi/systemd-bootx64.efi ${work_dir}/${layer}/efiboot/EFI/boot/loader.efi

    mkdir -p ${work_dir}/${layer}/efiboot/loader/entries
    cp ${script_path}/iso/${layer}/efiboot/loader/loader.conf ${work_dir}/${layer}/efiboot/loader/
    cp ${script_path}/iso/${layer}/efiboot/loader/entries/uefi-shell-v2-x86_64.conf ${work_dir}/${layer}/efiboot/loader/entries/
    cp ${script_path}/iso/${layer}/efiboot/loader/entries/uefi-shell-v1-x86_64.conf ${work_dir}/${layer}/efiboot/loader/entries/

    sed "s|%ARCHISO_LABEL%|${iso_label}|g;
         s|%INSTALL_DIR%|${install_dir}|g
         s|%ARCHISO_LAYER%|${layer}|g;
         s|%ARCHISO_LAYERS%|${archiso_layers}|g;" \
        ${script_path}/iso/${layer}/efiboot/loader/entries/archiso-x86_64-cd.conf > ${work_dir}/${layer}/efiboot/loader/entries/archiso-x86_64.conf

    cp ${work_dir}/${layer}/iso/EFI/shellx64_v2.efi ${work_dir}/${layer}/efiboot/EFI/
    cp ${work_dir}/${layer}/iso/EFI/shellx64_v1.efi ${work_dir}/${layer}/efiboot/EFI/

    umount -d ${work_dir}/${layer}/efiboot
}

_cleanup () {
    echo "Cleaning up what we can on airootfs..."

    # Delete pacman database sync cache files (*.tar.gz)
    if [[ -d "${work_dir}/airootfs/var/lib/pacman" ]]; then
        find "${work_dir}/airootfs/var/lib/pacman" -maxdepth 1 -type f -delete
    fi
    # Delete pacman database sync cache
    if [[ -d "${work_dir}/airootfs/var/lib/pacman/sync" ]]; then
        find "${work_dir}/airootfs/var/lib/pacman/sync" -delete
    fi
    # Delete pacman package cache
    if [[ -d "${work_dir}/airootfs/var/cache/pacman/pkg" ]]; then
        find "${work_dir}/airootfs/var/cache/pacman/pkg" -type f -delete
    fi
    # Delete all log files, keeps empty dirs.
    if [[ -d "${work_dir}/airootfs/var/log" ]]; then
        find "${work_dir}/airootfs/var/log" -type f -delete
    fi
    # Delete all temporary files and dirs
    if [[ -d "${work_dir}/airootfs/var/tmp" ]]; then
        find "${work_dir}/airootfs/var/tmp" -mindepth 1 -delete
    fi
    # Delete package pacman related files.
    find "${work_dir}" \( -name "*.pacnew" -o -name "*.pacsave" -o -name "*.pacorig" \) -delete
    echo "Done!"
}

make_final_cleanup() {
    # Delete initcpio image(s)
    if [[ -d "${work_dir}/airootfs/boot" ]]; then
        find "${work_dir}/airootfs/boot" -type f -name '*.img' -delete
    fi
    # Delete kernel(s)
    if [[ -d "${work_dir}/airootfs/boot" ]]; then
        find "${work_dir}/airootfs/boot" -type f -name 'vmlinuz*' -delete
    fi
}

# Build airootfs filesystem image
make_squash_layer() {
    _cleanup

    umount ${work_dir}/airootfs
    rm -rf ${work_dir}/overlaywork

    cp -a -l -f ${work_dir}/airootfs-${layer} ${work_dir}

    mkdir -p "${work_dir}/sfs/${install_dir}/${arch}"
    echo "Creating SquashFS image, this may take some time..."
    if [[ -z "${verbose}" ]]; then
        mksquashfs "${work_dir}/airootfs-${layer}" "${work_dir}/sfs/${install_dir}/${arch}/airootfs-${layer}.sfs" -noappend -comp "${sfs_comp}" -no-progress &> /dev/null
    else
        mksquashfs "${work_dir}/airootfs-${layer}" "${work_dir}/sfs/${install_dir}/${arch}/airootfs-${layer}.sfs" -noappend -comp "${sfs_comp}"
    fi
    echo "Done!"

    echo "Creating checksum file for self-test..."
    cd "${work_dir}/sfs/${install_dir}/${arch}"
    sha512sum "airootfs-${layer}.sfs" > "airootfs-${layer}.sha512"
    cd ${OLDPWD}
    echo "Done!"

}

make_prepare_layer() {
    mkdir -p ${work_dir}/airootfs-${layer} ${work_dir}/airootfs ${work_dir}/overlaywork
    mount -t overlay overlay \
        -olowerdir=$(_join : $(printf "${work_dir}/airootfs-%s\n" "$@")) \
        -oupperdir=${work_dir}/airootfs-${layer} \
        -oworkdir=${work_dir}/overlaywork \
        ${work_dir}/airootfs
}

# Build ISO
make_iso() {
    mkarchiso ${verbose} -w "${work_dir}/${layer}" -D "${install_dir}" -L "${iso_label}" -o "${out_dir}" iso "${iso_name}-${layer}${devbuild}-${iso_version}.iso"
}

make_unmount() {
    if mount | grep -q "${work_dir}/airootfs" > /dev/null; then
        umount "${work_dir}/airootfs"
    fi
}

while getopts 'N:V:L:D:w:o:g:vdh' arg; do
    case "${arg}" in
        N) iso_name="${OPTARG}" ;;
        V) iso_version="${OPTARG}" ;;
        L) iso_label="${OPTARG}" ;;
        D) install_dir="${OPTARG}" ;;
        w) work_dir="${OPTARG}" ;;
        o) out_dir="${OPTARG}" ;;
        g) gpg_key="${OPTARG}" ;;
        v) verbose="-v" ;;
        d) devbuild="-dev" && sfs_comp="gzip" ;;
        h) _usage 0 ;;
        *)
           echo "Invalid argument '${arg}'"
           _usage 1
           ;;
    esac
done
shift $((OPTIND -1))

if [[ $# -lt 1 ]]; then
    _usage 1
fi

trap make_unmount EXIT

mkdir -p ${work_dir}

for layer in $@; do
    . iso/${layer}/build-layer.sh
done
