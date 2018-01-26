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

package network

import (
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type InterfaceSettings struct {
	Interface string
	Mode      string
	Subnet    string
	Gateway   string
}

type NetworkConfigurator interface {
	ConfigureInterface(InterfaceSettings) error
}

type DummyConfigurator struct{}

func (*DummyConfigurator) ConfigureInterface(InterfaceSettings) error {
	return nil
}

type NetworkdConfigurator struct{}

func (nc *NetworkdConfigurator) ConfigureInterface(iface InterfaceSettings) error {
	var data string
	if iface.Mode == "dhcp" {
		data = fmt.Sprintf(`[Match]
Name=%s

[Network]
DHCP=ipv4
`, iface.Interface)
	} else {
		data = fmt.Sprintf(`[Match]
Name=%s

[Network]
Address=%s
`, iface.Interface, iface.Subnet)
		if iface.Gateway != "" {
			data += fmt.Sprintf("Gateway=%s\n", iface.Gateway)
		}
	}

	log.Debugf("Writing config file for %s", iface.Interface)
	if err := ioutil.WriteFile(fmt.Sprintf("/etc/systemd/network/%s.network", iface.Interface), []byte(data), 0666); err != nil {
		return errors.Wrapf(err, "could not write /etc/systemd/network/%s.network", iface.Interface)
	}

	log.Debugf("Restarting systemd-networkd")
	if err := restartNetworkd(); err != nil {
		return errors.Wrap(err, "could not restart systemd-networkd")
	}

	log.Debugf("Waiting for %s", iface.Interface)
	if err := waitForInterface(iface.Interface, 20*time.Second); err != nil {
		return errors.Wrapf(err, "error waiting for interface %s to come up", iface.Interface)
	}

	return nil
}

func restartNetworkd() error {
	cmd := exec.Command("/usr/bin/systemctl", "restart", "systemd-networkd")
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func waitForInterface(name string, timeout time.Duration) error {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return errors.Wrap(err, "error retrieving interface information")
	}

	start := time.Now()
	for time.Now().Before(start.Add(timeout)) {
		addrs, err := iface.Addrs()
		if err != nil {
			return errors.Wrap(err, "error retrieving interface addresses")
		}

		log.Debugf("Device has addresses: %v", addrs)

		num4Addrs := 0
		for _, addr := range addrs {
			var ip net.IP
			switch addr.(type) {
			case *net.IPAddr:
				ip = addr.(*net.IPAddr).IP.To4()
			case *net.IPNet:
				ip = addr.(*net.IPNet).IP.To4()
			}

			if ip != nil {
				num4Addrs++
			}
		}

		time.Sleep(3 * time.Second)

		if num4Addrs > 0 {
			// Done
			return nil
		}
	}

	return errors.Errorf("timed out waiting to obtain IP address")
}
