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

make_prepare_iso_worker() {
    mkdir -p "${work_dir}/worker/iso/${install_dir}/${arch}"
    ln -f ${work_dir}/sfs/${install_dir}/${arch}/airootfs-base.* ${work_dir}/worker/iso/${install_dir}/${arch}/
    ln -f ${work_dir}/sfs/${install_dir}/${arch}/airootfs-node.* ${work_dir}/worker/iso/${install_dir}/${arch}/
    ln -f ${work_dir}/sfs/${install_dir}/${arch}/airootfs-worker.* ${work_dir}/worker/iso/${install_dir}/${arch}/
}

run_once make_prepare_layer node base
run_once make_packages
run_once make_customize_airootfs
# run_once make_setup_mkinitcpio
# run_once make_boot
# run_once make_isolinux
run_once make_syslinux base node worker
# run_once make_final_cleanup
run_once make_squash_layer
# run_once make_prepare_iso_worker
# run_once make_iso
