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
	"io/ioutil"
	"math"
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

const TSTPath string = "tst"

//assertNearEqualFloat64 checks if the obeserved is within the epsilon of the expected
//if not, then fails the test
func assertNearEqualFloat64(expected float64, observed float64, name string, epsilon float64, t *testing.T) {
	if math.Abs(expected-observed) > epsilon {
		t.Errorf("%s was expected to be %f and  it is %f,whic is further than %e from it.", name, expected, observed, epsilon)
	}
}

//randomString generates a random string of a fixed lenght with a probability of an empty string
func randomString(length int, probabilityOfEmpty float64) string {
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	if rand.Float64() > probabilityOfEmpty {
		res := ""
		for i := 0; i < length; i++ {
			res += string(charset[rand.Intn(len(charset))])
		}
		return res

	}

	return ""

}

func TestNewDeviceTreeWithDifferentInputs(t *testing.T) {
	type args struct {
		xmldata []byte
	}
	tests := []struct {
		name        string
		testXMLFile string
		want        *DeviceTree
		wantErr     bool
	}{
		{name: "ThinkPad T430s", testXMLFile: TSTPath + "/compatability/test_01.xml", wantErr: false},
		{name: "Homebrew MSI server", testXMLFile: TSTPath + "/compatability/test_02.xml", wantErr: false},
		{name: "Desktop MS-7673", testXMLFile: TSTPath + "/compatability/test_03.xml", wantErr: false},
		{name: "Asus Desktop", testXMLFile: TSTPath + "/compatability/test_04.xml", wantErr: false},
		//----and the large set

		{name: "0.0.0.0.out", testXMLFile: TSTPath + "/compatability_large_set/0.0.0.0.out", wantErr: false},
		{name: "10.75.10.100.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.100.out", wantErr: false},
		{name: "10.75.10.101.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.101.out", wantErr: false},
		{name: "10.75.10.102.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.102.out", wantErr: false},
		{name: "10.75.10.103.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.103.out", wantErr: false},
		{name: "10.75.10.105.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.105.out", wantErr: false},
		{name: "10.75.10.107.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.107.out", wantErr: false},
		{name: "10.75.10.108.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.108.out", wantErr: false},
		{name: "10.75.10.109.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.109.out", wantErr: false},
		{name: "10.75.10.110.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.110.out", wantErr: false},
		{name: "10.75.10.111.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.111.out", wantErr: false},
		{name: "10.75.10.113.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.113.out", wantErr: false},
		{name: "10.75.10.114.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.114.out", wantErr: false},
		{name: "10.75.10.115.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.115.out", wantErr: false},
		{name: "10.75.10.116.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.116.out", wantErr: false},
		{name: "10.75.10.117.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.117.out", wantErr: false},
		{name: "10.75.10.119.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.119.out", wantErr: false},
		{name: "10.75.10.121.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.121.out", wantErr: false},
		{name: "10.75.10.122.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.122.out", wantErr: false},
		{name: "10.75.10.123.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.123.out", wantErr: false},
		{name: "10.75.10.124.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.124.out", wantErr: false},
		{name: "10.75.10.125.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.125.out", wantErr: false},
		{name: "10.75.10.127.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.127.out", wantErr: false},
		{name: "10.75.10.44.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.44.out", wantErr: false},
		{name: "10.75.10.45.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.45.out", wantErr: false},
		{name: "10.75.10.46.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.46.out", wantErr: false},
		{name: "10.75.10.47.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.47.out", wantErr: false},
		{name: "10.75.10.48.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.48.out", wantErr: false},
		{name: "10.75.10.49.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.49.out", wantErr: false},
		{name: "10.75.10.50.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.50.out", wantErr: false},
		{name: "10.75.10.51.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.51.out", wantErr: false},
		{name: "10.75.10.52.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.52.out", wantErr: false},
		{name: "10.75.10.53.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.53.out", wantErr: false},
		{name: "10.75.10.54.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.54.out", wantErr: false},
		{name: "10.75.10.55.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.55.out", wantErr: false},
		{name: "10.75.10.56.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.56.out", wantErr: false},
		{name: "10.75.10.57.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.57.out", wantErr: false},
		{name: "10.75.10.58.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.58.out", wantErr: false},
		{name: "10.75.10.59.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.59.out", wantErr: false},
		{name: "10.75.10.60.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.60.out", wantErr: false},
		{name: "10.75.10.61.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.61.out", wantErr: false},
		{name: "10.75.10.62.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.62.out", wantErr: false},
		{name: "10.75.10.63.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.63.out", wantErr: false},
		{name: "10.75.10.64.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.64.out", wantErr: false},
		{name: "10.75.10.65.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.65.out", wantErr: false},
		{name: "10.75.10.67.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.67.out", wantErr: false},
		{name: "10.75.10.69.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.69.out", wantErr: false},
		{name: "10.75.10.70.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.70.out", wantErr: false},
		{name: "10.75.10.71.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.71.out", wantErr: false},
		{name: "10.75.10.72.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.72.out", wantErr: false},
		{name: "10.75.10.73.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.73.out", wantErr: false},
		{name: "10.75.10.74.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.74.out", wantErr: false},
		{name: "10.75.10.76.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.76.out", wantErr: false},
		{name: "10.75.10.77.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.77.out", wantErr: false},
		{name: "10.75.10.78.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.78.out", wantErr: false},
		{name: "10.75.10.79.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.79.out", wantErr: false},
		{name: "10.75.10.80.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.80.out", wantErr: false},
		{name: "10.75.10.81.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.81.out", wantErr: false},
		{name: "10.75.10.82.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.82.out", wantErr: false},
		{name: "10.75.10.83.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.83.out", wantErr: false},
		{name: "10.75.10.84.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.84.out", wantErr: false},
		{name: "10.75.10.85.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.85.out", wantErr: false},
		{name: "10.75.10.86.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.86.out", wantErr: false},
		{name: "10.75.10.87.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.87.out", wantErr: false},
		{name: "10.75.10.88.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.88.out", wantErr: false},
		{name: "10.75.10.89.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.89.out", wantErr: false},
		{name: "10.75.10.90.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.90.out", wantErr: false},
		{name: "10.75.10.91.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.91.out", wantErr: false},
		{name: "10.75.10.92.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.92.out", wantErr: false},
		{name: "10.75.10.94.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.94.out", wantErr: false},
		{name: "10.75.10.95.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.95.out", wantErr: false},
		{name: "10.75.10.96.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.96.out", wantErr: false},
		{name: "10.75.10.97.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.97.out", wantErr: false},
		{name: "10.75.10.98.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.98.out", wantErr: false},
		{name: "10.75.10.99.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.99.out", wantErr: false},
		{name: "10.75.11.2.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.11.2.out", wantErr: false},
		{name: "10.75.11.3.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.11.3.out", wantErr: false},
		{name: "10.75.11.4.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.11.4.out", wantErr: false},
		{name: "10.75.11.6.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.11.6.out", wantErr: false},
		{name: "10.75.11.7.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.11.7.out", wantErr: false},
		{name: "10.75.12.1.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.12.1.out", wantErr: false},
		{name: "10.75.12.2.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.12.2.out", wantErr: false},
		{name: "10.75.12.3.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.12.3.out", wantErr: false},
		{name: "127.0.0.1.out", testXMLFile: TSTPath + "/compatability_large_set/127.0.0.1.out", wantErr: false},
		{name: "10.75.12.13.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.12.13.out", wantErr: false},
		{name: "::1.out", testXMLFile: TSTPath + "/compatability_large_set/::1.out", wantErr: false},
	}
	for _, tt := range tests {

		devTreeDeviceTree, err := LoadDeviceTree(tt.testXMLFile)

		if (err == nil) == tt.wantErr {
			t.Errorf("Did not expect an error for case %s", tt.name)
		}

		fmt.Println(devTreeDeviceTree.GetUUID())
	}
}

