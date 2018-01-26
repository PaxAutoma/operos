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

# Prepare an offline repo. This uses a dummy blank database in order to force
# pacman to download all of the dependencies as well. We set the pacman
# cachedir to the target dir to coerce it to place all the downloaded files
# there.

script_path=$(readlink -f ${0%/*})
repo_path=${script_path}/cache/pkg
tmp_path=$(mktemp -d)
trap "rm -rf $tmp_path" EXIT

pacman \
    --noconfirm \
    --dbpath ${tmp_path} \
    --cachedir ${repo_path} \
    -Syuw \
    base syslinux \
    $(grep -h -v ^# ${script_path}/iso/*/packages | sort -u)

# Create repo for installer. This cannot be called "local" because that seems
# to be a reserved name.
cd ${repo_path}
repo-add -n -R paxautoma-local.db.tar.gz *.pkg.tar.*
