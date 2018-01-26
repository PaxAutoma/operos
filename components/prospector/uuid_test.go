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
	"io/ioutil"
	"reflect"
	"testing"
)

func TestUUIDType_IsIdenticalTo(t *testing.T) {

	uuidA := UUIDType{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	uuidAA := UUIDType{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

	uuidB := UUIDType{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0}

	type args struct {
		uuid2 *UUIDType
	}
	tests := []struct {
		name  string
		uuid1 *UUIDType
		args  args
		want  bool
	}{
		{name: "equal test", uuid1: &uuidA, args: args{uuid2: &uuidAA}, want: true},
		{name: "same struct test", uuid1: &uuidA, args: args{uuid2: &uuidA}, want: true},
		{name: "same struct test", uuid1: &uuidA, args: args{uuid2: &uuidB}, want: false},
		{name: "same struct test", uuid1: &uuidA, args: args{uuid2: nil}, want: false},
	}
	for _, tt := range tests {
		if got := tt.uuid1.IsIdenticalTo(tt.args.uuid2); got != tt.want {
			t.Errorf("%q. UUIDType.IsIdenticalTo() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestUUIDType_BytesDiffer(t *testing.T) {
	uuidA := UUIDType{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	uuidAA := UUIDType{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

	uuidB := UUIDType{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0}

	uuidC := UUIDType{1, 22, 23, 24, 25, 26, 27, 28, 9, 10, 11, 12, 13, 14, 15, 0}

	type args struct {
		uuid2 *UUIDType
	}
	tests := []struct {
		name  string
		uuid1 *UUIDType
		args  args
		want  int
	}{
		{name: "equal test", uuid1: &uuidA, args: args{uuid2: &uuidAA}, want: 0},
		{name: "same struct test", uuid1: &uuidA, args: args{uuid2: &uuidA}, want: 0},
		{name: "different struct test", uuid1: &uuidA, args: args{uuid2: &uuidB}, want: 1},
		{name: "nil struct test", uuid1: &uuidA, args: args{uuid2: nil}, want: 16},
		{name: "half different struct test", uuid1: &uuidA, args: args{uuid2: &uuidC}, want: 15},
	}
	for _, tt := range tests {
		if got := tt.uuid1.BytesDiffer(tt.args.uuid2); got != tt.want {
			t.Errorf("%q. UUIDType.BytesDiffer() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_GetUUID(t *testing.T) {

	testFileName := TSTPath + "/json/laptop.json"
	data, err := ioutil.ReadFile(testFileName)

	if err != nil {
		t.Errorf("Could not load the file %s", testFileName)
	}

	type args struct {
		deviceTreeJSON []byte
	}
	tests := []struct {
		name    string
		args    args
		want    UUIDType
		wantErr bool
	}{
		{name: "laptop", args: args{deviceTreeJSON: data}, want: UUIDType{30, 127, 75, 45, 178, 152, 193, 79, 224, 58, 43, 77, 252, 126, 88, 186}, wantErr: false},
	}
	for _, tt := range tests {
		got, err := GetUUID(tt.args.deviceTreeJSON)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. getUUID() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. getUUID() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestUUIDToStringMinorChanges(t *testing.T) {
	uuid1 := UUIDType{12, 140, 155, 235, 96, 1, 124, 83, 224, 58, 43, 77, 181, 146, 16, 230}

	uuid2 := uuid1

	uuid2[BytesForMajorComponents] = uuid1[BytesForMajorComponents] * 2

	if !uuid1.HasTheSameMajorParts(&uuid2) {
		t.Error("The UUIDs should have the same major parts")
	}

	uuid1String := uuid1.ToString()
	uuid2String := uuid2.ToString()

	if uuid1String != uuid2String {
		t.Error("The UUIDs with the same major parts should have the same string representation")
	}

}

func TestUUIDToStringMajorChanges(t *testing.T) {
	uuid1 := UUIDType{12, 140, 155, 235, 96, 1, 124, 83, 224, 58, 43, 77, 181, 146, 16, 230}

	uuid2 := uuid1

	uuid2[BytesForMajorComponents-1] = uuid1[BytesForMajorComponents-1] * 2

	if uuid1.HasTheSameMajorParts(&uuid2) {
		t.Error("The UUIDs should not have the same major parts")
	}

	uuid1String := uuid1.ToString()
	uuid2String := uuid2.ToString()

	if uuid1String == uuid2String {
		t.Error("The UUIDs with the same major parts should not have the same string representation")
	}
}

func TestUUIDToStringEntropy(t *testing.T) {
	uuid1 := UUIDType{12, 140, 155, 235, 96, 1, 124, 83, 224, 58, 43, 77, 181, 146, 16, 230}
	uuid1String := uuid1.ToString()

	for i, v := range uuid1[:BytesForMajorComponents] {
		uuid2 := uuid1

		uuid2[i] = v + 10

		uuid2String := uuid2.ToString()

		diff := 0

		for j := range uuid2String {
			if uuid1String[j] != uuid2String[j] {
				diff = j
				break
			}
		}

		t.Logf("%s %s %d", uuid1String, uuid2String, diff)

		if diff < i {
			t.Error("The UUIDS should differ more. ")
		}

	}

}
