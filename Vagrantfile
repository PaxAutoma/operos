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

provision = <<EOF
set -x -e

cat > /etc/resolv.conf <<RESOLV
nameserver 10.128.128.1
RESOLV

cat >> /etc/pacman.conf <<PACMAN
[paxautoma]
SigLevel = Optional TrustAll
Server = https://arch-packages.paxautoma.com/
PACMAN

pacman --noconfirm --needed -Syu archiso docker rsync go git protobuf glide

# Go environment
mkdir -p /home/vagrant/go/{src,pkg,bin}
find /home/vagrant/go -xdev -type d -exec chown vagrant:vagrant {} \\;

cat > /home/vagrant/.bashrc <<BASHRC
export GOPATH=/home/vagrant/go
export PATH=/home/vagrant/go/bin:$PATH

# If not running interactively, skip the rest
[[ $- != *i* ]] && return

alias ls='ls --color=auto'
PS1='[\\u@\\h \\W]\\$ '
BASHRC

echo 'export GOPATH=/home/vagrant/go' >> /home/vagrant/.bashrc
echo 'export PATH=/home/vagrant/go/bin:\$PATH' >> /home/vagrant/.bashrc

if [ ! -e /home/vagrant/operos ]; then
  ln -s /home/vagrant/go/src/github.com/paxautoma/operos /home/vagrant/operos
fi

# Install Go generator for protobuf
sudo -u vagrant go get -v -u github.com/golang/protobuf/protoc-gen-go

# Enable docker
usermod -a -G docker vagrant
systemctl enable docker.service
EOF

Vagrant.configure("2") do |config|
  # https://github.com/elasticdog/packer-arch
  config.vm.box = "archlinux64"
  config.vm.box_check_update = true

  config.vm.network "private_network", ip: "192.168.33.5", auto_config: false

  config.vm.define "isobuild" do |isovm|
    isovm.vm.hostname = "isobuild"
    isovm.vm.provider "virtualbox" do |vb|
      vb.memory = "1024"
      vb.cpus = 4
    end

    isovm.vm.synced_folder ".", "/vagrant", disabled: true
    isovm.vm.synced_folder File.dirname(__FILE__), "/home/vagrant/go/src/github.com/paxautoma/operos"

    isovm.vm.provision "shell", inline: provision
    # If the Linux kernel is updated during provisioning, Docker will fail to 
    # start until the machine is reloaded. So we do that here for good measure.
    isovm.vm.provision :reload
  end

end
