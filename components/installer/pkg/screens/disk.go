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
	"github.com/jroimartin/gocui"
	"github.com/paxautoma/operos/components/common/widgets"
	installer "github.com/paxautoma/operos/components/installer/pkg"
)

func DiskSelectionScreen(screenSet *widgets.ScreenSet, context interface{}) *widgets.Screen {
	ctx := context.(*installer.InstallerContext)

	screen := widgets.NewScreen()
	screen.Title = "Controller disk"
	screen.Message = `
Please select which disk should be used to install the controller software.
Only disks with at least 50GB are shown.

` +
		widgets.BoldString(gocui.ColorWhite, "Warning: all data on this disk will be deleted!")

	items := make([]widgets.MenuItem, len(ctx.Disks))
	selectedIdx := -1
	for idx, disk := range ctx.Disks {
		items[idx] = &widgets.SimpleMenuItem{
			Text:  disk.String(),
			Value: disk.Name,
		}
		if disk.Name == ctx.Responses.ControllerDisk {
			selectedIdx = idx
		}
	}

	menu := widgets.NewMenu("menu-disk", items, 80, 6)
	menu.OnSelectItem = func(item widgets.MenuItem) error {
		smi := item.(*widgets.SimpleMenuItem)
		ctx.Responses.ControllerDisk = "/dev/" + smi.Value
		return nil
	}
	if selectedIdx >= 0 {
		menu.SelectItem(selectedIdx)
	}

	screen.Content = menu

	screen.OnNext = func() error {
		screenSet.Forward(1)
		return nil
	}

	screen.OnPrev = func() error {
		screenSet.Back(1)
		return nil
	}

	return screen
}