// //generateRandomDevice

func generateRandomDevice(maxDeviceTreeDepth int, maxChildren int, probabiltyOfChild float64, fieldLength int, probEmptyField float64) *Device {

	emptyDevice := true
	res := Device{}

	for emptyDevice == true {
		res = Device{
			ID:          randomString(fieldLength, probEmptyField),
			Class:       randomString(fieldLength, probEmptyField),
			Handle:      randomString(fieldLength, probEmptyField),
			Product:     randomString(fieldLength, probEmptyField),
			Description: randomString(fieldLength, probEmptyField),
			Vendor:      randomString(fieldLength, probEmptyField),
			Version:     randomString(fieldLength, probEmptyField),
			Serial:      randomString(fieldLength, probEmptyField),
			Slot:        randomString(fieldLength, probEmptyField),
		}

		emptyDevice = false

	}

	if maxDeviceTreeDepth > 1 { //generte the children by rolling the dice for each child
		for i := 0; i < maxChildren; i++ {
			if rand.Float64() < probabiltyOfChild {
				child := generateRandomDevice(maxDeviceTreeDepth-1, maxChildren, probabiltyOfChild, fieldLength, probEmptyField)
				res.Devices = append(res.Devices, child)
			}
		}
	}

	res.assignWeights(0)
	res.assignHash("")

	return &res
}

