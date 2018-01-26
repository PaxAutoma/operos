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

DIAG_FILES=$(shell find components/diag/ -type f)

.PHONY: diag-novm
diag-novm: iso/base/airootfs/usr/bin/operosdiagauth

iso/base/airootfs/usr/bin/operosdiagauth: $(DIAG_FILES) $(COMMON_FILES) components/common/gatekeeper/gatekeeper.pb.go vendor
	mkdir -p $(dir $@)
	go build -v -o $@ ./components/diag/operosdiagauth.go
	cp components/diag/operosdiag.sh $(dir $@)/operosdiag
	cp components/diag/installerdiag.sh iso/installer/airootfs/root/installerdiag.sh

clean: clean-diag

clean-diag:
	rm -f iso/base/airootfs/usr/bin/operosdiagauth
	rm -f iso/base/airootfs/usr/bin/operosdiag
	rm -f iso/installer/airootfs/root/installerdiag.sh
