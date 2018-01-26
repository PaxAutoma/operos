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

mkdir -p /mnt/etc/kubernetes/ssl

echo "${INSTALLER_CONTROLLER_CERT}" > /mnt/etc/kubernetes/ssl/ca.pem
echo "${INSTALLER_CONTROLLER_KEY}" > /mnt/etc/kubernetes/ssl/controller-ca-key.pem
echo "${INSTALLER_CONTROLLER_CERT}" > /mnt/etc/kubernetes/ssl/controller-ca.pem
echo "${INSTALLER_SERVER_KEY}" > /mnt/etc/kubernetes/ssl/apiserver-key.pem
echo "${INSTALLER_SERVER_CERT}" > /mnt/etc/kubernetes/ssl/apiserver.pem

mkdir -p /mnt/etc/paxautoma/tls