//TestNewDeviceTreeForCollisionProbability tests UUID would colide for 1,0000 machines of the complete random specs
//and parts
func TestNewDeviceTreeForCollisionProbability(t *testing.T) {
	//generate 25,000 random trees, srings random[7], 20% empty, 50% of having a child

	const numServers = 25000
	maxChildren := 4
	maxDeviceTreeDepth := 4
	probabiltyOfChild := 0.2

	fieldLength := 50     //the length of the device struc fields
	probEmptyField := 0.1 //probability that one of the fields in the device structure won't be filled
	collisionTolerance := 1e-10

	rand.Seed(20)

	var servers [numServers]*DeviceTree

	uUIDList := UUIDs{}

	for serverIndex := 0; serverIndex < numServers; serverIndex++ {
		system := generateRandomDevice(maxChildren, maxDeviceTreeDepth, probabiltyOfChild, fieldLength, probEmptyField)
		devTree := DeviceTree{System: system}

		servers[serverIndex] = &devTree

		devTree.assignWeights()
		UUID, err := devTree.GetUUID()

		if err != nil {
			t.Log(err)
			t.Logf("Failed to assign UUID")
		} else {
			uUIDList = append(uUIDList, UUID)
		}

	}

	findCollisions(uUIDList, collisionTolerance, t)

}

func findCollisions(uUIDList UUIDs, collisionTolerance float64, t *testing.T) {
	sort.Sort(uUIDList)
	numServers := len(uUIDList)

	collions := int(0)
	majorPartsCollisions := 0
	for i := 0; i < len(uUIDList); i++ {
		for j := 0; j < i; j++ {
			if uUIDList[i].IsIdenticalTo(&uUIDList[j]) {
				collions++
				fmt.Printf("A collision between uuid:%v and uuid:%v\n", uUIDList[i], uUIDList[j])
			}

			if uUIDList[i].HasTheSameMajorParts(&uUIDList[j]) {
				majorPartsCollisions++
				fmt.Printf("A major parts collision between uuid:%v and uuid:%v\n", uUIDList[i], uUIDList[j])
			}
		}
	}

	fmt.Printf("Observed %d collisions for %d random servers\n", collions, numServers)
	fmt.Printf("Observed %d major parts collisions for %d random servers\n", majorPartsCollisions, numServers)

	probOfCollision := float64(collions) / float64(numServers)
	probOfMajorPartsCollision := float64(majorPartsCollisions) / float64(numServers)

	assertNearEqualFloat64(0.0, probOfCollision, "Proability of colliison", collisionTolerance, t)
	assertNearEqualFloat64(0.0, probOfMajorPartsCollision, "Proability of colliison", collisionTolerance, t)
}

func (dev *Device) randomizeSerials(length int) {
	dev.Serial = randomString(length, 0.0)

	for _, childDev := range dev.Devices {
		childDev.randomizeSerials(length)
	}

	return
}

