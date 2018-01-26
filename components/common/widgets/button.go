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

	"github.com/jroimartin/gocui"
)

type Button struct {
	SimpleFocusable

	Label   string
	Bounds  image.Rectangle
	Handler func(string) error
	Visible bool
}

func NewButton(name string, label string, x int, y int, width int, height int) *Button {
	b := &Button{
		SimpleFocusable: SimpleFocusable{Name: name},
		Label:           label,
		Bounds:          image.Rect(x, y, x+width-1, y+height-1),
		Visible:         true,
	}
	return b
}

func (b *Button) Layout(g *gocui.Gui) error {
	return b.Render(g, image.Point{0, 0}, nil)
}

func (b *Button) Render(g *gocui.Gui, container image.Point, fs *FocusableSet) error {
	b.SimpleFocusable.g = g

	absPos := b.Bounds.Add(container)

	v, err := g.SetView(b.Name, absPos.Min.X, absPos.Min.Y, absPos.Max.X, absPos.Max.Y)
	if err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}

		g.SetKeybinding(b.Name, gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			if b.Handler != nil {
				return b.Handler(b.Name)
			}
			return nil
		})

		if fs != nil {
			fs.Add(b)
		}
	}

	v.Clear()
	if b.Visible {
		v.Frame = true
		fmt.Fprintf(v, CenterInBox(b.Label, b.Bounds.Dx(), b.Bounds.Dy()))
	} else {
		v.Frame = false
	}

	return nil
}

func (b *Button) WantsFocus() bool {
	return b.Visible
}

func (b *Button) GetHeight() int {
	return b.Bounds.Dy() + 1
}
