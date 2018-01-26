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

cat <<EOF > /mnt/etc/issue
Operos Controller (\l)

EOF

cat <<EOF > /mnt/usr/lib/os-release
NAME="Operos Controller"
PRETTY_NAME="Operos Controller"
VERSION="${OPEROS_VERSION}"
ID=operos
ID_LIKE=archlinux
ANSI_COLOR="0;36"
HOME_URL="https://www.paxautoma.com/"
EOF
