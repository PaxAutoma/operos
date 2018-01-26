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
	"encoding/binary"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"hash"
	"hash/crc32"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

//BlackListedDeviceClasses list the classes of devices which should be ignored
//const BlackListedDeviceClasses = []string{"volume"}
//TODO:add support for the blacklist of the devices

//CapabilityType describes the capabilities of a node
type CapabilityType struct {
	ID          string `xml:"id,attr" json:"id"`
	Description string `xml:"description" json:"description"`
}

//CapabilitiesList is wrapper for the list of capabilities
type CapabilitiesList struct {
	Capability []*CapabilityType `xml:"capability" json:"capability"`
}

//Setting represent a list of settings in a Node
type Setting struct {
	ID    string `xml:"id,attr" json:"id"`
	Value string `xm:"value,attr" json:"value"`
}

//SettingsList is wrapper for the list of settings
type SettingsList struct {
	Settings []*Setting `xml:"setting" json:"setting"`
}

//Int64WithUnit number with a unit
type Int64WithUnit struct {
	Units string `xml:"units,attr" json:"units"`
	Value string `xml:",chardata" json:"value"`
}

//Device is a node in a device tree
type Device struct {
	ID          string `xml:"id,attr" json:"id"`
	Class       string `xml:"class,attr" json:"class"`
	Handle      string `xml:"handle,attr" json:"handle"`
	Product     string `xml:"product" json:"product"`
	Description string `xml:"description" json:"description"`
	Vendor      string `xml:"vendor" json:"vendor"`
	Version     string `xml:"version" json:"version"`
	Serial      string `xml:"serial" json:"serial"`
	Slot        string `xml:"slot" json:"slot"`

	//TODO: take the UUID out as it won't be needed

	hash   uint32 //crc32 used for calculating the UUID
	weight int    `` //the weight of UUID in the tree

	Size     Int64WithUnit `xml:"size" json:"size,omitempty"`
	Capacity Int64WithUnit `xml:"capacity" json:"capacity,omitempty"`

	Capabilities  CapabilitiesList `xml:"capabilities" json:"capabilities"`
	Configuration SettingsList     `xml:"configuration" json:"configuration"`

	Devices []*Device `xml:"node" json:"Nodes,omitempty"`
}

//DeviceTree is tree representing a tree of thedevices available
type DeviceTree struct {

	//UUID   UUIDType `json:"uuid"`
	System *Device `xml:"node" json:"system"`
}

//GetUUID comptes a unqique deterministic 64bit UUID based of the uuid of the system.
//The latter is calculated recursively, by calculating the crc32  of the relevant fields. The crc code is scaled down to fit the range of [0..1).
//Each device has a unique weight which determines its "imporantce" in the hash. The scaled down hashes weighted and added together to calculate
//the hash of the system. The latter is stored as float64 and converted to uint64, according to its IEEE 754 bit representation.
//with some weight.
func (deviceTree *DeviceTree) GetUUID() (UUIDType, error) {
	//reset the UUID in case it was called earlier
	resUUID := GetZeroUUID()
	deviceTree.assignWeights()

	alternativeSerial := ""
	//In VMS there are no serial numbers
	if deviceTree.allSignificantDevicesLackSerials() {
		tmp, err := deviceTree.System.getMacAddress()
		if err != nil {
			return UUIDType{}, errors.New("Device tree lacks serials and MAC addresses")
		}
		alternativeSerial = tmp

	}

	if err := deviceTree.System.assignHash(alternativeSerial); err != nil {
		log.Printf("Failed to generate the UUID for a tree due to:%v", err)
		return UUIDType{}, errors.New("Failed to generate device tree")
	}

	deviceTree.System.updateUUID(&resUUID)

	return resUUID, nil
}

//Traverses the tree and finds if at least one significant device has a serial
func (deviceTree *DeviceTree) allSignificantDevicesLackSerials() bool {
	return deviceTree.System.allSignificantDevicesLackSerials()
}

func (device *Device) isSignificant() bool {
	if device.weight < BytesForMajorComponents {
		return true
	}
	return false

}

//Traverses the tree and finds if at least one significant device has a serial

func (device *Device) hasValidSerial() bool {
	return device.Serial != "" && device.Serial != "0"
}

func (device *Device) allSignificantDevicesLackSerials() bool {
	if device.isSignificant() == true {
		if device.hasValidSerial() {
			return false
		}
		for _, child := range device.Devices {
			if !child.isSignificant() {
				break
			}

			if child.allSignificantDevicesLackSerials() != true {
				return false
			}
		}
		return true

	}
	return true

}

//getMacAddress returns the the first MAC address of the nice NICs
//if there are not devices, returns an error
func (device *Device) getMacAddress() (string, error) {
	if device.Class == "network" {
		return device.Serial, nil
	}

	for _, child := range device.Devices {
		serial, err := child.getMacAddress()
		if err == nil {
			return serial, nil
		}
	}

	return "", fmt.Errorf("The network subdevices found")

}

// //GetFloat64UUID returns the float representtion of the UUID
// func (deviceTree *DeviceTree) GetFloat64UUID() float64 {
// 	return float64()
// }

