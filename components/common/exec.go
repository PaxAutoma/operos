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

package common

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

type CmdExecutor struct {
	Name           string
	SuccessMessage string
	FailMessage    string
	OnFinish       func(success bool)

	cmd *exec.Cmd
}

func NewCmdExecutor(cmd *exec.Cmd) *CmdExecutor {
	return &CmdExecutor{
		cmd: cmd,
	}
}

func (ce *CmdExecutor) Start(w io.Writer) error {
	ce.cmd.Stdout = w
	ce.cmd.Stderr = w

	log.Debugf("Starting process %s with env", ce.cmd.Path, spew.Sdump(ce.cmd.Env))
	err := ce.cmd.Start()

	if err != nil {
		log.Debugf("Failed to start process: %s", err.Error())
		return err
	}

	go func() {
		err := ce.cmd.Wait()
		success := true

		if ioerr, ok := err.(*exec.ExitError); ok {
			success = ioerr.Success()
		} else if err != nil {
			panic(err)
		}

		if success {
			fmt.Fprintln(w, ce.SuccessMessage)
		} else {
			fmt.Fprintln(w, ce.FailMessage)
		}

		if ce.OnFinish != nil {
			ce.OnFinish(success)
		}
	}()

	return nil
}
