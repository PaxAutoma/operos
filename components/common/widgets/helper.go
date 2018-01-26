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
	"strings"

	"github.com/jroimartin/gocui"
)

func CenterString(str string, width int) string {
	lines := strings.Split(str, "\n")

	result := make([]string, len(lines))
	for idx, line := range lines {
		if width <= len(line) {
			result[idx] = line
			continue
		}
		paddingLen := width - len(line)

		prefix := strings.Repeat(" ", paddingLen/2)
		postfix := strings.Repeat(" ", paddingLen/2+paddingLen%2)
		result[idx] += prefix + line + postfix
	}

	return strings.Join(result, "\n")
}

func CenterInBox(str string, width, height int) string {
	centeredText := CenterString(str, width)
	numLines := strings.Count(str, "\n") + 1

	if height <= numLines {
		return centeredText
	}

	numPadding := height - numLines
	padLine := strings.Repeat(" ", width) + "\n"
	paddingPre := strings.Repeat(padLine, numPadding/2)
	paddingPost := strings.Repeat(padLine, numPadding/2+numPadding%2)

	return paddingPre + centeredText + paddingPost
}

func ColorString(color gocui.Attribute, str string) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", 30+color-1, str)
}

func ReverseString(str string) string {
	return fmt.Sprintf("\033[7m%s\033[0m", str)
}

func BoldString(color gocui.Attribute, str string) string {
	return fmt.Sprintf("\033[%d;1m%s\033[0m", 30+color-1, str)
}
