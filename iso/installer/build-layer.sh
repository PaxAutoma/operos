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

make_devbuild() {
    cp ${script_path}/keys/testkey.pub ${work_dir}/airootfs/root/testkey.pub
    cp ${script_path}/iso/installer/300-devbuild.sh ${work_dir}/airootfs/root/install/300-devbuild.sh
}

make_prepare_iso_installer() {
    mkdir -p "${work_dir}/installer/iso/${install_dir}/${arch}"
    ln -f ${work_dir}/sfs/${install_dir}/${arch}/* ${work_dir}/installer/iso/${install_dir}/${arch}/
    cp -af ${work_dir}/worker/iso/${install_dir}/boot/syslinux-worker/ ${work_dir}/installer/iso/${install_dir}/boot/
    cp -af ${work_dir}/controller/iso/${install_dir}/boot/syslinux-controller/ ${work_dir}/installer/iso/${install_dir}/boot/
}

run_once make_prepare_layer base
run_once make_packages
run_once make_customize_airootfs
if [ -n "$devbuild" ]; then
    run_once make_devbuild
fi
run_once make_setup_mkinitcpio
run_once make_boot
run_once make_isolinux
run_once make_syslinux base installer
run_once make_efi base installer
run_once make_efiboot base installer
run_once make_final_cleanup
run_once make_squash_layer
run_once make_prepare_iso_installer
run_once make_iso
