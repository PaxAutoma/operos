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
	"bufio"
	"fmt"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type DiskInfo struct {
	Name  string
	Size  int
	Model string
}

func (d DiskInfo) StrSize() string {
	return StrSize(float64(d.Size))
}

func StrSize(size float64) string {
	units := []string{"K", "M", "G", "T", "P", "E", "Z", "Y"}
	unit := ""

	var idx int
	for idx = range units {
		if size <= math.Pow(1024, float64(idx+1)) {
			break
		}
		unit = units[idx]
	}

	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.1f", size/math.Pow(1024, float64(idx))), "0"), ".") + unit
}

func (d DiskInfo) String() string {
	return fmt.Sprintf("%s (%s): %s", d.Name, d.Model, d.StrSize())
}

func GetDiskList() ([]DiskInfo, error) {
	cmd := exec.Command("/bin/lsblk", "-r", "-o", "TYPE,NAME,TRAN,RO,MOUNTPOINT,SIZE,MODEL", "-b")
	lsblkOutput, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.Wrap(err, "could not run lsblk")
	}

	result := make([]DiskInfo, 0)
	lines := strings.Split(string(lsblkOutput), "\n")
	for _, line := range lines {
		parts := strings.Split(line, " ")
		// Want only read/write disks that are not currently mounted (i.e. are not the disk we're booting from)
		if parts[0] != "disk" || parts[2] == "usb" || parts[3] != "0" || parts[4] != "" {
			continue
		}

		size, _ := strconv.Atoi(parts[5])
		result = append(result, DiskInfo{
			Name:  parts[1],
			Model: strings.TrimSpace(strings.Replace(parts[6], "\\x20", " ", -1)),
			Size:  size,
		})
	}

	return result, nil
}

func GetNumCPUs() int {
	return runtime.NumCPU()
}

func GetTotalMemory() (int, error) {
	memInfo, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, errors.Wrap(err, "failed to read /proc/meminfo")
	}

	scanner := bufio.NewScanner(memInfo)
	for scanner.Scan() {
		if scanner.Err() != nil {
			return 0, errors.Wrap(err, "failed to read /proc/meminfo")
		}

		fields := strings.Fields(scanner.Text())
		if fields[0] == "MemTotal:" {
			kibibytesIsh, err := strconv.Atoi(fields[1])
			if err != nil {
				return 0, errors.Wrap(err, "invalid format of /proc/meminfo")
			}

			return int(float64(kibibytesIsh)/(1024*1024) + 0.5), nil
		}
	}

	return 0, errors.Errorf("failed to find MemTotal in /proc/meminfo")
}
