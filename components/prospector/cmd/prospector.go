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

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/paxautoma/operos/components/prospector"
)

//GetHostUUID function which returns the JSON representation of the device tree
//where the code is executed. Should be run wih the root priveledges, to
//ensure that all the device information is accessed correctly
func GetHostUUID(lshwXMLFile *string) {

	var uuid prospector.UUIDType

	if *lshwXMLFile != "" {

		xmldata, err := ioutil.ReadFile(*lshwXMLFile)

		if err != nil {
			fmt.Printf("error: %v", err)
			return
		}

		v, err := prospector.NewDeviceTree(xmldata, "xml")
		uuid, err = v.GetUUID()

	} else {
		uuidRef, err := prospector.GetUUIDForHost()
		uuid = *uuidRef
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to obatin generate UUIDs for blockdevices due to %s\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("UUID:%s\n", uuid.ToString())

}

//GetBlockDevices returns the JSON representation of the device tree
//where the code is executed. Should be run wih the root priveledges, to
//ensure that all the device information is accessed correctly
func GetBlockDevices() {

	devicedata, err := prospector.RunLSBLK()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to obatin generate UUIDs for blockdevices due to %s\n", err)
		os.Exit(1)
	}

	hostUUID, err := prospector.GetUUIDForHost()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to calculate UUID for host due to %s", err)
		os.Exit(1)
	}

	blockDevices, err := prospector.GenerateUUIDForBlockDevices(devicedata, hostUUID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to obatin generate UUIDs for blockdevices due to %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Device,UUID\n")

	keys := []string{}

	for deviceName := range *blockDevices {
		keys = append(keys, deviceName)
	}

	sort.Strings(keys)

	for _, deviceName := range keys {
		fmt.Printf("%s,%s\n", deviceName, (*blockDevices)[deviceName])
	}

}

//ShowDeviceTree returns the JSON representation of the device tree
//where the code is executed. Should be run wih the root priveledges, to
//ensure that all the device information is accessed correctly
func ShowDeviceTree() {
	xmldata, err := prospector.RunLSHW()

	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	v, err := prospector.NewDeviceTree(xmldata, "xml")

	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	devicedata, err := prospector.RunLSBLK()
	var blockDevices prospector.BlockDevices

	err = json.Unmarshal(devicedata, &blockDevices)
	if err != nil {
		return
	}

	out := new(prospector.Report)
	out.System = v
	out.Storage = blockDevices

	if jsonout, err := json.Marshal(out); err == nil {
		fmt.Println(string(jsonout))
	} else {
		return
	}
}

var getBlockDevices = flag.Bool("blk-device-uuid", false, "Generate the UUID for block devices")
var hostUUIDOnly = flag.String("host-uuid-only", "", "The name of the XML file which to generate UUID")

func main() {
	flag.Parse()

	if *getBlockDevices {
		GetBlockDevices()
	} else if *hostUUIDOnly != "" {
		GetHostUUID(hostUUIDOnly)
	} else {
		ShowDeviceTree()
	}

}
