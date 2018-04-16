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
	"sort"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/paxautoma/operos/components/common/widgets"
	installer "github.com/paxautoma/operos/components/installer/pkg"
	"github.com/paxautoma/operos/components/installer/pkg/network"
	log "github.com/sirupsen/logrus"
)

func HardwareProbeScreen(screenSet *widgets.ScreenSet, context interface{}) *widgets.Screen {
	ctx := context.(*installer.InstallerContext)

	screen := widgets.NewScreen()
	screen.Message = widgets.CenterInBox(`Please wait, probing hardware.
This can take up to 15 seconds`, 80, 25)
	screen.ShowNext(false)
	screen.ShowPrev(false)

	screen.OnInitialize = func(g *gocui.Gui) {
		go func() {
			ifaces, err := network.ProbePhysicalInterfaces()
			if err != nil {
				panic(err)
			}

			disks, err := installer.GetDiskList()
			if err != nil {
				panic(err)
			}

			goodDisks := []installer.DiskInfo{}
			for _, disk := range disks {
				if disk.Size >= 50*1024*1024*1024 {
					goodDisks = append(goodDisks, disk)
				}
			}

			numCPUs := installer.GetNumCPUs()

			totalRAM, err := installer.GetTotalMemory()
			if err != nil {
				panic(err)
			}

			log.Debugf("CPUs: %d, RAM: %d GiB", numCPUs, totalRAM)

			time.Sleep(1 * time.Second)

			g.Update(func(g *gocui.Gui) error {
				hwErrors := []string{}

				if len(ifaces) < 1 {
					hwErrors = append(hwErrors, fmt.Sprintf(
						"Required: at least 1 physical wired interface. This machine has %d.",
						len(ifaces)))
				}

				if numCPUs < 2 {
					hwErrors = append(hwErrors, fmt.Sprintf("Required: at least 2 CPUs. This machine has %d.", numCPUs))
				}

				if totalRAM < 2 {
					hwErrors = append(hwErrors, fmt.Sprintf("Required: at least 2 GB of RAM. This machine has %d GB.", totalRAM))
				}

				if len(goodDisks) < 1 {
					hwErrors = append(hwErrors, fmt.Sprintf("Required: a disk with at least 50GB capacity."))
				}

				if len(hwErrors) > 0 {
					screen.Title = "Hardware requirements"
					screen.Message = fmt.Sprintf(`
 This machine does not meet the hardware requirements for Operos Controller:
 
 - %s
 
 Installation cannot continue.`, strings.Join(hwErrors, "\n - "))
					screen.ShowPrev(true)
					screen.FocusableSet.Next()
					return nil
				}

				ctx.Interfaces.ByName = ifaces

				names := make([]string, len(ifaces))
				idx := 0
				for name := range ifaces {
					names[idx] = name
					idx++
				}
				sort.Strings(names)

				ctx.Interfaces.Ordered = make([]*network.InterfaceInfo, len(names))
				for idx, name := range names {
					ctx.Interfaces.Ordered[idx] = ifaces[name]
				}

				ctx.Disks = goodDisks

				screenSet.Forward(1)
				return nil
			})
		}()
	}

	screen.OnPrev = func() error {
		screenSet.Back(1)
		return nil
	}

	return screen
}
