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

FROM paxautoma/golang-dev:latest

RUN apk add --no-cache \
    nodejs \
    yarn \
    protobuf \
    protobuf-dev \
    g++ \
    linux-pam-dev

# these have to be installed in the $GOPATH instead of vendor because they come
# with command-line tools, which cannot be vendored
RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway && \
    go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger && \
    go get -u github.com/golang/protobuf/protoc-gen-go

WORKDIR /go/src/github.com/paxautoma/operos

