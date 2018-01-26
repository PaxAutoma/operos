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

package statustty

import (
	"io/ioutil"
	"net"
	"strings"
	"time"

	"github.com/paxautoma/operos/components/common"

	"github.com/pkg/errors"
)

type NetStatus struct {
	Statuses []IfaceStatus
	Error    error
}

func (ns *NetStatus) IsOK() bool {
	for _, v := range ns.Statuses {
		if v.IP == nil || !v.Up {
			return false
		}
	}
	return true
}

type IfaceSpec struct {
	Title  string
	Device string
}

type IfaceStatus struct {
	Spec IfaceSpec
	Up   bool
	IP   *net.IPNet
}

func (is IfaceStatus) NiceIP() string {
	if is.IP != nil {
		return is.IP.String()
	}
	return "-"
}

func SubscribeNetStatus(ifaces []IfaceSpec, closer <-chan struct{}) <-chan *NetStatus {
	ch := make(chan *NetStatus)

	go func() {
		defer common.LogPanic()

		ch <- GetNetStatus(ifaces)

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case _, ok := <-closer:
				if !ok {
					break
				}

			case <-ticker.C:
				ch <- GetNetStatus(ifaces)
			}
		}
	}()

	return ch
}

func GetNetStatus(ifaces []IfaceSpec) *NetStatus {
	result := make([]IfaceStatus, len(ifaces))

	for idx, iface := range ifaces {
		status, err := GetIfaceStatus(iface)
		if err != nil {
			return &NetStatus{Error: errors.Wrapf(err, "failed to get status for interface %s", iface)}
		}

		result[idx] = status
	}

	return &NetStatus{Statuses: result}
}

func GetIfaceStatus(spec IfaceSpec) (IfaceStatus, error) {
	iface, err := net.InterfaceByName(spec.Device)
	if err != nil {
		return IfaceStatus{}, errors.Wrapf(err, "could not lookup interface")
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return IfaceStatus{}, errors.Wrap(err, "could not obtain interface addresses")
	}

	status := IfaceStatus{
		Spec: spec,
		Up:   iface.Flags&net.FlagUp > 0,
	}
	for _, addr := range addrs {
		if ip, ok := addr.(*net.IPNet); ok && ip.IP.To4() != nil {
			status.IP = ip
			break
		}
	}

	return status, nil
}

func GetBootIface() (string, error) {
	cmdline, err := ioutil.ReadFile("/proc/cmdline")
	if err != nil {
		return "", errors.Wrap(err, "could not read /proc/cmdline")
	}

	return GetBootIfaceFromCmdLine(string(cmdline))
}

func GetBootIfaceFromCmdLine(cmdline string) (string, error) {
	args := strings.Split(cmdline, " ")
	for _, arg := range args {
		arg = strings.TrimSpace(arg)
		if arg == "" || !strings.Contains(arg, "=") {
			continue
		}

		parts := strings.SplitN(arg, "=", 2)
		k, v := parts[0], parts[1]

		if k == "BOOTIF" {
			if len(v) != 20 {
				return "", errors.Errorf("kernel argument BOOTIF has invalid format (length %d != 20)", len(v))
			}

			mac := strings.ToLower(strings.Replace(v[3:], "-", ":", -1))

			ifaces, err := net.Interfaces()
			if err != nil {
				return "", errors.Wrapf(err, "failed to list interfaces")
			}

			for _, iface := range ifaces {
				if iface.HardwareAddr.String() == mac {
					return iface.Name, nil
				}
			}

			return "", errors.Errorf("interface matching BOOTIF not found")
		}
	}

	return "", errors.Errorf("kernel commandline does not contain BOOTIF argument")
}
