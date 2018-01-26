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

package statustty

import (
	"fmt"
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetBootIfaceFromCmdLine(t *testing.T) {
	errorTests := []struct {
		name      string
		cmdline   string
		errString string
	}{
		{
			"NoArg_ReturnsError",
			"foo bar blah yahoo=hello foofoo=yes",
			"does not contain BOOT_IF",
		},
		{
			"NoValue_ReturnsError",
			"foo bar blah yahoo=hello foofoo=yes BOOTIF=",
			"invalid format",
		},
		{
			"ValueTooShort_ReturnsError",
			"foo bar blah yahoo=hello foofoo=yes BOOTIF=00-11-22-33-44-55",
			"invalid format",
		},
		{
			"NotRealMac_ReturnsError",
			"foo bar blah yahoo=hello foofoo=yes BOOTIF=01-00-11-22-33-44-55",
			"not found",
		},
	}

	for _, test := range errorTests {
		t.Run(test.name, func(t *testing.T) {
			_, err := GetBootIfaceFromCmdLine(test.cmdline)
			require.Error(t, err)
			require.Contains(t, err.Error(), test.errString)
		})
	}

	ifaces, err := net.Interfaces()
	require.NoError(t, err)
	realMac := "01-" + strings.Replace(ifaces[1].HardwareAddr.String(), ":", "-", -1)
	realIfName := ifaces[1].Name

	goodTests := []struct {
		name    string
		cmdline string
	}{
		{
			"RealMac_ReturnsInterfaceName",
			fmt.Sprintf("foo bar blah yahoo=hello foofoo=yes BOOTIF=%s", realMac),
		},
	}

	for _, test := range goodTests {
		t.Run(test.name, func(t *testing.T) {
			res, err := GetBootIfaceFromCmdLine(test.cmdline)
			require.NoError(t, err)
			require.Equal(t, res, realIfName)
		})
	}
}
