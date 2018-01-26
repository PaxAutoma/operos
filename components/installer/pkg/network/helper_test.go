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

package network

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIncrementIP_IncrementByOneNoRollover_ReturnsCorrectResult(t *testing.T) {
	ip := net.ParseIP("1.1.1.1")
	require.Equal(t, "1.1.1.2", IncrementIP(ip, 1).String())
}
func TestIncrementIP_IncrementByOneWithRollover_ReturnsCorrectResult(t *testing.T) {
	ip := net.ParseIP("1.1.1.255")
	require.Equal(t, "1.1.2.0", IncrementIP(ip, 1).String())
}

func TestIncrementIP_IncrementByNumberNoRollover_ReturnsCorrectResult(t *testing.T) {
	ip := net.ParseIP("1.1.1.1")
	require.Equal(t, "1.1.1.28", IncrementIP(ip, 27).String())
}

func TestIncrementIP_IncrementByNumberWithRollover_ReturnsCorrectResult(t *testing.T) {
	ip := net.ParseIP("1.1.1.1")
	require.Equal(t, "1.1.2.45", IncrementIP(ip, 300).String())
}

func TestIncrementIP_IncrementByNumberWithMultiRollover_ReturnsCorrectResult(t *testing.T) {
	ip := net.ParseIP("1.1.1.1")
	require.Equal(t, "1.1.40.17", IncrementIP(ip, 10000).String())
}
