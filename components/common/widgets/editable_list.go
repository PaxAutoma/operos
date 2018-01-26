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

type EditableListItem struct {
	Label        string
	Key          string
	Value        string
	ValidateFunc func(string, string) []error
	Dirty        bool
}

func NewEditableListItem(label string, key string, value string, validateFunc func(string, string) []error) *EditableListItem {
	return &EditableListItem{label, key, value, validateFunc, false}
}

func (eli *EditableListItem) Validate() []error {
	var errors []error

	if eli.ValidateFunc != nil {
		errors = eli.ValidateFunc(eli.Label, eli.Value)
		for _, err := range errors {
			if ve, ok := err.(*ValidationError); ok {
				ve.Show = eli.Dirty
			}
		}
	}
	return errors
}

type EditableList struct {
	Name         string
	Bounds       image.Rectangle
	Items        []*EditableListItem
	OnItemChange func(item *EditableListItem)
	Dirty        bool

	visible    bool
	textboxes  []*Textbox
	labelWidth int
}

func NewEditableList(name string, items []*EditableListItem, width int, height int) *EditableList {
	widestLabel := 0
	for _, item := range items {
		if len(item.Label) > widestLabel {
			widestLabel = len(item.Label)
		}
	}

	labelWidth := widestLabel + 4

	textboxes := make([]*Textbox, len(items))
	for idx, item := range items {
		textboxes[idx] = NewTextbox(fmt.Sprintf("tb-%s-%s", name, item.Key), item.Value, false, width-labelWidth-2)
	}

	el := &EditableList{
		Name:       name,
		Bounds:     image.Rect(0, 0, width-1, height-1),
		Items:      items,
		visible:    true,
		labelWidth: labelWidth,
		textboxes:  textboxes,
	}

	return el
}

func (el *EditableList) Render(g *gocui.Gui, container image.Point, fs *FocusableSet) error {
	absPos := el.Bounds.Add(container)
	v, err := g.SetView(el.Name, absPos.Min.X, absPos.Min.Y, absPos.Max.X, absPos.Max.Y)
	if err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}

		for idx, item := range el.Items {
			item := item

			el.textboxes[idx].OnBlur = func(value string) {
				item.Value = value
				item.Dirty = true

				if el.OnItemChange != nil {
					el.OnItemChange(item)
				}
			}

			el.textboxes[idx].OnEnter = func(value string) {
				fs.Next()
			}
		}
	}

	v.Frame = el.visible
	v.Clear()

	for idx, item := range el.Items {
		if el.visible {
			fmt.Fprintln(v, item.Label)
		}

		tbMinPt := absPos.Min.Add(image.Point{el.labelWidth + 1, idx + 1})
		el.textboxes[idx].Render(g, tbMinPt, fs)
	}

	return nil
}

func (el *EditableList) GetHeight() int {
	return el.Bounds.Dy() + 1
}

func (el *EditableList) Validate() []error {
	errors := make([]error, 0)
	for _, item := range el.Items {
		errors = append(errors, item.Validate()...)
	}
	return errors
}

func (el *EditableList) SetVisibility(visible bool) {
	el.visible = visible
	for _, tb := range el.textboxes {
		tb.Visible = visible
	}
}
