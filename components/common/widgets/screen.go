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

type Renderable interface {
	Render(*gocui.Gui, image.Point, *FocusableSet) error
	GetHeight() int
}

type Screen struct {
	Title        string
	Message      string
	OnPrev       func() error
	OnNext       func() error
	OnInitialize func(g *gocui.Gui)
	Content      Renderable
	FocusableSet *FocusableSet

	buttonNext *Button
	buttonPrev *Button
	g          *gocui.Gui
}

func NewScreen() *Screen {
	s := &Screen{}
	s.buttonNext = NewButton("button-next", "Next", 71, 22, 8, 3)
	s.buttonNext.Handler = func(string) error {
		if s.OnNext != nil {
			return s.OnNext()
		}
		return nil
	}

	s.buttonPrev = NewButton("button-prev", "Back", 1, 22, 8, 3)
	s.buttonPrev.Handler = func(string) error {
		if s.OnPrev != nil {
			return s.OnPrev()
		}
		return nil
	}

	return s
}

func (s *Screen) Layout(g *gocui.Gui) error {
	s.g = g

	termWidth, termHeight := g.Size()
	bounds := image.Rect(termWidth/2-41, termHeight/2-13, termWidth/2+40, termHeight/2+13)

	v, err := g.SetView("outer-frame", bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Max.Y)
	initialize := false
	if err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}

		initialize = true
	}

	inner := bounds.Inset(1)
	contentStart := inner.Min

	if initialize {
		s.FocusableSet = NewFocusableSet()
		v.Frame = true
	}

	v.Clear()

	if s.Title != "" {
		fmt.Fprintln(v, BoldString(gocui.ColorWhite, s.Title))
		fmt.Fprintln(v)
		contentStart.Y += 2
	}

	if s.Message != "" {
		fmt.Fprintln(v, s.Message)
		contentStart.Y += strings.Count(s.Message, "\n") + 2
	}

	if s.Content != nil {
		err = s.Content.Render(g, contentStart, s.FocusableSet)
		if err != nil {
			return err
		}
	}

	s.buttonNext.Render(g, inner.Min, s.FocusableSet)
	s.buttonPrev.Render(g, inner.Min, s.FocusableSet)

	if initialize {
		s.FocusableSet.Next()

		s.Focus()

		g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			return gocui.ErrQuit
		})

		if s.OnInitialize != nil {
			s.OnInitialize(g)
		}
	}

	return nil
}

func (s *Screen) ShowNext(enabled bool) {
	s.buttonNext.Visible = enabled
}

func (s *Screen) ShowPrev(enabled bool) {
	s.buttonPrev.Visible = enabled
}

func (s *Screen) Focus() {
	if s.g == nil {
		return
	}

	if curItem := s.FocusableSet.GetCurrent(); curItem != nil {
		curItem.Focus()
	}

	s.g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		s.FocusableSet.Next()
		return nil
	})
}

func (s *Screen) Blur() {
	if s.g == nil {
		return
	}

	s.g.DeleteKeybinding("", gocui.KeyTab, gocui.ModNone)
}

type ScreenCreator func(screenSet *ScreenSet, context interface{}) *Screen

type ScreenSet struct {
	Screens   []ScreenCreator
	g         *gocui.Gui
	activeIdx int
	context   interface{}
}

func NewScreenSet(g *gocui.Gui, context interface{}) *ScreenSet {
	return &ScreenSet{
		g:       g,
		context: context,
	}
}

func (ss *ScreenSet) Forward(skip int) {
	ss.move(skip)
}

func (ss *ScreenSet) Back(skip int) {
	ss.move(-skip)
}

func (ss *ScreenSet) Restart() {
	ss.move(-ss.activeIdx)
}

func (ss *ScreenSet) move(skip int) {
	newIdx := ss.activeIdx + skip

	if newIdx >= 0 && newIdx < len(ss.Screens) {
		ss.activeIdx = newIdx
		ss.g.SetManager(ss.Screens[ss.activeIdx](ss, ss.context))
	}
}

func (ss *ScreenSet) Start() {
	if len(ss.Screens) < 1 {
		return
	}

	ss.g.SetManager(ss.Screens[ss.activeIdx](ss, ss.context))
}
