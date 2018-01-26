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

FROM alpine:3.6

RUN apk add --no-cache linux-pam
RUN mkdir -p /opt/waterfront
COPY dist/ /opt/waterfront/

RUN ln -sf /etc-host/shadow /etc/shadow && \
    ln -sf /etc-host/passwd /etc/passwd && \
    mkdir -p /etc/paxautoma && \
    ln -sf /etc-host/paxautoma/settings /etc/paxautoma/settings

EXPOSE 2780 2781
WORKDIR /opt/waterfront

# Waterfront needs to run as root and have the host's /etc mounted in at
# /etc-host so that it can perform PAM authentication.
#
# See http://jenkins-ci.361315.n4.nabble.com/Using-UNIX-PAM-authentication-from-a-non-root-user-tp378559p378563.html

CMD ["./waterfront"]
