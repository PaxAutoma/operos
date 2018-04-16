/*
Copyright 2018 Pax Automa Systems, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package installer

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jroimartin/gocui"
	"github.com/paxautoma/operos/components/installer/pkg/network"
)

type InstallerResponses struct {
	OrgInfo struct {
		Cluster      string
		Organization string
		Department   string
		City         string
		Province     string
		Country      string
	}
	PrivateInterface        string
	PublicNetwork           network.InterfaceSettings
	PrivateSubnet           string
	PodSubnet               string
	ServiceSubnet           string
	PrivateGateway          string
	PublicHostname          string
	DNSDomain               string
	StorageSystemPercentage int
	ControllerDisk          string
	RootPassword            string

	ControllerIP     string
	DNSIP            string
	KubeAPIServiceIP string
}

func (ir InstallerResponses) StorageDataPercentage() int {
	return 100 - ir.StorageSystemPercentage
}

func (it InstallerResponses) PublicIPInfo() string {
	if it.PublicNetwork.Mode == "dhcp" {
		return "DHCP"
	}
	return fmt.Sprintf("Static IP: %s, Gateway: %s", it.PublicNetwork.Subnet, it.PublicNetwork.Gateway)
}

type InstallerContext struct {
	Interfaces struct {
		ByName  map[string]*network.InterfaceInfo
		Ordered []*network.InterfaceInfo
	}
	Disks             []DiskInfo
	Responses         InstallerResponses
	Versions          []string
	ControllerCert    string
	ControllerKey     string
	ServerCert        string
	ServerKey         string
	G                 *gocui.Gui
	Net               network.NetworkConfigurator
	GatekeeperAddress string
	GatekeeperTLS     bool
	InstallID         string
	OperosVersion     string
}

var DefaultContext InstallerContext

func init() {
	// Initialize default values in the context
	DefaultContext.Responses.PublicNetwork.Mode = "dhcp"
	DefaultContext.Responses.PrivateSubnet = "192.168.33.10/24"
	DefaultContext.Responses.PrivateGateway = "192.168.33.10"
	DefaultContext.Responses.PodSubnet = "10.10.0.0/16"
	DefaultContext.Responses.ServiceSubnet = "10.11.0.0/16"
	DefaultContext.Responses.DNSDomain = "cluster.local"
	DefaultContext.Responses.StorageSystemPercentage = 50

	DefaultContext.InstallID = uuid.New().String()
}
