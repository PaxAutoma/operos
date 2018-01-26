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

package widgets

import (
	"fmt"
	"image"
	"strings"

	"github.com/jroimartin/gocui"
)

type Dialog struct {
	Name       string
	Message    string
	AcceptText string
	CancelText string
	Width      int
	OnAccept   func() error
	OnCancel   func() error

	g *gocui.Gui
}

func NewDialog(name string, message string, width int) *Dialog {
	return &Dialog{
		Name:       name,
		Message:    message,
		AcceptText: "Yes",
		CancelText: "No",
		Width:      width,
	}
}

func (d *Dialog) Layout(g *gocui.Gui) error {
	d.g = g

	messageHeight := strings.Count(strings.TrimSpace(d.Message), "\n") + 1
	height := messageHeight + 6
	maxWidth, maxHeight := g.Size()

	if v, err := g.SetView(d.Name, (maxWidth-d.Width)/2, (maxHeight-height)/2, (maxWidth+d.Width)/2, (maxHeight+height)/2); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}

		v.Frame = true
		fmt.Fprint(v, "\n", d.Message)

		fs := NewFocusableSet()

		butAccept := NewButton(fmt.Sprintf("%s-accept", d.Name), d.AcceptText, 0, 0, len(d.AcceptText)+4, 3)
		butAccept.Render(g, image.Point{(maxWidth-d.Width)/2 + 2, (maxHeight+height)/2 - 3}, fs)
		butAccept.Handler = func(string) error {
			if d.OnAccept != nil {
				return d.OnAccept()
			}
			return nil
		}

		butCancel := NewButton(fmt.Sprintf("%s-cancel", d.Name), d.CancelText, 0, 0, len(d.CancelText)+4, 3)
		butCancel.Render(g, image.Point{(maxWidth-d.Width)/2 + len(d.CancelText) + 8, (maxHeight+height)/2 - 3}, fs)
		butCancel.Handler = func(string) error {
			if d.OnCancel != nil {
				return d.OnCancel()
			}
			return nil
		}

		fs.Next()

		g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			fs.Next()
			return nil
		})
	}

	return nil
}

func (d *Dialog) Close() {
	if d.g == nil {
		return
	}

	d.g.DeleteView(fmt.Sprintf("%s-accept", d.Name))
	d.g.DeleteView(fmt.Sprintf("%s-cancel", d.Name))
	d.g.DeleteView(d.Name)
	d.g.DeleteKeybinding("", gocui.KeyTab, gocui.ModNone)
}
