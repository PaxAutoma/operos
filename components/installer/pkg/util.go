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

package installer

import (
	"os/exec"

	"github.com/jroimartin/gocui"
)

type userError interface {
	IsUserError() bool
}

type temporary interface {
	Temporary() bool
}

func IsUserError(err error) bool {
	if userErr, ok := err.(userError); ok {
		return userErr.IsUserError()
	}
	return false
}

func IsTemporaryError(err error) bool {
	if tempErr, ok := err.(temporary); ok {
		return tempErr.Temporary()
	}
	return false
}

func Reboot() error {
	cmd := exec.Command("/sbin/reboot")
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	return gocui.ErrQuit
}
