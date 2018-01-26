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

type Textbox struct {
	SimpleFocusable
	Key         string
	Frame       bool
	OnEnter     func(value string)
	OnBlur      func(value string)
	OnChange    func(value string)
	OnArrowUp   func()
	OnArrowDown func()
	Value       string
	Visible     bool
	Mask        rune
	initValue   string
	width       int
}

func NewTextbox(name string, value string, frame bool, width int) *Textbox {
	return &Textbox{
		SimpleFocusable: SimpleFocusable{Name: name},
		Frame:           frame,
		Value:           value,
		Visible:         true,
		initValue:       value,
		width:           width,
	}
}

func (tb *Textbox) Render(g *gocui.Gui, container image.Point, fs *FocusableSet) error {
	tb.SimpleFocusable.g = g

	bounds := image.Rect(-1, -1, tb.width, 1).Add(container)
	if tb.Frame {
		bounds = bounds.Add(image.Point{1, 1})
		bounds.Max.X -= 2
	}

	if tb.Visible {
		v, err := g.SetView(tb.Name, bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Max.Y)
		if err != nil {
			if err != gocui.ErrUnknownView {
				panic(err)
			}

			v.Frame = tb.Frame
			v.Editable = true
			v.Wrap = false
			v.Editor = gocui.EditorFunc(func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
				tb.editor(g, v, key, ch, mod)
			})

			v.Mask = tb.Mask

			v.Clear()
			fmt.Fprintf(v, tb.Value)

			fs.Add(tb)
		}
	} else {
		if err := g.DeleteView(tb.Name); err != nil && err != gocui.ErrUnknownView {
			panic(err)
		}
	}

	return nil
}

func (tb *Textbox) GetHeight() int {
	if tb.Frame {
		return 3
	}
	return 1
}

func (tb *Textbox) Destroy(g *gocui.Gui) {
	err := g.DeleteView(tb.Name)
	if err != nil {
		panic(err)
	}
	g.Cursor = false
}

func (tb *Textbox) editor(g *gocui.Gui, v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	cx, _ := v.Cursor()
	ox, _ := v.Origin()
	value := strings.TrimRight(v.Buffer(), "\n")

	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	case key == gocui.KeyDelete:
		v.EditDelete(false)
	case key == gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
	case key == gocui.KeyEnter:
		if tb.OnEnter != nil {
			tb.OnEnter(value)
		}
	case key == gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0, false)
	case key == gocui.KeyArrowRight:
		if cx+ox < len(value) {
			v.MoveCursor(1, 0, false)
		}
	case key == gocui.KeyArrowUp:
		if tb.OnArrowUp != nil {
			tb.OnArrowUp()
		}
	case key == gocui.KeyArrowDown:
		if tb.OnArrowDown != nil {
			tb.OnArrowDown()
		}
	case key == gocui.KeyHome:
		v.MoveCursor(-cx-ox, 0, false)
	case key == gocui.KeyEnd:
		v.MoveCursor(len(value)-ox-cx, 0, false)
	}

	tb.Value = strings.TrimRight(v.Buffer(), "\n")
	if tb.Value != value && tb.OnChange != nil {
		tb.OnChange(value)
	}
}

func (tb *Textbox) WantsFocus() bool {
	return tb.Visible
}

func (tb *Textbox) Focus() {
	tb.g.Cursor = true
	tb.SimpleFocusable.Focus()

	v, err := tb.g.View(tb.Name)
	if err != nil {
		panic(err)
	}

	maxX, _ := v.Size()
	cx := len(tb.Value)
	ox := 0
	if cx >= maxX {
		ox = cx - maxX + 1
		cx = maxX - 1
	}
	v.SetCursor(cx, 0)
	v.SetOrigin(ox, 0)
}

func (tb *Textbox) Blur() {
	tb.g.Cursor = false
	tb.SimpleFocusable.Blur()

	v, err := tb.g.View(tb.Name)
	if err != nil {
		panic(err)
	}

	v.SetOrigin(0, 0)

	if tb.OnBlur != nil {
		tb.OnBlur(tb.Value)
	}
}
