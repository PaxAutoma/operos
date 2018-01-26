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
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStrSize(t *testing.T) {
	require.Equal(t, "1", StrSize(1))
	require.Equal(t, "100", StrSize(100))
	require.Equal(t, "1024", StrSize(1024))
	require.Equal(t, "1K", StrSize(1025))
	require.Equal(t, "1.2K", StrSize(1200))
	require.Equal(t, "120K", StrSize(122880))
	require.Equal(t, "1.2M", StrSize(1300000))
	require.Equal(t, "1.2M", StrSize(1300001))
	require.Equal(t, "775.7T", StrSize(852852852852852))
}
