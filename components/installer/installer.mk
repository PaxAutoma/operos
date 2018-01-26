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

INSTALLER_FILES=$(shell find components/installer/ -name "*.go")

.PHONY: installer-novm
installer-novm: iso/installer/airootfs/root/installer
iso/installer/airootfs/root/installer: .buildnum vendor $(INSTALLER_FILES) $(COMMON_FILES)
	mkdir -p $(dir $@)
	go build \
		-v -i \
		-ldflags "-X main.operosVersion=$(ISO_VERSION)" \
		-o $@ ./components/installer/cmd/installer.go

clean: clean-installer

clean-installer:
	rm -f iso/installer/airootfs/root/installer
