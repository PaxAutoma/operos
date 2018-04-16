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

type Par struct {
	SimpleFocusable

	Text      string
	Bounds    image.Rectangle
	FgColor   gocui.Attribute
	BgColor   gocui.Attribute
	Focusable bool
}

func NewPar(name string, text string) *Par {
	return &Par{
		SimpleFocusable: SimpleFocusable{name, nil},
		Text:            text,
		FgColor:         gocui.ColorDefault,
		BgColor:         gocui.ColorDefault,
	}
}

func (p *Par) Render(g *gocui.Gui, container image.Point, fs *FocusableSet) error {
	absBounds := p.Bounds.Add(container)

	v, err := g.SetView(p.Name, absBounds.Min.X-1, absBounds.Min.Y-1, absBounds.Max.X+1, absBounds.Max.Y+1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}

		p.g = g

		if fs != nil {
			fs.Add(p)
		}

		v.Clear()
		fmt.Fprint(v, p.Text)
		v.SetOrigin(0, 0)

		moveOriginY := func(dy int) {
			numLines := strings.Count(p.Text, "\n") + 1
			ox, oy := v.Origin()
			newY := oy + dy
			if newY > numLines-(p.Bounds.Dy()-1) {
				newY = numLines - (p.Bounds.Dy() - 1)
			}
			if newY < 0 {
				newY = 0
			}
			v.SetOrigin(ox, newY)
		}

		p.g.SetKeybinding(p.Name, gocui.KeyArrowDown, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			moveOriginY(1)
			return nil
		})

		p.g.SetKeybinding(p.Name, gocui.KeyPgdn, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			moveOriginY(p.Bounds.Dy())
			return nil
		})

		p.g.SetKeybinding(p.Name, gocui.KeyArrowUp, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			moveOriginY(-1)
			return nil
		})

		p.g.SetKeybinding(p.Name, gocui.KeyPgup, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			moveOriginY(-p.Bounds.Dy())
			return nil
		})
	}

	v.FgColor = p.FgColor
	v.BgColor = p.BgColor
	v.Frame = p.Focusable
	v.Wrap = false
	v.Autoscroll = !p.Focusable

	return nil
}

func (p *Par) Destroy(g *gocui.Gui) {
	g.DeleteView(p.Name)
}

func (p *Par) GetHeight() int {
	return p.Bounds.Dy()
}

func (p *Par) WantsFocus() bool {
	return p.Focusable
}

type execResult struct {
	n   int
	err error
}

func (p *Par) Write(data []byte) (n int, err error) {
	v, err := p.g.View(p.Name)
	if err != nil {
		return 0, err
	}

	ch := make(chan execResult)

	p.g.Update(func(g *gocui.Gui) error {
		n, err := v.Write(data)
		ch <- execResult{n, err}

		return nil
	})

	result := <-ch
	return result.n, result.err
}

func (p *Par) SetText(text string) {
	p.Text = text

	if p.g == nil {
		return
	}

	p.g.Update(func(g *gocui.Gui) error {
		v, err := g.View(p.Name)
		if err != nil {
			return err
		}

		v.Clear()
		fmt.Fprint(v, p.Text)

		return nil
	})
}
