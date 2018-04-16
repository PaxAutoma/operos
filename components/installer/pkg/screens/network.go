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

package screens

import (
	"fmt"
	"image"
	"net"
	"strings"

	"github.com/jroimartin/gocui"

	"github.com/paxautoma/operos/components/common/widgets"
	installer "github.com/paxautoma/operos/components/installer/pkg"
	"github.com/paxautoma/operos/components/installer/pkg/network"
)

func NetworkSettingsPrivateScreen(screenSet *widgets.ScreenSet, context interface{}) *widgets.Screen {
	ctx := context.(*installer.InstallerContext)

	screen := widgets.NewScreen()
	screen.Title = "Network settings: private interface"
	screen.Message = `
An Operos Controller must be attached to two networks: private and the public.
The private network connects the Controller to the worker nodes, and also
connects the worker nodes to each other. When worker nodes boot, they obtain
their IP address (via DHCP) and software image via this network. The controller
will run a DHCP server on this interface, so no other DHCP servers should be
listening.

Which interface should be used for the private network?`

	items := make([]widgets.MenuItem, len(ctx.Interfaces.Ordered))
	selectedIdx := -1
	for idx, iface := range ctx.Interfaces.Ordered {
		items[idx] = formatInterfaceMenuItem(*iface)
		if iface.Name == ctx.Responses.PrivateInterface {
			selectedIdx = idx
		}
	}

	menu := widgets.NewMenu("menu-private-ifs", items, 80, 8)
	menu.OnSelectItem = func(item widgets.MenuItem) error {
		smi := item.(*widgets.SimpleMenuItem)
		ctx.Responses.PrivateInterface = smi.Value
		return nil
	}
	if selectedIdx >= 0 {
		menu.SelectItem(selectedIdx)
	}

	screen.Content = menu

	screen.OnPrev = func() error {
		screenSet.Back(1)
		return nil
	}

	screen.OnNext = func() error {
		if ctx.Interfaces.ByName[ctx.Responses.PrivateInterface].DhcpOffer != "" {
			message := fmt.Sprintf(`
 A DHCP server was detected on the network connected to %s.
 
 This will likely interfere with the operation of the Operos
 controller.

 Are you sure you wish to continue?`, ctx.Responses.PrivateInterface)

			dialog := widgets.NewDialog("confirm-dhcp", message, 70)
			dialog.OnAccept = func() error {
				dialog.Close()
				screenSet.Forward(1)
				return nil
			}

			dialog.OnCancel = func() error {
				dialog.Close()
				screen.Focus()
				return nil
			}

			screen.Blur()
			return dialog.Layout(ctx.G)
		}

		screenSet.Forward(1)
		return nil
	}

	return screen
}

func NetworkSettingsPublicIfaceScreen(screenSet *widgets.ScreenSet, context interface{}) *widgets.Screen {
	ctx := context.(*installer.InstallerContext)

	screen := widgets.NewScreen()
	screen.Title = "Network settings: public interface"
	screen.Message = `The public network connects the Controller to the outside world. This interface
will be used to access the Controller UI and API.

Which interface should be used for the public network?`

	items := make([]widgets.MenuItem, len(ctx.Interfaces.Ordered))
	idx := 0
	selectedIdx := -1
	for _, iface := range ctx.Interfaces.Ordered {
		if iface.Name != ctx.Responses.PrivateInterface {
			items[idx] = formatInterfaceMenuItem(*iface)
			if iface.Name == ctx.Responses.PublicNetwork.Interface {
				selectedIdx = idx
			}
			idx++
		}
	}
	items[idx] = &widgets.SimpleMenuItem{
		Text:  "<disabled>",
		Value: "",
	}

	menu := widgets.NewMenu("menu-public-ifs", items, 80, 8)
	menu.OnSelectItem = func(item widgets.MenuItem) error {
		smi := item.(*widgets.SimpleMenuItem)

		if smi.Value == "" {
			ctx.Responses.PublicNetwork.Mode = "disabled"
			// force user to set the gateway for the private network
			ctx.Responses.PrivateGateway = ""
		} else if ctx.Responses.PublicNetwork.Mode == "disabled" {
			ctx.Responses.PublicNetwork = installer.DefaultContext.Responses.PublicNetwork
		}
		ctx.Responses.PublicNetwork.Interface = smi.Value
		return nil
	}

	if selectedIdx >= 0 {
		menu.SelectItem(selectedIdx)
	}

	screen.Content = menu
	screen.OnPrev = func() error {
		screenSet.Back(1)
		return nil
	}
	screen.OnNext = func() error {
		if ctx.Responses.PublicNetwork.Mode == "disabled" {
			screenSet.Forward(2)
			return nil
		}

		if ctx.Responses.PublicNetwork.Mode == "dhcp" && ctx.Interfaces.ByName[ctx.Responses.PublicNetwork.Interface].DhcpOffer == "" {
			message := fmt.Sprintf(`
 A DHCP server was not detected on the network connected to %s;
 however, the interface is set to use DHCP.

 Are you sure you wish to continue?`, ctx.Responses.PublicNetwork.Interface)

			dialog := widgets.NewDialog("confirm-dhcp", message, 70)
			dialog.OnAccept = func() error {
				dialog.Close()
				screenSet.Forward(1)
				return nil
			}

			dialog.OnCancel = func() error {
				dialog.Close()
				screen.Focus()
				return nil
			}

			screen.Blur()
			return dialog.Layout(ctx.G)
		}

		screenSet.Forward(1)
		return nil
	}

	return screen
}