//TestNewDeviceTreeForCollisionProbabilitySerialsOnly tests UUID would colide for 1,0000 machines of the same spec
//with parts having different serial numbers
func TestNewDeviceTreeForCollisionProbabilitySerialsOnly(t *testing.T) {

	testXMLFile := TSTPath + "/collisionprobability/TestNewDeviceTreeForCollisionProbabilitySerialsOnly.xml"
	const numServers = 1000
	serialNumLen := 10
	collisionTolerance := 1e-10

	xmldata, err := ioutil.ReadFile(testXMLFile)

	if err != nil {
		t.Errorf("Could not read in the input file %s", testXMLFile)
	}

	rand.Seed(20)

	//var servers [numServers]*DeviceTree

	uUIDList := UUIDs{}

	for serverIndex := 0; serverIndex < numServers; serverIndex++ {
		//randomize the serial numbers
		devTree, err := NewDeviceTree(xmldata, "xml")
		if (err != nil) != false {
			t.Errorf("Failed to create a valid device tree from %s", testXMLFile)
			return
		}

		devTree.System.randomizeSerials(serialNumLen)

		devTree.assignWeights()
		UUID, err := devTree.GetUUID()

		if err != nil {
			t.Errorf("Failed to assign UUID")
		}
		//servers[serverIndex] = devTree

		uUIDList = append(uUIDList, UUID)
	}

	findCollisions(uUIDList, collisionTolerance, t)

}

func TestSensativityToPartsChanged(t *testing.T) {

	serialNumLen := 10

	type args struct {
		xmlFilename string
	}
	tests := []struct {
		name        string
		testXMLFile string
		wantErr     bool
	}{
		{name: "ThinkPad T430s", testXMLFile: TSTPath + "/compatability/test_01.xml", wantErr: false},
		{name: "Homebrew MSI server", testXMLFile: TSTPath + "/compatability/test_02.xml", wantErr: false},
		{name: "Desktop MS-7673", testXMLFile: TSTPath + "/compatability/test_03.xml", wantErr: false},
		{name: "Asus Desktop", testXMLFile: TSTPath + "/compatability/test_04.xml", wantErr: false},
		//----and the large set

		// {name: "0.0.0.0.out", testXMLFile: TSTPath + "/compatability_large_set/0.0.0.0.out", wantErr: false},
		{name: "10.75.10.100.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.100.out", wantErr: false},
		{name: "10.75.10.101.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.101.out", wantErr: false},
		{name: "10.75.10.102.out", testXMLFile: TSTPath + "/compatability_large_set/10.75.10.102.out", wantErr: false},
	}
	for _, tt := range tests {
		devTree, err := LoadDeviceTree(tt.testXMLFile)

		UUID, err := devTree.GetUUID()

		if err != nil {
			t.Errorf("Failed to assign UUID")
		}

		if err != nil {
			t.Errorf("%q. LoadDeviceTree() error = %v, wantErr %v", tt.name, err, tt.testXMLFile)
			continue
		}

		for _, depth := range []int{0, 1, 2, 3, 4} {
			devTree2, _ := LoadDeviceTree(tt.testXMLFile)

			device := devTree2.System
			i := 0
			for ; i < depth; i++ {
				if len(device.Devices) > 0 {
					device = device.Devices[rand.Intn(len(device.Devices))]
				} else {
					break
				}

			}

			device.Serial = randomString(serialNumLen, 0.0)
			UUID2, err := devTree2.GetUUID()

			if err != nil {
				t.Errorf("Failed to assign UUID")
			}

			expectedUUIDDifference := BytesPerUUID - i*BytePerTreeLevel
			if UUID.BytesDiffer(&UUID2) != expectedUUIDDifference {
				t.Errorf("Test %s Depth %d oldUUID %v newUUID %v Expected Difference %d Oberved difference %v \n",
					tt.name,
					i,
					UUID,
					UUID,
					expectedUUIDDifference,
					UUID.BytesDiffer(&UUID2))

			}

		}

	}
}

func Test_GetUUID_VMS(t *testing.T) {
	node, err := LoadDeviceTree("tst/vbox/controller.xml")
	if err != nil {
		t.Errorf("Could not read the file, %s", err)
	}

	nodeNewNic, err := LoadDeviceTree("tst/vbox/controller_different_nic.xml")
	if err != nil {
		t.Errorf("Could not read the file, %s", err)
	}

	nodeUUID, err := node.GetUUID()
	fmt.Println(node.System.string(" "))
	if err != nil {
		t.Errorf("Failed to compute UUID, %s", err)
	}

	nodeNewNicUUID, err := nodeNewNic.GetUUID()
	if err != nil {
		t.Errorf("Failed to compute UUID, %s", err)
	}
	fmt.Println(nodeNewNic.System.string(" "))

	if nodeUUID.ToString() == nodeNewNicUUID.ToString() {
		t.Errorf("UUIDs for VMs with different nics should be different")
	}

}

