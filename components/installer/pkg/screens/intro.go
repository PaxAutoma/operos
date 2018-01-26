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

func IntroScreen(screenSet *widgets.ScreenSet, context interface{}) *widgets.Screen {
	screen := widgets.NewScreen()
	screen.Title = "Welcome to the Pax Automa Operos Installer."
	screen.Message = "Please select one of the following options to get started:"
	screen.ShowPrev(false)
	screen.ShowNext(false)

	menu := widgets.NewMenu("menu-main", []widgets.MenuItem{
		&widgets.SimpleMenuItem{"Install", "install"},
		&widgets.SimpleMenuItem{"Exit to shell", "shell"},
		&widgets.SimpleMenuItem{"Reboot", "reboot"},
	}, 80, 5)
	menu.EnterToSelect = true

	menu.OnSelectItem = func(item widgets.MenuItem) error {
		smi := item.(*widgets.SimpleMenuItem)
		switch smi.Value {
		case "install":
			screenSet.Forward(1)
		case "reboot":
			return installer.Reboot()
		case "shell":
			return gocui.ErrQuit
		}

		return nil
	}

	screen.Content = menu
	return screen
}
