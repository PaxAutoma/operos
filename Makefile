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

BUILD_NUM=x
ISO_VERSION=$(shell cat operos-version).$(BUILD_NUM)

IMG_PREFIX_WORKER=operos-worker
IMG_PREFIX_CONTROLLER=operos-controller

.PHONY: all setup
all: setup isobuild
setup: vagrant packages images

include components/common.mk
include components/installer/installer.mk
include components/statustty/statustty.mk
include components/diag/diag.mk
include components/prospector/prospector.mk
include components/teamster/teamster.mk
include components/ipxe/ipxe.mk
include components/waterfront/waterfront.mk

# Write the .buildnum file if the build number changed since the last build.
# This file can serve as a prerequisite for builds that depend on the build
# number.
$(shell \
	if [ ! -e .buildnum ] || [ "`cat .buildnum`" != "$(BUILD_NUM)" ]; then \
		echo $(BUILD_NUM) > .buildnum; \
	fi > /dev/null \
)

WORKER_FILES=$(shell find iso/worker/ -type f)
CONTROLLER_FILES=$(shell find iso/controller/ -type f)

.PHONY: vagrant
vagrant:
	if ! vagrant plugin list | grep -q reload; then \
		vagrant plugin install vagrant-reload; \
	fi
	vagrant up

.PHONY: packages images test isobuild

packages images test isobuild isobuild-dev installer statustty diag prospector teamster ipxe:
	time vagrant ssh -c "cd ~/operos && make $@-novm BUILD_NUM=$(BUILD_NUM)"

.PHONY: packages-novm
packages-novm:
	sudo ./update-pkgs.sh

.PHONY: images-novm
images-novm:
	sudo ./update-imgs.sh

isobuild: waterfront
isobuild-dev: waterfront

isobuild-novm: installer-novm statustty-novm diag-novm teamster-novm prospector-novm ipxe-novm
	cp versions iso/installer/airootfs/root/versions
	sudo ./vagrant-build.sh -V $(ISO_VERSION) -L OPEROS_$(subst .,_,$(ISO_VERSION)) base node worker controller installer

isobuild-dev-novm: installer-novm statustty-novm diag-novm teamster-novm prospector-novm ipxe-novm keys/testkey
	cp versions iso/installer/airootfs/root/versions
	sudo ./vagrant-build.sh -V $(ISO_VERSION) -L OPEROS_$(subst .,_,$(ISO_VERSION)) -d base node worker controller installer

keys/testkey:
	mkdir -p keys
	ssh-keygen -f $@ -N ""

.PHONY: test-novm
test-novm:
	go test -v ./src/...

vendor: glide.lock
	rm -rf vendor
	glide install
	# this file does not compile under alpine (depends on glibc) and is only
	# used during testing of the go-crypt package, so can be safely deleted.
	rm vendor/github.com/amoghe/go-crypt/version.go

clean:
	rm -f keys/testkey
	rm -f keys/testkey.pub
	rm -rf vendor
	rm -f out/*
	rm -rf operos-credentials/

distclean: clean
	rm -rf cache/