func Test_GetUUIDVMSWorkers(t *testing.T) {
	node, err := LoadDeviceTree("tst/vbox/node1.xml")
	if err != nil {
		t.Errorf("Could not read the file, %s", err)
	}

	nodeNewNic, err := LoadDeviceTree("tst/vbox/node2.xml")
	if err != nil {
		t.Errorf("Could not read the file, %s", err)
	}

	nodeUUID, err := node.GetUUID()
	if err != nil {
		t.Errorf("Failed to compute UUID, %s", err)
	}

	nodeNewNicUUID, err := nodeNewNic.GetUUID()
	if err != nil {
		t.Errorf("Failed to compute UUID, %s", err)
	}

	if nodeUUID.ToString() == nodeNewNicUUID.ToString() {
		t.Errorf("UUIDs for VMs with different nics should be different")
	}
}

func Test_GetUUIDVMSRealHW(t *testing.T) {
	node, err := LoadDeviceTree("tst/vbox/real_hw.xml")
	if err != nil {
		t.Errorf("Could not read the file, %s", err)
	}

	nodeNewNic, err := LoadDeviceTree("tst/vbox/real_hw_new_nic.xml")
	if err != nil {
		t.Errorf("Could not read the file, %s", err)
	}

	nodeUUID, err := node.GetUUID()
	if err != nil {
		t.Errorf("Failed to compute UUID, %s", err)
	}

	nodeNewNicUUID, err := nodeNewNic.GetUUID()
	if err != nil {
		t.Errorf("Failed to compute UUID, %s", err)
	}

	if nodeUUID.ToString() != nodeNewNicUUID.ToString() {
		t.Errorf("UUIDs for physical HW with different nics should be the same")
	}

}

func Test_runlshw(t *testing.T) {
	tests := []struct {
		name string

		wantErr bool
	}{
		{name: "basic test",
			wantErr: false},
	}
	for _, tt := range tests {
		devTreeXML, err := RunLSHW()
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. runlshw() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}

		devTree, err := NewDeviceTree(devTreeXML, "xml")

		fmt.Print(devTree)

		if err != nil {
			t.Errorf("Failed to construct a valid device tree from the output of lshw: %s", devTreeXML)
		}

	}
}

func TestDevice_assignWeights(t *testing.T) {
	type fields struct {
		ID            string
		Class         string
		Handle        string
		Product       string
		Description   string
		Vendor        string
		Version       string
		Serial        string
		Slot          string
		uuid          UUIDType
		hash          uint32
		weight        int
		Size          Int64WithUnit
		Capacity      Int64WithUnit
		Capabilities  CapabilitiesList
		Configuration SettingsList
		Devices       []*Device
	}
	type args struct {
		currentWegith int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		device := &Device{
			ID:            tt.fields.ID,
			Class:         tt.fields.Class,
			Handle:        tt.fields.Handle,
			Product:       tt.fields.Product,
			Description:   tt.fields.Description,
			Vendor:        tt.fields.Vendor,
			Version:       tt.fields.Version,
			Serial:        tt.fields.Serial,
			Slot:          tt.fields.Slot,
			hash:          tt.fields.hash,
			weight:        tt.fields.weight,
			Size:          tt.fields.Size,
			Capacity:      tt.fields.Capacity,
			Capabilities:  tt.fields.Capabilities,
			Configuration: tt.fields.Configuration,
			Devices:       tt.fields.Devices,
		}
		device.assignWeights(tt.args.currentWegith)
	}
}

