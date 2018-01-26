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
	"testing"

	"github.com/stretchr/testify/require"
)

type FakeFocusable struct {
	Active        bool
	DoesWantFocus bool
	Sequence      []string
}

func NewFakeFocusable(wantsFocus bool) *FakeFocusable {
	return &FakeFocusable{DoesWantFocus: wantsFocus}
}

func (o *FakeFocusable) WantsFocus() bool {
	return o.DoesWantFocus
}

func (o *FakeFocusable) Focus() {
	o.Sequence = append(o.Sequence, "focus")
	o.Active = true
}

func (o *FakeFocusable) Blur() {
	o.Sequence = append(o.Sequence, "blur")
	o.Active = false
}

func Test_Next_SingleItemWantsFocus_GetsFocus(t *testing.T) {
	focusableSet := NewFocusableSet()

	f1 := NewFakeFocusable(true)
	focusableSet.Add(f1)

	require.False(t, f1.Active)
	focusableSet.Next()
	require.True(t, f1.Active)
}

func Test_Next_SingleItemDoesNotWantFocus_DoesNotGetFocus(t *testing.T) {
	fs := NewFocusableSet()

	f1 := NewFakeFocusable(false)

	fs.Add(f1)
	fs.Next()

	require.False(t, f1.Active)
}

func Test_Next_NoItems(t *testing.T) {
	fs := NewFocusableSet()
	require.Equal(t, fs.GetCurrent(), nil)
	fs.Next()
	require.Equal(t, fs.GetCurrent(), nil)
}

func Test_Next_ManyItemsSuccessiveNext_ActivatesCorrectFocusable(t *testing.T) {
	focusableSet := NewFocusableSet()

	f1 := NewFakeFocusable(true)
	f2 := NewFakeFocusable(true)
	f3 := NewFakeFocusable(true)

	focusableSet.Add(f1)
	focusableSet.Add(f2)
	focusableSet.Add(f3)

	require.False(t, f1.Active)
	require.False(t, f2.Active)
	require.False(t, f3.Active)

	focusableSet.Next()
	require.True(t, f1.Active)
	require.False(t, f2.Active)
	require.False(t, f3.Active)

	focusableSet.Next()
	require.False(t, f1.Active)
	require.True(t, f2.Active)
	require.False(t, f3.Active)

	focusableSet.Next()
	require.False(t, f1.Active)
	require.False(t, f2.Active)
	require.True(t, f3.Active)
}

func Test_Next_SingleItemCalledMultipleTimes_CallsBlurAndFocusInOrder(t *testing.T) {
	focusableSet := NewFocusableSet()
	f1 := NewFakeFocusable(true)
	focusableSet.Add(f1)

	require.False(t, f1.Active)
	focusableSet.Next()
	require.True(t, f1.Active)

	require.Equal(t, f1.Sequence, []string{
		"focus",
	})
	f1.Sequence = nil

	focusableSet.Next()
	require.Equal(t, f1.Sequence, []string{
		"blur",
		"focus",
	})
	f1.Sequence = nil
}

func Test_Next_CalledWhenLastItemActive_WrapsAround(t *testing.T) {
	focusableSet := NewFocusableSet()

	f1 := NewFakeFocusable(true)
	f2 := NewFakeFocusable(true)
	f3 := NewFakeFocusable(true)

	focusableSet.Add(f1)
	focusableSet.Add(f2)
	focusableSet.Add(f3)

	focusableSet.Next()
	focusableSet.Next()
	focusableSet.Next()
	focusableSet.Next()

	require.True(t, f1.Active)
	require.False(t, f2.Active)
	require.False(t, f3.Active)
}

func Test_Next_ItemsThatDontWantFocus_SkipsThoseItems(t *testing.T) {
	focusableSet := NewFocusableSet()

	f1 := NewFakeFocusable(true)
	f2 := NewFakeFocusable(false)
	f3 := NewFakeFocusable(true)
	f4 := NewFakeFocusable(false)
	f5 := NewFakeFocusable(false)
	f6 := NewFakeFocusable(true)
	focusableSet.Add(f1)
	focusableSet.Add(f2)
	focusableSet.Add(f3)
	focusableSet.Add(f4)
	focusableSet.Add(f5)
	focusableSet.Add(f6)

	focusableSet.Next()
	focusableSet.Next()

	require.False(t, f1.Active)
	require.False(t, f2.Active)
	require.True(t, f3.Active)
	require.False(t, f4.Active)
	require.False(t, f5.Active)
	require.False(t, f6.Active)

	focusableSet.Next()
	require.False(t, f1.Active)
	require.False(t, f2.Active)
	require.False(t, f3.Active)
	require.False(t, f4.Active)
	require.False(t, f5.Active)
	require.True(t, f6.Active)
}

func Test_GetCurrent_InitiallyNoItemIsActive_ReturnsNil(t *testing.T) {
	focusableSet := NewFocusableSet()
	focusableSet.Add(NewFakeFocusable(true))

	require.Nil(t, focusableSet.GetCurrent())
}

func Test_GetCurrent_Always_ReturnsCurrentItem(t *testing.T) {
	focusableSet := NewFocusableSet()

	f1 := NewFakeFocusable(true)
	f2 := NewFakeFocusable(true)
	f3 := NewFakeFocusable(true)

	focusableSet.Add(f1)
	focusableSet.Add(f2)
	focusableSet.Add(f3)

	focusableSet.Next()
	require.Equal(t, focusableSet.GetCurrent(), f1)

	focusableSet.Next()
	require.Equal(t, focusableSet.GetCurrent(), f2)

	focusableSet.Next()
	require.Equal(t, focusableSet.GetCurrent(), f3)

	focusableSet.Next()
	require.Equal(t, focusableSet.GetCurrent(), f1)
}
