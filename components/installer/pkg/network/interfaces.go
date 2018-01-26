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
	"net"
	"os"
	"os/exec"
	"path"
	"strconv"
	"time"

	"github.com/d2g/dhcp4"
	"github.com/pkg/errors"
	"github.com/rlisagor/dhcp4client"
	log "github.com/sirupsen/logrus"
)

const sysfsnetPath = "/sys/class/net"

func GetPhysicalInterfaces() ([]net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, errors.Wrap(err, "could not list interfaces")
	}

	result := make([]net.Interface, 0)
	for _, iface := range ifaces {
		isPhysical, err := checkSysfsNetFileExists(iface.Name, "device")
		if err != nil {
			return nil, errors.Wrapf(err, "could not determine whether interface %s is virtual", iface.Name)
		}
		if !isPhysical {
			log.Debugf("hiding interface %s because it's virtual", iface.Name)
			continue
		}

		isWifi, err := checkSysfsNetFileExists(iface.Name, "phy80211")
		if err != nil {
			return nil, errors.Wrapf(err, "could not determine whether interface %s is wireless", iface.Name)
		}
		if isWifi {
			log.Debugf("hiding interface %s because it's wireless", iface.Name)
			continue
		}

		isWifi, err = checkSysfsNetFileExists(iface.Name, "wireless")
		if err != nil {
			return nil, errors.Wrapf(err, "could not determine whether interface %s is wireless", iface.Name)
		}
		if isWifi {
			log.Debugf("hiding interface %s because it's wireless", iface.Name)
			continue
		}

		result = append(result, iface)
	}

	return result, nil
}

func checkSysfsNetFileExists(ifname, subpath string) (bool, error) {
	filePath := path.Join(sysfsnetPath, ifname, subpath)
	if _, err := os.Stat(filePath); err != nil {
		if !os.IsNotExist(err) {
			return false, errors.Wrapf(err, "failed to stat file %s", filePath)
		}
		return false, nil
	}
	return true, nil
}

func ProbeInterfaceDhcp(iface net.Interface, timeout time.Duration) (dhcp4.Packet, error) {
	sock, err := dhcp4client.NewPacketSock(iface.Index)
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize DHCP client")
	}

	opts := []func(*dhcp4client.Client) error{
		dhcp4client.Connection(sock),
		dhcp4client.Timeout(timeout),
		dhcp4client.HardwareAddr(iface.HardwareAddr),
	}

	client, err := dhcp4client.New(opts...)
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize DHCP client")
	}

	discoverPacket, err := client.SendDiscoverPacket()
	if err != nil {
		return nil, errors.Wrap(err, "could not send DHCP discover packet")
	}

	offerPacket, err := client.GetOffer(&discoverPacket)
	if err != nil {
		return nil, errors.Wrap(err, "unable to receive DHCP offer")
	}

	return offerPacket, nil
}

func ParseOffer(offerPacket dhcp4.Packet) string {
	options := offerPacket.ParseOptions()
	mask, maskExists := options[dhcp4.OptionCode(1)]

	result := offerPacket.YIAddr().String()

	if maskExists {
		ones, _ := net.IPMask(mask).Size()
		result += "/" + strconv.Itoa(ones)
	}

	return result
}

type InterfaceInfo struct {
	Name      string
	Mac       string
	DhcpOffer string
	Err       error
}

type Timeoutable interface {
	Timeout() bool
}

func ProbePhysicalInterfaces() (map[string]*InterfaceInfo, error) {
	ifaces, err := GetPhysicalInterfaces()
	if err != nil {
		return nil, err
	}

	ch := make(chan InterfaceInfo)

	for _, iface := range ifaces {
		go func(iface net.Interface) {
			ch <- ProbeInterface(iface)
		}(iface)
	}

	output := make(map[string]*InterfaceInfo)
	for range ifaces {
		ifaceInfo := <-ch
		output[ifaceInfo.Name] = &ifaceInfo
	}

	return output, nil
}

func BringUpInterface(name string) error {
	log.Debugf("Bringing up interface %s", name)

	cmd := exec.Command("ip", "link", "set", name, "up")
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return errors.Wrapf(err, "ip link command failed (%s)", exitErr.Stderr)
		}
		return errors.Wrap(err, "ip link command failed")
	}

	time.Sleep(3 * time.Second)
	return nil
}

func ProbeInterface(iface net.Interface) InterfaceInfo {
	log.Debugf("Probing DHCP on %s", iface.Name)
	info := InterfaceInfo{
		Name: iface.Name,
		Mac:  iface.HardwareAddr.String(),
	}

	if (iface.Flags & net.FlagUp) == 0 {
		if err := BringUpInterface(iface.Name); err != nil {
			info.Err = err
			return info
		}
	}

	offerPacket, err := ProbeInterfaceDhcp(iface, 5*time.Second)
	if err != nil {
		log.Debugf("Probe on %s: error (%s)", iface.Name, err.Error())
		terr, ok := err.(Timeoutable)
		if !ok || !terr.Timeout() {
			info.Err = err
		}
	} else {
		info.DhcpOffer = ParseOffer(offerPacket)
		log.Debugf("Probe on %s: offer (%s)", iface.Name, info.DhcpOffer)
	}

	return info
}
