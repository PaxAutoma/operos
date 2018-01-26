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
	"encoding/hex"
	"errors"
	"fmt"
)

//it is just an array of bytes
type UUIDType [BytesPerUUID]byte

//BytePerTreeLevel determines how many bits in UUID is allocated per each tree level. 32 bits per tree level
const BytePerTreeLevel = 4

//BytesForMajorComponents determines how many leading bytes ([0:BytesForMajorComponents-1])in the UUID uniqely identify the server
//the rest is used detecting the changes withing a server
const BytesForMajorComponents = 8

//BytesPerUUID is the number of bytes used for specifying a UUID
const BytesPerUUID = 16

//UUIDs is an array of UUIDType, implemented data interface
type UUIDs []UUIDType

func (slice UUIDs) Len() int {
	return len(slice)
}

func (slice UUIDs) Less(i, j int) bool {
	ind := 0
	//skipe the equal parts
	for ; ind < BytesPerUUID && slice[i][ind] == slice[j][ind]; ind++ {

	}

	if ind < BytesPerUUID && slice[i][ind] < slice[j][ind] {
		return true
	}
	return false

}

func (slice UUIDs) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

//BytesDiffer tells from which byte in UUID the difference begins, the higher the less of the impact it has
func (uuid1 *UUIDType) BytesDiffer(uuid2 *UUIDType) int {
	if uuid2 == nil {
		return len(uuid1)
	}

	diffLen := len(uuid1)

	for index, element := range uuid1 {
		if uuid2[index] != element {
			break
		}
		diffLen--
	}

	return diffLen
}

//IsIdenticalTo tells if 2 UUIDs are identical
func (uuid1 *UUIDType) IsIdenticalTo(uuid2 *UUIDType) bool {
	if uuid1.BytesDiffer(uuid2) == 0 {
		return true
	}
	return false
}

//HasTheSameMajorParts tells if there were any changes in the major components
func (uuid1 *UUIDType) HasTheSameMajorParts(uuid2 *UUIDType) bool {
	if uuid1.BytesDiffer(uuid2) <= BytesForMajorComponents {
		return true
	}
	return false
}

//GetUUID returns the UUID for a JSON device tree
func GetUUID(deviceTreeJSON []byte) (UUIDType, error) {

	devTree, err := NewDeviceTree(deviceTreeJSON, "json")

	if err != nil {
		return UUIDType{}, errors.New("Failed to load and parse JSON for the device tree")
	}

	resUUID, err := devTree.GetUUID()

	if err != nil {
		return UUIDType{}, fmt.Errorf("Failed to assign an UUID for device tree %v", deviceTreeJSON)
	}

	return resUUID, nil

}

func (u *UUIDType) ToHexString() string {
	return hex.EncodeToString(u[:])
}

func UUIDTypeFromHexString(b []byte) (*UUIDType, error) {
	uuid := new(UUIDType)

	if n, err := hex.Decode(uuid[:], b); err != nil {
		return nil, err
	} else if n != BytesPerUUID {
		return nil, fmt.Errorf("Failed to decode BytesPerUUID:%i bytes from input", BytesPerUUID)

	}
	return uuid, nil
}

func (u *UUIDType) ToString() string {

	uuid := make([]byte, 16)
	j := 0

	var offest byte = 0
	//basically mapping major bytes to 16 bytes
	//so it does not look like it was concatentaed
	for i := range uuid {
		uuid[i] = u[j] + offest
		if j == BytesForMajorComponents-1 {
			j = 0
			offest += 55
		} else {
			j++
		}
	}

	return string(UUIDFromBytes(uuid))

}

//UUIDFromBytes converts 16bytes into smth like f8e48ce2-6c32-6d81-2f1b-c319a369a4b8
func UUIDFromBytes(u []byte) string {
	var offsets = [...]int{0, 2, 4, 6, 9, 11, 14, 16, 19, 21, 24, 26, 28, 30, 32, 34}

	const hexString = "0123456789abcdef"

	r := make([]byte, 36)
	for i, b := range u {
		r[offsets[i]] = hexString[b>>4]
		r[offsets[i]+1] = hexString[b&0xF]
	}
	r[8] = '-'
	r[13] = '-'
	r[18] = '-'
	r[23] = '-'
	return string(r)
}

//GetEmptyUUID returns a UUID with all zeros
func GetZeroUUID() UUIDType {
	resUUID := UUIDType{}

	for i := 0; i < BytesPerUUID; i++ {
		resUUID[i] = 0
	}
	return resUUID
}
