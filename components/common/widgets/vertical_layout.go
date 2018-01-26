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
	"image"

	"github.com/jroimartin/gocui"
)

type VerticalLayout struct {
	Items   []Renderable
	Padding int
}

func NewVerticalLayout() *VerticalLayout {
	return &VerticalLayout{Padding: 1}
}

func (vl *VerticalLayout) Render(g *gocui.Gui, container image.Point, fs *FocusableSet) error {
	pos := container
	for _, item := range vl.Items {
		if err := item.Render(g, pos, fs); err != nil {
			return err
		}
		pos.Y += item.GetHeight() + vl.Padding
	}

	return nil
}

func (vl *VerticalLayout) GetHeight() int {
	if len(vl.Items) < 1 {
		return 0
	}

	total := 0
	for _, item := range vl.Items {
		total += item.GetHeight()
	}

	total += vl.Padding * (len(vl.Items) - 1)
	return total
}