func TestDeviceTree_assignWeights(t *testing.T) {

	DevABC := Device{
		weight: -100,
	}

	DevABE := Device{
		weight: -100,
	}

	DevAB := Device{
		weight:  -100,
		Devices: []*Device{&DevABC, &DevABE},
	}

	DevA := Device{
		weight:  -100,
		Devices: []*Device{&DevAB},
	}

	type fields struct {
		UUID   UUIDType
		System *Device
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{name: "single node test", fields: fields{UUID: UUIDType{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			System: &DevA}},
	}
	for _, tt := range tests {
		deviceTree := &DeviceTree{
			System: tt.fields.System,
		}
		deviceTree.assignWeights()

		if DevABC.weight != 2*BytePerTreeLevel {
			t.Errorf("Expected weight %d and got %d instead", 2*BytePerTreeLevel, DevABC.weight)
		}

		if DevABE.weight != 2*BytePerTreeLevel {
			t.Errorf("Expected weight %d and got %d instead", 2*BytePerTreeLevel, DevABE.weight)
		}

		if DevAB.weight != BytePerTreeLevel {
			t.Errorf("Expected weight %d and got %d instead", BytePerTreeLevel, DevABE.weight)
		}
	}

}

func TestDevice_updateUUID(t *testing.T) {
	zeroUUID := UUIDType{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	zeroUUIDOut := UUIDType{0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0}
	zeroUUIDConst := UUIDType{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	type fields struct {
		ID            string
		Class         string
		Handle        string
		Product       string
		Description   string
		Vendor        string
		Version       string
		Serial        string
		Slot          string
		uuid          UUIDType
		hash          uint32
		weight        int
		Size          Int64WithUnit
		Capacity      Int64WithUnit
		Capabilities  CapabilitiesList
		Configuration SettingsList
		Devices       []*Device
	}
	type args struct {
		uuid *UUIDType
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		wantUUID           *UUIDType
		wantBytesDifferent int
	}{
		{name: "4-8 bytes, change second 32uint", fields: fields{hash: 2<<31 - 1, weight: 4}, args: args{&zeroUUID}, wantUUID: &zeroUUIDOut, wantBytesDifferent: 0},
		{name: "4-8 bytes, change second 32uint compare to zeros", fields: fields{hash: 2<<31 - 1, weight: 4}, args: args{&zeroUUID}, wantUUID: &zeroUUIDConst, wantBytesDifferent: 12},
		{name: "8-12 bytes", fields: fields{hash: 2<<31 - 1, weight: 8}, args: args{&zeroUUID}, wantUUID: &zeroUUIDConst, wantBytesDifferent: 8},
		{name: "15-19 bytes test, only 1 byte should matter", fields: fields{hash: 2<<31 - 1, weight: 15}, args: args{&zeroUUID}, wantUUID: &zeroUUIDConst, wantBytesDifferent: 1},
	}
	for _, tt := range tests {
		device := &Device{
			ID:            tt.fields.ID,
			Class:         tt.fields.Class,
			Handle:        tt.fields.Handle,
			Product:       tt.fields.Product,
			Description:   tt.fields.Description,
			Vendor:        tt.fields.Vendor,
			Version:       tt.fields.Version,
			Serial:        tt.fields.Serial,
			Slot:          tt.fields.Slot,
			hash:          tt.fields.hash,
			weight:        tt.fields.weight,
			Size:          tt.fields.Size,
			Capacity:      tt.fields.Capacity,
			Capabilities:  tt.fields.Capabilities,
			Configuration: tt.fields.Configuration,
			Devices:       tt.fields.Devices,
		}
		device.updateUUID(tt.args.uuid)

		if tt.args.uuid.BytesDiffer(tt.wantUUID) != tt.wantBytesDifferent {
			t.Errorf("Expected %v UUID and got %v UUID bytes different is %v and it should be %v \n", tt.wantUUID, tt.args.uuid, tt.args.uuid.BytesDiffer(tt.wantUUID), tt.wantBytesDifferent)
		}

		//reset
		zeroUUID = UUIDType{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	}
}

func TestDeviceTree_GetUUID(t *testing.T) {
	type fields struct {
		System *Device
	}
	tests := []struct {
		name    string
		fields  fields
		want    UUIDType
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		deviceTree := &DeviceTree{
			System: tt.fields.System,
		}
		got, err := deviceTree.GetUUID()
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. DeviceTree.GetUUID() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. DeviceTree.GetUUID() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestReportUnmarshallinig(t *testing.T) {

	body, err := ioutil.ReadFile("tst/report_teamster/report_lower_case_n.json")
	report := new(Report)

	if err := json.Unmarshal(body, &report); err != nil {
		t.Error("Failed to unmarshall the report")
		return
	}

	uuid, err := report.System.GetUUID()

	if err != nil {
		t.Error(err)
		return
	}
	uuidString := uuid.ToString()
	uuidStringExpected := "3d217465-7272-824e-7458-ab9ca9a9b985"
	if uuidString != uuidStringExpected {
		t.Errorf("UUID for node is %s got %s instead", uuidStringExpected, uuidString)
	}

}