//updateUUID traverses the device tree. It add the hash of each device to the relevant parts in the passed UUID
func (device *Device) updateUUID(uuid *UUIDType) {

	localUUID := make([]byte, BytePerTreeLevel)
	binary.LittleEndian.PutUint32(localUUID, device.hash)
	localUUIDIndex := 0

	//add localUUID to the relevant bits of the global UUID
	for i := device.weight; i < device.weight+BytePerTreeLevel && i < len(uuid); i++ {
		uuid[i] = uuid[i] + localUUID[localUUIDIndex]
		localUUIDIndex++
	}

	for _, child := range device.Devices {
		child.updateUUID(uuid)
	}

}

func GetHasher() hash.Hash32 {
	return crc32.NewIEEE()
}

func (device *Device) assignHash(alternativeSerial string) error {

	hasher := GetHasher()

	if device.ID == "" &&
		device.Serial == "" &&
		device.Class == "" &&
		device.Product == "" &&
		device.Vendor == "" {
		//TODO: Do this only for the major parts
		return fmt.Errorf("ID, Serial, Class, Product and Vendor are empty. Cannot generate hash for device %s", device)
	}

	if device.Class != "system" {
		//ignore the hostname which is an id
		//for system
		hasher.Write([]byte(device.ID))
	}

	if device.isSignificant() && !device.hasValidSerial() {
		hasher.Write([]byte(alternativeSerial))
	} else {
		hasher.Write([]byte(device.Serial))
	}
	hasher.Write([]byte(device.Class))
	hasher.Write([]byte(device.Product))
	//the vendor name which starts with Linux need to be ignored
	//as it would change with the version of linux
	if !strings.HasPrefix(strings.ToLower(device.Vendor), "linux") {
		hasher.Write([]byte(device.Vendor))
	}

	device.hash = hasher.Sum32()

	for _, childDevice := range device.Devices {

		err := childDevice.assignHash(alternativeSerial)
		if err != nil {
			return err
		}
	}

	return nil

}

//assignWeights assigns weight according to the rules
//TODO: switch from naive implementation for more case specific
//For example, one can provide a map like /motherboard/memory, 0.2
//to manually override the weight assignment
func (device *Device) assignWeights(currentWegith int) {

	device.weight = currentWegith
	childWeight := currentWegith + BytePerTreeLevel

	//
	for _, childDevice := range device.Devices {
		childDevice.assignWeights(childWeight) //to make sure that the buckets have at least)
	}
	return

}

//assignWeights assigns weight according to the rules
func (deviceTree *DeviceTree) assignWeights() {
	deviceTree.System.assignWeights(0)
}

//LoadDeviceTree constructs a device tree from an XML file
func LoadDeviceTree(xmlFilename string) (*DeviceTree, error) {
	xmldata, err := ioutil.ReadFile(xmlFilename)

	if err != nil {
		log.Printf("Could not read in the input file %s due to %s", xmlFilename, err)
		return nil, errors.New("Failed to load a device tree due to IO error")
	}

	devTree, err := NewDeviceTree(xmldata, "xml")
	if err != nil {
		log.Printf("Could not parse XML from %s due to %s", xmlFilename, err)
		return nil, errors.New("Failed to load a device tree due to parsing")
	}

	return devTree, nil

}

//NewDeviceTree contructs a device tree from the XML or JSON produced by ran with lshw --xml option.
//format can be "xml" or "json"
func NewDeviceTree(data []byte, format string) (*DeviceTree, error) {

	deviceTree := DeviceTree{}
	var err error

	if format == "xml" {
		err = xml.Unmarshal(data, &deviceTree)
	} else {
		err = json.Unmarshal(data, &deviceTree)
	}
	if err != nil {
		err = fmt.Errorf("Failed to construct a valid device tree due to error: %v", err)
		return nil, err
	}

	if err != nil {
		return nil, errors.New("Failed to assign UUIDs")
	}

	return &deviceTree, nil

}

func (deviceTree *DeviceTree) String() string {
	return (deviceTree.System.String())

}

func (device *Device) string(ident string) string {
	res := ""
	for _, child := range device.Devices {
		res = res + child.string(ident+"  ")
	}
	return (fmt.Sprintf("%sdevice:%s\t class:%s\t crc32:%d\t  weight:%d\n",
		ident,
		device.Description,
		device.Class,
		device.hash,
		device.weight) + res)
}

func (device *Device) String() string {
	return device.string("")
}

//RunLSHW executes the lshw and returns its output in XML
func RunLSHW() ([]byte, error) {
	out, err := exec.Command("lshw", "-xml").Output()
	if err != nil {
		log.Printf("Could not run lshw: %s", err)
		return nil, errors.New("failed to exectue lshw")
	}
	return out, nil
}

//ToJSONreturns a JSON representation of the device tree
//returns []byte,error. If error is nil, everything went fine
//if not...well there is an error
func (deviceTree *DeviceTree) ToJSON() ([]byte, error) {
	json, err := json.Marshal(deviceTree)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}
	return json, nil

}