func formatInterfaceMenuItem(iface network.InterfaceInfo) *widgets.SimpleMenuItem {
	label := fmt.Sprintf("%s [%s]", iface.Name, iface.Mac)
	if iface.DhcpOffer != "" {
		label += fmt.Sprintf(" - DHCP detected (%s)", iface.DhcpOffer)
	}
	return &widgets.SimpleMenuItem{
		Text:  label,
		Value: iface.Name,
	}
}

func validatePublicNetwork(pn network.InterfaceSettings) []error {
	result := make([]error, 0)

	subnetIP, _, err := net.ParseCIDR(pn.Subnet)
	if err != nil {
		return result
	}

	gateway := net.ParseIP(pn.Gateway)
	if gateway == nil {
		return result
	}

	if subnetIP.Equal(gateway) {
		result = append(result, widgets.NewValidationError("Gateway", "should not have the same IP as the controller host"))
	}

	return result
}

func NetworkSettingsPublicIpsScreen(screenSet *widgets.ScreenSet, context interface{}) *widgets.Screen {
	ctx := context.(*installer.InstallerContext)

	screen := widgets.NewScreen()
	screen.Title = "Network settings: public interface"
	screen.Message = `Please configure the public interface. This is required during installation in
order to obtain a TLS certificate and license from Pax Automa.`

	menuDhcp := widgets.NewMenu("menu-dhcp", []widgets.MenuItem{
		&widgets.SimpleMenuItem{"Dynamic configuration via DHCP", "dhcp"},
		&widgets.SimpleMenuItem{"Static configuration", "static"},
	}, 80, 4)

	menuStatic := widgets.NewEditableList("public-static-settings", []*widgets.EditableListItem{
		widgets.NewEditableListItem("IP address/prefix", "subnet", ctx.Responses.PublicNetwork.Subnet, widgets.ValidateIPNet),
		widgets.NewEditableListItem("Gateway", "gateway", ctx.Responses.PublicNetwork.Gateway, widgets.ValidateIP),
	}, 80, 4)

	errorList := widgets.NewPar("par-errors", "")
	errorList.Bounds = image.Rect(1, 0, 79, 3)
	errorList.FgColor = gocui.ColorRed

	var valid bool
	validate := func() {
		if ctx.Responses.PublicNetwork.Mode == "dhcp" {
			valid = true
			errorList.SetText("")
		} else {
			errors := menuStatic.Validate()
			errors = append(errors, validatePublicNetwork(ctx.Responses.PublicNetwork)...)
			errorList.SetText(widgets.JoinValidationErrors(errors))
			valid = len(errors) == 0
		}

		screen.ShowNext(valid)
	}

	validate()

	menuDhcp.OnSelectItem = func(item widgets.MenuItem) error {
		smi := item.(*widgets.SimpleMenuItem)
		if smi.Value == "static" {
			menuStatic.SetVisibility(true)
			ctx.Responses.PublicNetwork.Mode = "static"
		} else {
			menuStatic.SetVisibility(false)
			ctx.Responses.PublicNetwork.Mode = "dhcp"
		}

		validate()

		return nil
	}

	menuStatic.OnItemChange = func(item *widgets.EditableListItem) {
		switch item.Key {
		case "subnet":
			ctx.Responses.PublicNetwork.Subnet = item.Value
		case "gateway":
			ctx.Responses.PublicNetwork.Gateway = item.Value
		}

		validate()
	}

	vl := widgets.NewVerticalLayout()
	vl.Items = []widgets.Renderable{menuDhcp, menuStatic, errorList}
	screen.Content = vl

	screen.OnPrev = func() error {
		screenSet.Back(1)
		return nil
	}

	screen.OnNext = func() error {
		if !valid {
			return nil
		}
		screenSet.Forward(1)
		return nil
	}

	return screen
}

