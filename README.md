# Operos

Operos is a Linux-based operating system that brings hyperscaler-grade
infrastructure automation to organizations of all sizes: scheduled containers,
software defined networking, and converged storage automatically provisioned on
commodity x86 servers.

Operos combines a number of open source technologies into a single cohesive
cloud-native platform:

- [Kubernetes](https://kubernetes.io/) for container orchestration
- [Ceph](http://ceph.com/) for distributed storage
- [Calico](https://www.projectcalico.org/) for software-defined container networking
- [Prometheus](https://prometheus.io/) for metrics collection
- [isc-dhcpd](https://www.isc.org/downloads/dhcp)/[NginX](https://www.nginx.com/)/[SYSLINUX](http://www.syslinux.org) for hardware provisioning
- [Arch Linux](https://www.archlinux.org/) as the platform

In addition to the above, Operos includes several original components:

- [Teamster](components/teamster) and [Prospector](components/prospector) for node management
- [Waterfront](components/waterfront) as an additional GUI
- A fast [installer](components/installer) for provisioning controllers

For more information about Operos, see [its home page](https://www.paxautoma.com/operos/).

## Get Operos

The easiest way to get started with Operos is to download a binary ISO image:

[Download the latest ISO binary here](https://www.paxautoma.com/download/operos-iso).

Read [the installation instructions](https://www.paxautoma.com/operos/docs/0.2.0.html).

## Building from source

1. Run `make` to build everything from scratch. See below for how to rebuild
   various parts of the system.

2. You should now see an installer ISO in the `out` directory.

## Pre-requisites

- You will need the archlinux64 box for Vagrant. This can be created via:
  [packer-arch](https://github.com/elasticdog/packer-arch).

        git clone git@github.com:elasticdog/packer-arch.git
        cd packer-arch
        ./wrapacker
        vagrant box add -f --name archlinux64 output/packer_arch_virtualbox.box 

## Running the generated ISO

To run the ISO, create virtual machines in VirtualBox. You'll need one machine
for the controller and one or more workers. The controller node needs at least
2GB of RAM and 2 CPUs. The worker nodes need 2GB of RAM and one CPU.

The controller should have at least two network interfaces:

- The first (external) interface should be connected externally. This can be
  done via bridged or NAT network types.
- The second one should be connected to a VirtualBox host-only network (e.g.
  vboxnet1). This will be used for cross-node communication. Disable any DHCP
  servers on this network (in VirtualBox settings) as the controller will run
  its own DHCP server.

The worker should have at least one network interface, connected to the same
host-only network.

After the controller installed, the Kubernetes API can be accessed via the
provided [kubectl](kubectl) script (note that the kubectl binary must be
installed on the machine). This script will automatically fetch the user
credentials from the controller if this has not already been done.

## Version number

- The version number is formatted as: `x.y.z`. The `x.y` portion is defined in
  the file [operos-version](operos-version). `z` is intended to be the build
  number in the CI system. This can be set via the make variable `BUILD_NUM`:

        make isobuild BUILD_NUM=123

  This value defaults to `x`, to indicate an unofficial build.

## Docker image and Arch package cache

The Docker images and Arch packages used during builds are cached in the build
tree. To refresh, use:

    # Refresh Arch package cache
    make packages
    # Refresh Docker image cache
    make images

The versions of Docker images to be used are specified in [versions](versions).
The cache must be built at least once before running the build. It can also be
rebuild any time to obtain the latest packages and images.

## Rebuilding the ISO only

To rebuild only the ISO, skipping the cache updates, use:

    make isobuild

## Development build

There is a special, development build of the Operos ISO that can be built
using:

    make isobuild-dev

Differences between the development and production builds:

- An SSH key is automatically generated (`keys/testkey[.pub]`) and set as an
  authorized key on all nodes, controller and worker. This makes it easy to log
  into the nodes without having to enter a password, for example:

        ssh -i keys/testkey root@192.168.33.10

- When creating the images, gzip compression is used (instead of xz for
  production). This takes less time, but produces larger images.
