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

type MenuItem interface{}

type Menu struct {
	SimpleFocusable
	Bounds         image.Rectangle
	Items          []MenuItem
	OnSelectItem   func(item MenuItem) error
	RenderStrategy func(item MenuItem, selected bool, active bool, width int) string
	EnterToSelect  bool
	Visible        bool
	selectedIdx    int
	origin         int
}

type SimpleMenuItem struct {
	Text  string
	Value string
}

func BasicRenderStrategy(item MenuItem, selected bool, active bool, width int) string {
	smi := item.(*SimpleMenuItem)
	itemText := smi.Text + strings.Repeat(" ", width-len(smi.Text))

	switch {
	case selected && active:
		return ReverseString(itemText)
	case selected && !active:
		return ColorString(gocui.ColorGreen, itemText)
	default:
		return itemText
	}
}

func NewMenu(name string, items []MenuItem, width int, height int) *Menu {
	menu := &Menu{
		SimpleFocusable: SimpleFocusable{Name: name},
		Items:           items,
		Bounds:          image.Rect(0, 0, width-1, height-1),
		selectedIdx:     -1,
		origin:          0,
		RenderStrategy:  BasicRenderStrategy,
		Visible:         true,
	}

	return menu
}

func (m *Menu) Render(g *gocui.Gui, container image.Point, fs *FocusableSet) error {
	m.SimpleFocusable.g = g

	absPos := m.Bounds.Add(container)

	v, err := g.SetView(m.Name, absPos.Min.X, absPos.Min.Y, absPos.Max.X, absPos.Max.Y)
	if err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}

		g.SetKeybinding(m.Name, gocui.KeyArrowUp, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			m.MoveUp(1)
			return nil
		})

		g.SetKeybinding(m.Name, gocui.KeyArrowDown, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			m.MoveDown(1)
			return nil
		})

		g.SetKeybinding(m.Name, gocui.KeyPgup, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			m.MoveUp(m.Bounds.Dy() - 1)
			return nil
		})

		g.SetKeybinding(m.Name, gocui.KeyPgdn, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			m.MoveDown(m.Bounds.Dy() - 1)
			return nil
		})

		g.SetKeybinding(m.Name, gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			if m.EnterToSelect && len(m.Items) > 0 {
				if m.OnSelectItem != nil {
					return m.OnSelectItem(m.SelectedItem())
				}
			} else {
				fs.Next()
			}

			return nil
		})

		if fs != nil {
			fs.Add(m)
		}
	}

	v.Clear()

	if m.Visible {
		v.Frame = true

		active := g.CurrentView() == v

		// The first time this widget becomes active, select the first item
		if active && len(m.Items) > 0 && m.selectedIdx < 0 {
			m.selectedIdx = 0
			if !m.EnterToSelect && m.OnSelectItem != nil {
				m.OnSelectItem(m.SelectedItem())
			}
		}

		for idx, item := range m.Items {
			selected := (m.selectedIdx == idx) && (!m.EnterToSelect || active)
			label := m.RenderStrategy(item, selected, active, m.Bounds.Dx()-2)

			switch {
			case m.origin > 0 && idx == m.origin:
				label += "▲"
			case m.origin+m.Bounds.Dy()-1 < len(m.Items) && idx == m.origin+m.Bounds.Dy()-2:
				label += "▼"
			default:
				label += " "
			}

			fmt.Fprintln(v, label)
		}
		v.SetOrigin(0, m.origin)
	} else {
		v.Frame = false
	}

	return nil
}

func (menu *Menu) MoveUp(delta int) {
	if menu.selectedIdx > 0 {
		menu.selectedIdx -= delta
		if menu.selectedIdx < 0 {
			menu.selectedIdx = 0
		}

		if menu.origin > menu.selectedIdx {
			menu.origin = menu.selectedIdx
		}

		if menu.OnSelectItem != nil && !menu.EnterToSelect {
			menu.OnSelectItem(menu.SelectedItem())
		}
	}
}

func (menu *Menu) MoveDown(delta int) {
	if menu.selectedIdx < len(menu.Items)-1 {
		menu.selectedIdx += delta
		if menu.selectedIdx > len(menu.Items)-1 {
			menu.selectedIdx = len(menu.Items) - 1
		}

		if menu.origin+menu.Bounds.Dy()-2 < menu.selectedIdx {
			menu.origin = menu.selectedIdx - menu.Bounds.Dy() + 2
		}

		if menu.OnSelectItem != nil && !menu.EnterToSelect {
			menu.OnSelectItem(menu.SelectedItem())
		}
	}
}

func (m *Menu) WantsFocus() bool {
	return m.Visible
}

func (m *Menu) SelectItem(idx int) error {
	if idx >= 0 && idx < len(m.Items) {
		m.selectedIdx = idx
	}
	return nil
}

func (m *Menu) SelectedItem() MenuItem {
	if m.selectedIdx >= 0 {
		return m.Items[m.selectedIdx]
	}
	return nil
}

func (m *Menu) GetHeight() int {
	return m.Bounds.Dy() + 1
}
