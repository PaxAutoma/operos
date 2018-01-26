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

WATERFRONT_VERSION=latest
WATERFRONT_BUILDIMG=waterfront:build
WATERFRONT_SHIPIMG=waterfront:$(WATERFRONT_VERSION)

WATERFRONT_PATH=components/waterfront

WATERFRONT_SERVER_FILES=$(shell find $(WATERFRONT_PATH)/server -name "*.go")
WATERFRONT_CLIENT_FILES=$(shell \
	find $(WATERFRONT_PATH)/client -path \
		 $(WATERFRONT_PATH)/client/node_modules -prune -o \
		 -type f -print)

WATERFRONT_PROTOS=\
	$(WATERFRONT_PATH)/server/pkg/waterfront/waterfront.pb.go \
	$(WATERFRONT_PATH)/server/pkg/waterfront/waterfront.pb.gw.go \
	components/teamster/pkg/teamster/teamster.pb.go


waterfront: waterfront-ship-container

.PHONY: waterfront-ship-container
waterfront-ship-container: iso/controller/airootfs/opt/docker-images/waterfront.tar.gz

iso/controller/airootfs/opt/docker-images/waterfront.tar.gz: waterfront-docker
	docker build -t $(WATERFRONT_SHIPIMG) -f $(WATERFRONT_PATH)/containers/ship/Dockerfile $(WATERFRONT_PATH)
	mkdir -p $(dir $@)
	docker save $(WATERFRONT_SHIPIMG) | gzip > $@

.PHONY: waterfront-docker
waterfront-docker: waterfront-build-container
	docker run \
		--rm \
		-v `pwd`:/go/src/github.com/paxautoma/operos \
		$(WATERFRONT_BUILDIMG) \
		make waterfront-nodocker

.PHONY: waterfront-build-container
waterfront-build-container:
	docker build --pull -t $(WATERFRONT_BUILDIMG) -f $(WATERFRONT_PATH)/containers/build/Dockerfile $(WATERFRONT_PATH)

.PHONY: waterfront
waterfront-nodocker: waterfront-client waterfront-server

waterfront-client: $(WATERFRONT_PATH)/dist/client/app.js

$(WATERFRONT_PATH)/dist/client/app.js: $(WATERFRONT_PATH)/client/node_modules $(WATERFRONT_CLIENT_FILES)
	mkdir -p $(WATERFRONT_PATH)/dist/client
	cp -r $(WATERFRONT_PATH)/client/src/static/ $(WATERFRONT_PATH)/dist/client
	cp $(WATERFRONT_PATH)/client/src/index.html $(WATERFRONT_PATH)/dist/client
	cd $(WATERFRONT_PATH)/client && yarn run build

$(WATERFRONT_PATH)/client/node_modules: $(WATERFRONT_PATH)/client/yarn.lock
	cd $(WATERFRONT_PATH)/client && yarn install

.PHONY: waterfront-server
waterfront-server: $(WATERFRONT_PATH)/dist/waterfront
$(WATERFRONT_PATH)/dist/waterfront: vendor $(WATERFRONT_PROTOS) $(WATERFRONT_SERVER_FILES)
	go build -o $@ -v $(WATERFRONT_PATH)/server/cmd/main.go


.PHONY: waterfront-clean
waterfront-clean:
	rm -rf $(WATERFRONT_PATH)/dist/
	rm -f $(WATERFRONT_PROTOS)
	rm -f iso/controller/airootfs/opt/docker-images/waterfront.tar.gz

.PHONY: waterfront-distclean
waterfront-distclean: waterfront-clean
	rm -rf $(WATERFRONT_PATH)/client/node_modules

clean: waterfront-clean
distclean: waterfront-distclean


.PHONY: waterfront-dev
waterfront-dev:
	docker run \
		-it \
		--rm \
		-v `pwd`:/go/src/github.com/paxautoma/operos \
		$(WATERFRONT_BUILDIMG)
