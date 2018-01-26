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
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/paxautoma/operos/components/teamster/pkg/cluster"
	"github.com/paxautoma/operos/components/teamster/pkg/tarball"

	"github.com/pkg/errors"
)

type WorkerContext struct {
	Node        *cluster.Node
	Cluster     *cluster.OperosCluster
	ShadowFile  string
	RootAccount string
}

var WorkerManifest = tarball.Manifest{
	tarball.ManifestFile{
		tar.Header{
			Name: "etc/paxautoma/settings",
			Mode: 0600,
		},
		ClusterSettings,
	},
	tarball.ManifestFile{
		tar.Header{
			Name: "etc/kubernetes/ssl/ca.pem",
			Mode: 0600,
		},
		KubernetesCertificateAuthorityCert,
	},
	tarball.ManifestFile{
		tar.Header{
			Name: "etc/kubernetes/ssl/worker.pem",
			Mode: 0600,
		},
		KubernetesWorkerCertificate,
	},
	tarball.ManifestFile{
		tar.Header{
			Name: "etc/kubernetes/ssl/worker-key.pem",
			Mode: 0600,
		},
		KubernetesWorkerKey,
	},
	tarball.ManifestFile{
		tar.Header{
			Name: "etc/keyfile",
			Mode: 0600,
		},
		LuksKeyFile,
	},
	tarball.ManifestFile{
		tar.Header{
			Name: "etc/hostname",
			Mode: 0644,
		},
		HostnameFile,
	},
	tarball.ManifestFile{
		tar.Header{
			Name: "root/.ssh/authorized_keys",
			Mode: 0600,
		},
		AuthorizedKeysFile,
	},
	tarball.ManifestFile{
		tar.Header{
			Name: "etc/rootpasshash",
			Mode: 0600,
		},
		RootPassHashFile,
	},
	tarball.ManifestFile{
		tar.Header{
			Name: "etc/ceph/ceph.conf",
			Mode: 0600,
		},
		CephConfig,
	},
	tarball.ManifestFile{
		tar.Header{
			Name: "etc/paxautoma/osd-loadout",
			Mode: 0600,
		},
		OSDLoadout,
	},
	tarball.ManifestFile{
		tar.Header{
			Name: "etc/ceph/keyring",
			Mode: 0600,
		},
		KubeCephKeyring,
	},
}

func KubeCephKeyring(ctx interface{}, data *bytes.Buffer) error {
	info := ctx.(*WorkerContext)
	data.Write(info.Node.Cluster.Secrets["secret-ceph-kube-keyring"])
	return nil
}

func OSDLoadout(ctx interface{}, data *bytes.Buffer) error {
	info := ctx.(*WorkerContext)

	for bdev_uuid, osd := range info.Node.OSDs {
		if osd.Id != "" {
			data.WriteString(fmt.Sprintf("%s:%s:%s\n", bdev_uuid, osd.Id, osd.Key))
		}
	}

	return nil
}

func CephConfig(ctx interface{}, data *bytes.Buffer) error {
	info := ctx.(*WorkerContext)
	data.Write(info.Node.Cluster.CephConfig)
	return nil
}

func LuksKeyFile(ctx interface{}, data *bytes.Buffer) error {
	info := ctx.(*WorkerContext)
	data.Write(info.Node.LuksKeyFile)
	return nil
}

func KubernetesCertificateAuthorityCert(ctx interface{}, data *bytes.Buffer) error {
	info := ctx.(*WorkerContext)
	if info.Node.Cluster.CA_Bundle == nil {
		return errors.New("Cluster has not been certified")
	}

	data.Write(info.Node.Cluster.CA_Bundle)
	return nil
}

func KubernetesWorkerCertificate(ctx interface{}, data *bytes.Buffer) error {
	info := ctx.(*WorkerContext)
	if info.Node.KubeletCertificate == nil {
		return errors.New("Node has not been certified")
	}
	data.Write(info.Node.KubeletCertificate)
	return nil
}

func KubernetesWorkerKey(ctx interface{}, data *bytes.Buffer) error {
	info := ctx.(*WorkerContext)
	if info.Node.KubeletPrivateKey == nil {
		return errors.New("Node has not been certified")
	}
	data.Write(info.Node.KubeletPrivateKey)
	return nil
}

func ClusterSettings(ctx interface{}, data *bytes.Buffer) error {
	info := ctx.(*WorkerContext)
	for key, value := range info.Node.Cluster.Vars {
		data.WriteString(fmt.Sprintf("%s=\"%s\"\n", key, value))
	}
	return nil
}

func HostnameFile(ctx interface{}, data *bytes.Buffer) error {
	info := ctx.(*WorkerContext)
	data.WriteString(info.Node.Id)
	return nil
}

func AuthorizedKeysFile(ctx interface{}, data *bytes.Buffer) error {
	info := ctx.(*WorkerContext)
	for _, key := range info.Cluster.WorkerAuthorizedKeys {
		data.Write(key)
		if key[len(key)-1] != '\n' {
			data.WriteByte('\n')
		}
	}
	return nil
}

func RootPassHashFile(ctx interface{}, data *bytes.Buffer) error {
	info := ctx.(*WorkerContext)

	shadow, err := os.Open(info.ShadowFile)
	if err != nil {
		return errors.Wrap(err, "couldn't open shadow file")
	}

	scanner := bufio.NewScanner(shadow)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return errors.Wrap(err, "couldn't read shadow file")
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Split(line, ":")
		if parts[0] == info.RootAccount {
			data.Write([]byte(parts[1]))
			return nil
		}
	}

	return errors.Errorf("shadow file does not contain account %s", info.RootAccount)
}