func validateNetworkIps(ctx *installer.InstallerContext) []error {
	result := []error{}

	_, privateSubnet, err := net.ParseCIDR(ctx.Responses.PrivateSubnet)
	if err != nil {
		return result
	}

	_, podSubnet, err := net.ParseCIDR(ctx.Responses.PodSubnet)
	if err != nil {
		return result
	}

	_, serviceSubnet, err := net.ParseCIDR(ctx.Responses.ServiceSubnet)
	if err != nil {
		return result
	}

	var publicSetting string
	if ctx.Responses.PublicNetwork.Mode == "dhcp" {
		publicSetting = ctx.Interfaces.ByName[ctx.Responses.PublicNetwork.Interface].DhcpOffer
	} else {
		publicSetting = ctx.Responses.PublicNetwork.Subnet
	}

	_, publicSubnet, err := net.ParseCIDR(publicSetting)
	if err != nil {
		return result
	}

	subnets := map[string]*net.IPNet{
		"private": privateSubnet,
		"public":  publicSubnet,
		"pod":     podSubnet,
		"service": serviceSubnet,
	}

	check := func(name1, name2 string) {
		if subnets[name1].Contains(subnets[name2].IP) || subnets[name2].Contains(subnets[name1].IP) {
			result = append(result, widgets.NewValidationError(
				fmt.Sprintf("%s subnet", strings.ToTitle(name1[0:1])+name1[1:]),
				fmt.Sprintf("should not overlap with %s subnet", name2)))
		}
	}

	check("private", "public")
	check("private", "pod")
	check("private", "service")
	check("public", "pod")
	check("public", "service")
	check("pod", "service")

	return result
}

func NetworkSettingsIpsScreen(screenSet *widgets.ScreenSet, context interface{}) *widgets.Screen {
	ctx := context.(*installer.InstallerContext)

	screen := widgets.NewScreen()
	screen.Title = "Network settings: IPs and domains"
	screen.Message = `The following IPs and domain settings will be used for the cluster. You can
edit the values below.`

	gatewayLabel := "Gateway for workers"
	if ctx.Responses.PublicNetwork.Mode == "disabled" {
		gatewayLabel = "Gateway"
	}

	menu := widgets.NewEditableList("eli-ips", []*widgets.EditableListItem{
		widgets.NewEditableListItem("Private subnet", "private-subnet", ctx.Responses.PrivateSubnet, widgets.ValidateIPNet),
		widgets.NewEditableListItem("Pod subnet", "pod-subnet", ctx.Responses.PodSubnet, widgets.ValidateIPNet),
		widgets.NewEditableListItem("Service subnet", "service-subnet", ctx.Responses.ServiceSubnet, widgets.ValidateIPNet),
		widgets.NewEditableListItem(gatewayLabel, "private-gateway", ctx.Responses.PrivateGateway, widgets.ValidateIP),
		widgets.NewEditableListItem("Public hostname", "public-hostname", "", nil),
		widgets.NewEditableListItem("DNS domain", "dns-domain", ctx.Responses.DNSDomain, widgets.ValidateNotEmpty),
	}, 80, 8)

	errorList := widgets.NewPar("par-errors", "")
	errorList.Bounds = image.Rect(1, 0, 79, 3)
	errorList.FgColor = gocui.ColorRed

	var valid bool
	validate := func() {
		errors := menu.Validate()
		errors = append(errors, validateNetworkIps(ctx)...)
		errorList.SetText(widgets.JoinValidationErrors(errors))
		valid = len(errors) == 0
		screen.ShowNext(valid)
	}

	validate()

	menu.OnItemChange = func(item *widgets.EditableListItem) {
		switch item.Key {
		case "private-subnet":
			ctx.Responses.PrivateSubnet = item.Value
		case "pod-subnet":
			ctx.Responses.PodSubnet = item.Value
		case "service-subnet":
			ctx.Responses.ServiceSubnet = item.Value
		case "dns-domain":
			ctx.Responses.DNSDomain = item.Value
		case "private-gateway":
			ctx.Responses.PrivateGateway = item.Value
		case "public-hostname":
			ctx.Responses.PublicHostname = item.Value
		}

		validate()
	}

	vl := widgets.NewVerticalLayout()
	vl.Items = []widgets.Renderable{menu, errorList}
	screen.Content = vl

	screen.OnPrev = func() error {
		if ctx.Responses.PublicNetwork.Mode == "disabled" {
			screenSet.Back(2)
		} else {
			screenSet.Back(1)
		}
		return nil
	}

	screen.OnNext = func() error {
		if valid {
			privateIP, _, err := net.ParseCIDR(ctx.Responses.PrivateSubnet)
			if err != nil {
				panic(err)
			}

			//ctx.Responses.ControllerIP = network.IncrementIP(privateSubnet.IP, 10).String()
			ctx.Responses.ControllerIP = privateIP.String()

			_, serviceSubnet, err := net.ParseCIDR(ctx.Responses.ServiceSubnet)
			if err != nil {
				panic(err)
			}
			ctx.Responses.KubeAPIServiceIP = network.IncrementIP(serviceSubnet.IP, 1).String()
			ctx.Responses.DNSIP = network.IncrementIP(serviceSubnet.IP, 2).String()

			screenSet.Forward(1)
		}
		return nil
	}

	return screen
}
