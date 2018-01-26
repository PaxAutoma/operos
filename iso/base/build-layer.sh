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

# Setup custom pacman.conf with current cache directories.
make_pacman_conf() {
    local _cache_dirs
    _cache_dirs=($(pacman -v 2>&1 | grep '^Cache Dirs:' | sed 's/Cache Dirs:\s*//g'))
    sed -r "s|^#?\\s*CacheDir.+|CacheDir = $(echo -n ${_cache_dirs[@]})|g" ${script_path}/iso/pacman.conf > ${work_dir}/pacman.conf
}

# Base installation (airootfs)
make_basefs() {
    mkarchiso ${verbose} -w "${work_dir}" -C "${work_dir}/pacman.conf" -D "${install_dir}" init
}

make_prepare_layer_base() {
    mkdir -p ${work_dir}/airootfs-base ${work_dir}/airootfs
    mount -o bind ${work_dir}/airootfs-base ${work_dir}/airootfs
}

run_once make_prepare_layer_base
run_once make_pacman_conf
run_once make_basefs
run_once make_packages
run_once make_customize_airootfs
run_once make_squash_layer
