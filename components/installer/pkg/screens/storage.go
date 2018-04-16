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
	"image"
	"strconv"

	"github.com/jroimartin/gocui"
	"github.com/paxautoma/operos/components/common/widgets"
	installer "github.com/paxautoma/operos/components/installer/pkg"
)

func StorageSettingsScreen(screenSet *widgets.ScreenSet, context interface{}) *widgets.Screen {
	ctx := context.(*installer.InstallerContext)

	screen := widgets.NewScreen()
	screen.Title = "Storage"
	screen.Message = `When workers boot from the Controller, their storage will be automatically
joined into the cluster. Each of the disks will be partitioned and divided into
two portions: system storage and data storage.

- System storage contains Docker images and file systems for running containers
- Data storage contains the distributed store used for persistent volumes

What percentage of the disk should be used for system storage?`

	tb := widgets.NewTextbox("tb-storage", "50", true, 80)

	errorList := widgets.NewPar("par-errors", "")
	errorList.Bounds = image.Rect(1, 0, 79, 3)
	errorList.FgColor = gocui.ColorRed

	var valid bool
	validate := func(value string) {
		errors := widgets.ValidateIntMinMax("Storage percentage", value, 20, 80)
		errorList.Text = widgets.JoinValidationErrors(errors)
		valid = len(errors) == 0
		screen.ShowNext(valid)
	}
	validate(strconv.Itoa(ctx.Responses.StorageSystemPercentage))

	tb.OnBlur = func(value string) {
		validate(value)

		if valid {
			if intValue, err := strconv.Atoi(value); err != nil {
				ctx.Responses.StorageSystemPercentage = intValue
			}
		}
	}

	tb.OnEnter = func(value string) {
		tb.OnBlur(value)
		if valid {
			screen.FocusableSet.Next()
		}
	}

	vl := widgets.NewVerticalLayout()
	vl.Items = []widgets.Renderable{tb, errorList}
	screen.Content = vl

	screen.OnPrev = func() error {
		screenSet.Back(1)
		return nil
	}

	screen.OnNext = func() error {
		screenSet.Forward(1)
		return nil
	}

	return screen
}
