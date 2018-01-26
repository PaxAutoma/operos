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

	"github.com/paxautoma/operos/components/common/widgets"
	installer "github.com/paxautoma/operos/components/installer/pkg"
)

func PasswordScreen(screenSet *widgets.ScreenSet, context interface{}) *widgets.Screen {
	ctx := context.(*installer.InstallerContext)

	screen := widgets.NewScreen()
	screen.Title = "Root password"
	screen.Message = `
This password will be required to log into the UI and the console. Please type
it twice to avoid mistakes.
`
	screen.ShowNext(false)

	pwdLabel1 := widgets.NewPar("pwd-1", "Password:")
	pwdLabel1.Bounds = image.Rect(0, 0, 20, 0)
	pwdTextBox1 := widgets.NewTextbox("pwd-txt-1", "", true, 80)
	pwdTextBox1.Mask = '*'

	pwdLabel2 := widgets.NewPar("pwd-2", "Confirm password:")
	pwdLabel2.Bounds = image.Rect(0, 0, 20, 0)
	pwdTextBox2 := widgets.NewTextbox("pwd-txt-2", "", true, 80)
	pwdTextBox2.Mask = '*'

	check := func(_ string) {
		if pwdTextBox1.Value == pwdTextBox2.Value && pwdTextBox1.Value != "" && pwdTextBox2.Value != "" {
			screen.ShowNext(true)
		} else {
			screen.ShowNext(false)
		}
	}

	pwdTextBox1.OnChange = check
	pwdTextBox2.OnChange = check
	pwdTextBox1.OnEnter = func(value string) {
		screen.FocusableSet.Next()
	}
	pwdTextBox2.OnEnter = func(value string) {
		screen.FocusableSet.Next()
	}

	screen.OnPrev = func() error {
		screenSet.Back(1)
		return nil
	}

	screen.OnNext = func() error {
		ctx.Responses.RootPassword = pwdTextBox1.Value
		screenSet.Forward(1)
		return nil
	}

	vl := widgets.NewVerticalLayout()
	vl.Items = []widgets.Renderable{pwdLabel1, pwdTextBox1, pwdLabel2, pwdTextBox2}
	screen.Content = vl

	return screen
}
