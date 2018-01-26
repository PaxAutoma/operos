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

make_docker() {
    mkdir -p ${work_dir}/airootfs/opt/docker-images
    cp -af --no-preserve=ownership ${script_path}/cache/docker/* ${work_dir}/airootfs/opt/docker-images/
}

run_once make_prepare_layer node base
run_once make_packages
run_once make_customize_airootfs
run_once make_syslinux base node controller
run_once make_docker
run_once make_final_cleanup
run_once make_squash_layer
