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

package identity

import (
	"archive/tar"
	"bytes"
	"fmt"
	"strings"

	"github.com/paxautoma/operos/components/teamster/pkg/tarball"
)

const kubeConfigTempl = `apiVersion: v1
kind: Config
clusters:
  - name: {CLUSTER}
    cluster: {API_SERVER}
      certificate-authority: ca.pem
contexts:
  - name: {CLUSTER}
    context:
      cluster: {CLUSTER}
      user: {USER}
users:
  - name: {USER}
    user:
      client-certificate: cert.pem
      client-key: key.pem
current-context: {CLUSTER}
`

const readme = `Operos Kubernetes credentials
-----------------------------

This directory contains the credentials and configuration necessary to connect
to the Kubernetes cluster managed by Operos. To use:

1. Install kubectl (https://kubernetes.io/docs/tasks/tools/install-kubectl/)
2. Run:

	kubectl --kubeconfig=kubeconfig <command>
`

type ClientContext struct {
	Cert      []byte
	Key       []byte
	Bundle    []byte
	InstallID string
	User      string
	ServerURL string
}

var ClientManifest = tarball.Manifest{
	tarball.ManifestFile{
		tar.Header{
			Name: "operos-credentials/cert.pem",
			Mode: 0600,
		},
		func(ctx interface{}, data *bytes.Buffer) error {
			info := ctx.(ClientContext)
			data.Write(info.Cert)
			return nil
		},
	},
	tarball.ManifestFile{
		tar.Header{
			Name: "operos-credentials/key.pem",
			Mode: 0600,
		},
		func(ctx interface{}, data *bytes.Buffer) error {
			info := ctx.(ClientContext)
			data.Write(info.Key)
			return nil
		},
	},
	tarball.ManifestFile{
		tar.Header{
			Name: "operos-credentials/ca.pem",
			Mode: 0600,
		},
		func(ctx interface{}, data *bytes.Buffer) error {
			info := ctx.(ClientContext)
			data.Write(info.Bundle)
			return nil
		},
	},
	tarball.ManifestFile{
		tar.Header{
			Name: "operos-credentials/kubeconfig",
			Mode: 0644,
		},
		func(ctx interface{}, data *bytes.Buffer) error {
			info := ctx.(ClientContext)
			kc := strings.Replace(kubeConfigTempl, "{CLUSTER}", info.InstallID, -1)
			kc = strings.Replace(kc, "{USER}", info.User, -1)
			if info.ServerURL != "" {
				kc = strings.Replace(kc, "{API_SERVER}", fmt.Sprintf("\n      server: %s", info.ServerURL), -1)
			}
			data.WriteString(kc)
			return nil
		},
	},
	tarball.ManifestFile{
		tar.Header{
			Name: "operos-credentials/README.txt",
			Mode: 0644,
		},
		func(ctx interface{}, data *bytes.Buffer) error {
			data.WriteString(readme)
			return nil
		},
	},
}
