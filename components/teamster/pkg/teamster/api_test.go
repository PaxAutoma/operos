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

package teamster

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	fmt "fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/coreos/etcd/clientv3"
	"github.com/paxautoma/operos/components/teamster/pkg/cluster"
)

var (
	clusterName string
	etcdCluster string
)

func init() {
	clusterName = os.Getenv("CLUSTER_NAME")
	etcdCluster = fmt.Sprintf("%s:2379", os.Getenv("ETCD_IP"))
}

func setupAPI() (*TeamsterAPI, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdCluster},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		return nil, errors.Wrap(err, "could not connect to etcd")
	}

	oc, err := cluster.InstantiateCluster(client, 5*time.Second, clusterName)

	if err != nil {
		return nil, errors.Wrap(err, "could not instantiate cluster object")
	}

	return NewTeamsterAPI(oc, "../../acceptance-test/data/shadow", "root"), nil
}

func TestGenClientCert(t *testing.T) {
	api, err := setupAPI()
	require.NoError(t, err)

	t.Run("ValidInput_GeneratesValidOutput", func(t *testing.T) {
		// Make the request
		req, err := http.NewRequest("GET", "/clientcert?user=mytestuser&group=mytestgroup", nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		handler := api.GetHttpHandler()
		handler.ServeHTTP(rr, req)

		// General response checks
		require.Equal(t, http.StatusOK, rr.Code)
		resp, err := readTarball(rr.Body)
		require.NoError(t, err)
		require.Contains(t, resp, "operos-credentials/cert.pem")
		require.Contains(t, resp, "operos-credentials/key.pem")
		require.Contains(t, resp, "operos-credentials/ca.pem")
		require.Contains(t, resp, "operos-credentials/kubeconfig")

		// Check certificate
		certBlock, _ := pem.Decode(resp["operos-credentials/cert.pem"])
		require.NotNil(t, certBlock)
		cert, err := x509.ParseCertificate(certBlock.Bytes)
		require.NoError(t, err)

		require.False(t, cert.IsCA)
		require.Contains(t, cert.ExtKeyUsage, x509.ExtKeyUsageClientAuth)
		require.Equal(t, "mytestuser", cert.Subject.CommonName)
		require.Contains(t, cert.Issuer.CommonName, "systest-teamster")

		name := pkix.Name{}
		name.FillFromRDNSequence(&pkix.RDNSequence{cert.Subject.Names})
		require.Contains(t, name.Organization, "mytestgroup")

		// Check key
		keyBlock, _ := pem.Decode(resp["operos-credentials/key.pem"])
		require.NotNil(t, certBlock)
		key, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
		require.NoError(t, err)
		require.NoError(t, key.Validate())
	})

	t.Run("MissingUsername_ReturnsBadRequest", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/clientcert?group=mytestgroup", nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(api.GenClientCert)
		handler.ServeHTTP(rr, req)

		require.Equal(t, rr.Code, http.StatusBadRequest)
	})

	t.Run("MissingGroup_ReturnsBadRequest", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/clientcert?user=mytestuser", nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(api.GenClientCert)
		handler.ServeHTTP(rr, req)

		require.Equal(t, rr.Code, http.StatusBadRequest)
	})
}

func TestWhoami(t *testing.T) {
	api, err := setupAPI()
	require.NoError(t, err)

	t.Run("ValidInput_GeneratesValidOutput", func(t *testing.T) {
		body, err := ioutil.ReadFile("../../acceptance-test/data/node001.json")
		require.NoError(t, err)

		// Make the request
		req, err := http.NewRequest("POST", "/whoami", bytes.NewReader(body))
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		handler := api.GetHttpHandler()
		handler.ServeHTTP(rr, req)

		// General response checks
		require.Equal(t, http.StatusOK, rr.Code)
		resp, err := readTarball(rr.Body)
		require.NoError(t, err)

		require.Contains(t, resp, "etc/paxautoma/settings")
		require.Contains(t, resp, "etc/kubernetes/ssl/ca.pem")
		require.Contains(t, resp, "etc/kubernetes/ssl/worker.pem")
		require.Contains(t, resp, "etc/kubernetes/ssl/worker-key.pem")
		require.Contains(t, resp, "etc/keyfile")
		require.Contains(t, resp, "etc/hostname")
		require.Contains(t, resp, "root/.ssh/authorized_keys")
		require.Contains(t, resp, "/etc/rootpasshash")

		keys := strings.Split(string(resp["root/.ssh/authorized_keys"]), "\n")
		require.Contains(t, keys, "first test key")
		require.Contains(t, keys, "second test key")
	})
}

func readTarball(buf *bytes.Buffer) (result map[string][]byte, err error) {
	gzReader, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}

	tarReader := tar.NewReader(gzReader)

	result = make(map[string][]byte)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		body, err := ioutil.ReadAll(tarReader)
		if err != nil {
			return nil, err
		}

		result[header.Name] = body
	}

	return result, nil
}
