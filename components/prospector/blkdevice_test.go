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

func TestGenerateUUIDForBlockDevices(t *testing.T) {
	type args struct {
		lsblkDataFile string
		hostUUID      *UUIDType
	}
	tests := []struct {
		name    string
		args    args
		want    *map[string]string
		wantErr bool
	}{
		{name: "sda",
			args: args{
				lsblkDataFile: TSTPath + "/lsblk/test1.json",
				hostUUID:      &UUIDType{30, 127, 75, 45, 178, 152, 193, 79, 224, 58, 43, 77, 252, 126, 88, 186},
			},
			wantErr: false,
			want: &map[string]string{
				"sda": "1e7f4b2d-b298-c14f-adf2-b905447146f7",
			},
		},
		{name: "nmve",
			args: args{
				lsblkDataFile: TSTPath + "/lsblk/test2.json",
				hostUUID:      &UUIDType{30, 127, 75, 45, 178, 152, 193, 79, 224, 58, 43, 77, 252, 126, 88, 186},
			},
			wantErr: false,
			want: &map[string]string{
				"nvme0n1": "1e7f4b2d-b298-c14f-c6cf-dbef7b6b932b",
			},
		},
		{name: "iSCSI",
			args: args{
				lsblkDataFile: TSTPath + "/lsblk/test3.json",
				hostUUID:      &UUIDType{30, 127, 75, 45, 178, 152, 193, 79, 224, 58, 43, 77, 252, 126, 88, 186},
			},
			wantErr: false,
			want: &map[string]string{
				"sda": "1e7f4b2d-b298-c14f-838d-1f6a6b4fe790",
			},
		},

		{name: "mutiple disks",
			args: args{
				lsblkDataFile: TSTPath + "/lsblk/test4.json",
				hostUUID:      &UUIDType{30, 127, 75, 45, 178, 152, 193, 79, 224, 58, 43, 77, 252, 126, 88, 186},
			},
			wantErr: false,
			want: &map[string]string{

				"sdd": "1e7f4b2d-b298-c14f-0211-638615e1ee59",
				"sde": "1e7f4b2d-b298-c14f-b65b-b7e4ad18deee",
				"sda": "1e7f4b2d-b298-c14f-739a-7b3ebdc5dd70",
				"sdb": "1e7f4b2d-b298-c14f-40de-1cb49d82045e",
			},
		},

		{name: "vbox 2",
			args: args{
				lsblkDataFile: TSTPath + "/lsblk/test5.json",
				hostUUID:      &UUIDType{30, 127, 75, 45, 178, 152, 193, 79, 224, 58, 43, 77, 252, 126, 88, 186},
			},
			wantErr: false,
			want: &map[string]string{
				"xvda": "1e7f4b2d-b298-c14f-7e2c-12f37b708fa4",
				"xvdc": "1e7f4b2d-b298-c14f-fc14-554c1c0de31e",
			},
		},

		{name: "vbox 3",
			args: args{
				lsblkDataFile: TSTPath + "/lsblk/test6.json",
				hostUUID:      &UUIDType{30, 127, 75, 45, 178, 152, 193, 79, 224, 58, 43, 77, 252, 126, 88, 186},
			},
			wantErr: false,
			want: &map[string]string{
				"sda": "1e7f4b2d-b298-c14f-d96b-9f45fac84e25",
			},
		},

		{name: "vbox 4",
			args: args{
				lsblkDataFile: TSTPath + "/lsblk/test7.json",
				hostUUID:      &UUIDType{30, 127, 75, 45, 178, 152, 193, 79, 224, 58, 43, 77, 252, 126, 88, 186},
			},
			wantErr: false,
			want: &map[string]string{
				"sda": "1e7f4b2d-b298-c14f-be4d-1a40fb8793bb",
			},
		},

		{name: "vbox 5",
			args: args{
				lsblkDataFile: TSTPath + "/lsblk/test8.json",
				hostUUID:      &UUIDType{30, 127, 75, 45, 178, 152, 193, 79, 224, 58, 43, 77, 252, 126, 88, 186},
			},
			wantErr: false,
			want: &map[string]string{
				"sda": "1e7f4b2d-b298-c14f-bb0b-2f628e69a83a",
			},
		},
	}
	for _, tt := range tests {

		lsblkData, err := ioutil.ReadFile(tt.args.lsblkDataFile)
		if err != nil {
			t.Errorf("Failed to open file %s", tt.args.lsblkDataFile)
		}
		got, err := GenerateUUIDForBlockDevices(lsblkData, tt.args.hostUUID)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. GenerateUUIDForBlockDevices() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. GenerateUUIDForBlockDevices() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
