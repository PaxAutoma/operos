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
	"github.com/jroimartin/gocui"
)

type Focusable interface {
	WantsFocus() bool
	Focus()
	Blur()
}

type SimpleFocusable struct {
	Name string

	g *gocui.Gui
}

func (sf *SimpleFocusable) WantsFocus() bool {
	return true
}

func (sf *SimpleFocusable) Focus() {
	if sf.g != nil {
		sf.g.SetCurrentView(sf.Name)
	}
}

func (sf *SimpleFocusable) Blur() {

}

type FocusableSet struct {
	Items   []Focusable
	Current int
}

func NewFocusableSet() *FocusableSet {
	return &FocusableSet{Current: -1}
}

func (fs *FocusableSet) Add(f Focusable) {
	found := false
	for _, item := range fs.Items {
		if item == f {
			found = true
		}
	}
	if !found {
		fs.Items = append(fs.Items, f)
	}
}

func (fs *FocusableSet) Next() {
	fs.move(true)
}

func (fs *FocusableSet) Prev() {
	fs.move(false)
}

func (fs *FocusableSet) move(forward bool) {
	step := 1
	if !forward {
		step = -1
	}

	if cur := fs.GetCurrent(); cur != nil {
		cur.Blur()
	}

	numTries := 0
	candidate := fs.Current

	for numTries < len(fs.Items) {
		if candidate+step >= len(fs.Items) {
			candidate = 0
		} else if candidate+step < 0 {
			candidate = len(fs.Items) - 1
		} else {
			candidate += step
		}

		if fs.Items[candidate].WantsFocus() {
			fs.Current = candidate
			break
		}

		numTries++
	}

	if cur := fs.GetCurrent(); cur != nil {
		cur.Focus()
	}
}

func (fs *FocusableSet) GetCurrent() Focusable {
	if len(fs.Items) < 1 || fs.Current < 0 {
		return nil
	}
	return fs.Items[fs.Current]
}
