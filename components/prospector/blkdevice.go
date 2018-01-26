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

package prospector

import (
	"encoding/json"
	"fmt"
	"hash/crc64"
	"os/exec"
	"strings"

	"github.com/cloudflare/cfssl/log"
)

//main function which returns the JSON representation of the device tree
//where the code is executed. Should be run wih the root priveledges, to
//ensure that all the device information is accessed correctly

type BlockDevice struct {
	MountPoint string `json:"mountpoint"`
	Name       string `json:"name"`
	KName      string `json:"kname"`
	Model      string `json:"model"`
	Serial     string `json:"serial"`
	Size       string `json:"size"`
	Rota       string `json:"rota"`
	Type       string `json:"type"`
}

type BlockDevices struct {
	BlockDevices []*BlockDevice `json:"blockdevices"`
}

func GetUUIDForHost() (*UUIDType, error) {
	xmldata, err := RunLSHW()

	if err != nil {
		log.Errorf("error: %v", err)
		return nil, fmt.Errorf("Failed to execute lshw")
	}

	v, err := NewDeviceTree(xmldata, "xml")

	if err != nil {
		log.Errorf("error: %v", err)
		return nil, fmt.Errorf("Failed to construct a device tree")
	}

	uuid, err := v.GetUUID()

	if err != nil {
		log.Errorf("error: %v", err)
		return nil, fmt.Errorf("Failed to obtain the node UUID")
	}

	return &uuid, nil
}

func UUIDStringForBlkDevice(blkdev *BlockDevice, hostUUID *UUIDType) (*string, error) {

	table := crc64.MakeTable(crc64.ECMA)
	hasher := crc64.New(table)

	hasher.Write([]byte(strings.Trim(blkdev.Model, " ")))
	hasher.Write([]byte(strings.Trim(blkdev.Serial, " ")))
	hasher.Write([]byte(strings.Trim(blkdev.Size, " ")))

	deviceHash := hasher.Sum([]byte{})
	bytesPerDiskCombinedHash := 16
	uuidBytes := make([]byte, bytesPerDiskCombinedHash)

	j := 0
	for i := 0; i < bytesPerDiskCombinedHash-len(deviceHash); i++ {
		uuidBytes[i] = (*hostUUID)[j]
		j++
	}

	for _, v := range deviceHash {
		uuidBytes[j] = v
		j++
	}

	res := UUIDFromBytes(uuidBytes)

	return &res, nil

}

//GenerateUUIDForBlockDevices returns a maps of device name:UUID
func GenerateUUIDForBlockDevices(lsblkData []byte, hostUUID *UUIDType) (*map[string]string, error) {

	blockDevices := BlockDevices{}

	res := make(map[string]string)

	err := json.Unmarshal(lsblkData, &blockDevices)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse lsblk output due to %s", err)
	}



	for _, blkDevice := range blockDevices.BlockDevices {
		//check only disks
		if blkDevice.Type == "disk" {

			uuid, err := UUIDStringForBlkDevice(blkDevice, hostUUID)
			if err != nil {
				return nil, fmt.Errorf("Failed to compute UUID for %s becase of %s", blkDevice.Name, err)
			}
			res[blkDevice.Name] = *uuid

		}
	}
	return &res, nil

}

//RunLSBLK executes the lsblk and returns its output in JSON
func RunLSBLK() ([]byte, error) {
	out, err := exec.Command("lsblk",
		"-o",
		"MOUNTPOINT,NAME,KNAME,MODEL,SERIAL,SIZE,ROTA,TYPE",
		"-b",
		"--json").Output()
	if err != nil {
		return nil, fmt.Errorf("Could not run lsblk: %s", err)
	}
	return out, nil
}
