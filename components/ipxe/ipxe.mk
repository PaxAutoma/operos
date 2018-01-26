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

.PHONY: ipxe-novm
ipxe-novm: iso/controller/airootfs/tftpboot/mapfile iso/controller/airootfs/tftpboot/ipxe-x64.efi iso/controller/airootfs/tftpboot/ipxe.pxe

clean: ipxe-clean

components/ipxe/ipxe:
	cd ./components/ipxe ; \
	git clone git://git.ipxe.org/ipxe.git

iso/controller/airootfs/tftpboot/mapfile:
	mkdir -p $(dir $@)
	cp ./components/ipxe/mapfile iso/controller/airootfs/tftpboot/mapfile
	chown 14 iso/controller/airootfs/tftpboot/mapfile

iso/controller/airootfs/tftpboot/ipxe-x86.efi: components/ipxe/ipxe
	cd ./components/ipxe/ipxe/src ; \
	$(MAKE) bin-i386-efi/ipxe.efi ;
	cp ./components/ipxe/ipxe/src/bin-i386-efi/ipxe.efi iso/controller/airootfs/tftpboot/ipxe-x86.efi
	chown 14 iso/controller/airootfs/tftpboot/ipxe-x86.efi

iso/controller/airootfs/tftpboot/ipxe-x64.efi: components/ipxe/ipxe
	cd ./components/ipxe/ipxe/src ; \
	$(MAKE) bin-x86_64-efi/ipxe.efi
	cp ./components/ipxe/ipxe/src/bin-x86_64-efi/ipxe.efi iso/controller/airootfs/tftpboot/ipxe-x64.efi
	chown 14 iso/controller/airootfs/tftpboot/ipxe-x64.efi

iso/controller/airootfs/tftpboot/ipxe.pxe: components/ipxe/ipxe
	cd ./components/ipxe/ipxe/src ; \
	$(MAKE) bin/ipxe.pxe
	cp ./components/ipxe/ipxe/src/bin/ipxe.pxe iso/controller/airootfs/tftpboot/ipxe.pxe
	chown 14 iso/controller/airootfs/tftpboot/ipxe.pxe

ipxe-clean:
	rm -rf ./components/ipxe/ipxe
	rm -f iso/controller/airootfs/tftpboot/ipxe.pxe
	rm -f iso/controller/airootfs/tftpboot/ipxe-x64.efi
	rm -f iso/controller/airootfs/tftpboot/ipxe-x86.efi
	rm -f iso/controller/airootfs/tftpboot/mapfile
