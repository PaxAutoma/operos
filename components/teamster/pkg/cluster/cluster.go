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

package cluster

import (
	"context"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/cloudflare/cfssl/config"
	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/helpers"
	"github.com/cloudflare/cfssl/signer"
	"github.com/cloudflare/cfssl/signer/local"
	"github.com/coreos/etcd/clientv3"
	"github.com/pkg/errors"
	"github.com/paxautoma/operos/components/prospector"
)

type NodeOSD struct {
	Id     string
	Key    string
	Weight string
}

type Node struct {
	Id                 string
	Fingerprint        *prospector.UUIDType
	LatestReport       *prospector.Report
	KubeletPrivateKey  []byte
	KubeletCertificate []byte
	LuksKeyFile        []byte
	Cluster            *OperosCluster
	OSDs               map[string]*NodeOSD
}

type OperosCluster struct {
	InstallID            string
	Signer               signer.Signer
	CA_Bundle            []byte
	CACert               []byte
	Nodes                map[string]*Node
	Vars                 map[string]string
	WorkerAuthorizedKeys [][]byte
	etcd                 *clientv3.Client
	etcdRequestTimeout   time.Duration
	CephConfig           []byte
	Secrets              map[string][]byte
}

func (cluster *OperosCluster) loadNode(nodeid string) *Node {
	ctx, cancel := context.WithTimeout(context.Background(), cluster.etcdRequestTimeout)
	resp, err := cluster.etcd.Get(ctx, fmt.Sprintf("nodes/%s/%s/", cluster.InstallID, nodeid), clientv3.WithPrefix())
	cancel()
	if err != nil {
		log.Printf("error: Unable to get node vars from etcd", err)
		return nil
	}

	node := new(Node)
	node.Id = nodeid
	node.OSDs = make(map[string]*NodeOSD)

	for _, ev := range resp.Kvs {
		keyv := strings.Split(string(ev.Key), "/")
		switch keyv[3] {
		case "latestreport":
			report := new(prospector.Report)
			if err = json.Unmarshal(ev.Value, &report); err != nil {
				log.Printf("unable to unmarshal latest report from prospector: %s", err)
			}
			node.LatestReport = report
		case "secret-kubelet-key":
			node.KubeletPrivateKey = ev.Value
		case "secret-kubelet-cert":
			node.KubeletCertificate = ev.Value
		case "secret-luks-keyfile":
			node.LuksKeyFile = ev.Value
		case "fingerprint":
			if node.Fingerprint, err = prospector.UUIDTypeFromHexString(ev.Value); err != nil {
				log.Printf("unable to decode node fingerprint: %s", err)
			}
		case "osd":
			osd_uuid := keyv[4]
			var osd *NodeOSD
			var exists bool
			if osd, exists = node.OSDs[osd_uuid]; !exists {
				osd = new(NodeOSD)
				node.OSDs[osd_uuid] = osd
			}
			switch keyv[5] {
			case "Id":
				osd.Id = string(ev.Value)
			case "Key":
				osd.Key = string(ev.Value)
			case "Weight":
				osd.Weight = string(ev.Value)
			}
		}
	}

	node.Cluster = cluster
	cluster.Nodes[node.Id] = node
	return node
}

func (cluster *OperosCluster) storeNode(node *Node) error {
	ctx, cancel := context.WithTimeout(context.Background(), cluster.etcdRequestTimeout)
	node_key := fmt.Sprintf("nodes/%s/%s", cluster.InstallID, node.Id)

	if _, ok := cluster.Nodes[node.Id]; ok {
		return errors.New(fmt.Sprintf("Node %s already exists in the cluster %s", cluster.InstallID, node.Id))
	}

	//XXX: TODO: wrap this in a transaction

	serialized_report, err := json.Marshal(node.LatestReport)
	if err != nil {
		return err
	}

	_, err = cluster.etcd.Put(ctx, fmt.Sprintf("%s/%s", node_key, "latestreport"), string(serialized_report))
	if err != nil {
		return err
	}

	_, err = cluster.etcd.Put(ctx, fmt.Sprintf("%s/%s", node_key, "fingerprint"), node.Fingerprint.ToHexString())
	if err != nil {
		return err
	}

	_, err = cluster.etcd.Put(ctx, fmt.Sprintf("%s/%s", node_key, "secret-kubelet-key"), string(node.KubeletPrivateKey))

	if err != nil {
		return err
	}

	_, err = cluster.etcd.Put(ctx, fmt.Sprintf("%s/%s", node_key, "secret-kubelet-cert"), string(node.KubeletCertificate))

	if err != nil {
		return err
	}

	_, err = cluster.etcd.Put(ctx, fmt.Sprintf("%s/%s", node_key, "secret-luks-keyfile"), string(node.LuksKeyFile))

	if err != nil {
		return err
	}

	_, err = cluster.etcd.Delete(ctx, fmt.Sprintf("%s/osd", node_key), clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for osd_uuid, osd := range node.OSDs {
		_, err = cluster.etcd.Put(ctx, fmt.Sprintf("%s/osd/%s/Id", node_key, osd_uuid), string(osd.Id))

		if err != nil {
			return err
		}

		_, err = cluster.etcd.Put(ctx, fmt.Sprintf("%s/osd/%s/Weight", node_key, osd_uuid), string(osd.Weight))

		if err != nil {
			return err
		}

		_, err = cluster.etcd.Put(ctx, fmt.Sprintf("%s/osd/%s/Key", node_key, osd_uuid), string(osd.Key))

		if err != nil {
			return err
		}
	}

	keys := make([]string, 0, len(cluster.Nodes)+1)
	for k := range cluster.Nodes {
		keys = append(keys, k)
	}

	keys = append(keys, node.Id)
	sort.Strings(keys)

	nodeList := strings.Join(keys, ",")
	_, err = cluster.etcd.Put(ctx, fmt.Sprintf("cluster/%s/nodeids", cluster.InstallID), nodeList)

	if err != nil {
		return err
	}

	cluster.Nodes[node.Id] = node

	cancel()
	return nil
}

func CephCreateOsd(uuid string) (string, error) {
	cmdName := "/usr/bin/ceph"
	cmdArgs := []string{"osd", "create", uuid}

	cmdOut, err := exec.Command(cmdName, cmdArgs...).Output()
	return string(cmdOut), err
}

func CephCreateOSDKey(osd_name string) error {
	cmdName := "/usr/bin/ceph"
	cmdArgs := []string{"auth", "add", osd_name, "osd", "allow *", "mon", "allow rwx"}

	_, err := exec.Command(cmdName, cmdArgs...).Output()
	return err
}

func CephOSDGetKey(osd_name string) (string, error) {
	cmdName := "/usr/bin/ceph"
	cmdArgs := []string{"auth", "get-key", osd_name}

	cmdOut, err := exec.Command(cmdName, cmdArgs...).Output()
	return string(cmdOut), err
}

func CephOSDCrushAdd(osd_name string, weight string, hostname string) (string, error) {
	cmdName := "/usr/bin/ceph"
	cmdArgs := []string{"osd", "crush", "add", osd_name, weight, fmt.Sprintf("host=%s", hostname)}

	cmdOut, err := exec.Command(cmdName, cmdArgs...).Output()
	return string(cmdOut), err
}

func CephAddHostToCrush(host string) error {
	cmdName := "/usr/bin/ceph"
	cmdArgs := []string{"osd", "crush", "add-bucket", host, "host"}

	_, err := exec.Command(cmdName, cmdArgs...).Output()
	return err
}

func CephCrushMoveHostDefault(host string) error {
	cmdName := "/usr/bin/ceph"
	cmdArgs := []string{"osd", "crush", "move", host, "root=default"}

	_, err := exec.Command(cmdName, cmdArgs...).Output()
	return err
}

func CephOSDPurge(osd_name string) error {
	cmdName := "/usr/bin/ceph"
	cmdArgs := []string{"osd", "purge", osd_name, "--yes-i-really-mean-it"}

	_, err := exec.Command(cmdName, cmdArgs...).Output()
	return err
}

func (node *Node) AddOSD(osd_uuid string, osd *NodeOSD) error {
	var osd_id string
	var osd_key string
	if raw_osd_id, err := CephCreateOsd(osd_uuid); err != nil {
		log.Printf("Failed to create osd: %s", raw_osd_id, err)
		return err
	} else {
		osd_id = strings.TrimSpace(raw_osd_id)
	}

	osd.Id = osd_id
	osd_name := fmt.Sprintf("osd.%s", osd_id)

	if err := CephCreateOSDKey(osd_name); err != nil {
		log.Printf("Failed to genereate a key for osd: %s", err)
		return err
	}

	if raw_osd_key, err := CephOSDGetKey(osd_name); err != nil {
		log.Printf("Failed to get key for osd: (%s) %s", raw_osd_key, err)
		return err
	} else {
		osd_key = strings.TrimSpace(raw_osd_key)
	}

	osd.Key = osd_key

	if out, err := CephOSDCrushAdd(osd_name, osd.Weight, node.Id); err != nil {
		log.Printf("Adding osd %s as %s with weight %s on %s to crushmap failed; will not be active: %s %s", osd_uuid, osd_name, osd.Weight, node.Id, err, out)
	}

	return nil
}

func (node *Node) InventoryOSDs(blkdevices []*prospector.BlockDevice) (map[string]*NodeOSD, error) {
	osds := make(map[string]*NodeOSD)

	for _, blkDevice := range blkdevices {
		if blkDevice.Type == "disk" {
			var uuid *string
			var err error
			if uuid, err = prospector.UUIDStringForBlkDevice(blkDevice, node.Fingerprint); err != nil {
				log.Printf("Failed to compute UUID for %s because of %s", blkDevice.Name, err)
			}
			osd := new(NodeOSD)
			if sz_bytes, err := strconv.ParseFloat(blkDevice.Size, 64); err != nil {
				log.Printf("Failed to calculate weight for OSD setting to zero in_b:%s", blkDevice.Size)
				osd.Weight = "0.0"
			} else {
				osd.Weight = fmt.Sprintf("%f", sz_bytes/float64(1000000000000))
			}
			osds[*uuid] = osd
		}

	}
	return osds, nil
}

func (cluster *OperosCluster) AddNode(id *prospector.UUIDType, uuid string, report *prospector.Report) (*Node, error) {
	node := new(Node)
	node.Id = uuid
	node.Fingerprint = id
	node.LatestReport = report

	log.Printf("Adding node %s to cluster %s", node.Id, cluster.InstallID)

	cn := fmt.Sprintf("Operos Cluster (%s) Node (%s)", cluster.InstallID, uuid)
	groups := []string{cluster.Vars["OPEROS_CLUSTER_ORG"]}
	c, p, err := cluster.requestAndSign(cn, groups)
	if err != nil {
		return nil, err
	}

	node.KubeletCertificate = c
	node.KubeletPrivateKey = p

	keyfile, err := cluster.generateLuksKeyFile()
	if err != nil {
		return nil, err
	}
	node.LuksKeyFile = keyfile
	CephAddHostToCrush(node.Id)
	CephCrushMoveHostDefault(node.Id)

	if node.OSDs, err = node.InventoryOSDs(node.LatestReport.Storage.BlockDevices); err != nil {
		log.Printf("Unable to inventory storage from node %s to cluster %s: %s", node.Id, cluster.InstallID, err)
	} else {
		for osd_uuid, osd := range node.OSDs {
			if err = node.AddOSD(osd_uuid, osd); err != nil {
				log.Printf("Faild to add OSD:%s from node %s to cluster %s: %s", osd_uuid, node.Id, cluster.InstallID, err)
			}
		}
	}
	err = cluster.storeNode(node)
	if err != nil {
		log.Printf("storing node %s failed: %s", uuid, err)
		return nil, err
	}

	node.Cluster = cluster

	return node, nil
}

func (cluster *OperosCluster) requestAndSign(cn string, o []string) ([]byte, []byte, error) {
	req := csr.New()
	req.CN = cn
	req.Names = make([]csr.Name, len(o)+1)
	req.Names[0] = csr.Name{
		C:  cluster.Vars["OPEROS_CLUSTER_COUNTRY"],
		L:  cluster.Vars["OPEROS_CLUSTER_CITY"],
		OU: cluster.Vars["OPEROS_CLUSTER_ORG"],
		ST: cluster.Vars["OPEROS_CLUSTER_PROVINCE"],
	}
	for idx, org := range o {
		req.Names[idx+1] = csr.Name{
			O: org,
		}
	}

	keyr := csr.NewBasicKeyRequest()
	keyr.A = "rsa"
	keyr.S = 2048

	req.KeyRequest = keyr

	var key, csrBytes []byte
	csrBytes, key, err := csr.ParseRequest(req)
	if err != nil {
		return nil, nil, err
	}

	signReq := signer.SignRequest{
		Request: string(csrBytes),
	}

	var cert []byte
	cert, err = cluster.Signer.Sign(signReq)
	if err != nil {
		return nil, nil, err
	}
	return cert, key, nil
}

func (node *Node) RemoveOSD(osd_uuid string, osd *NodeOSD) error {
	osd_name := fmt.Sprintf("osd.%s", osd.Id)

	if err := CephOSDPurge(osd_name); err != nil {
		return err
	}

	return nil
}

func (cluster *OperosCluster) UpdateNode(node *Node, id *prospector.UUIDType, uuid string, report *prospector.Report) error {
	if osds, err := node.InventoryOSDs(report.Storage.BlockDevices); err != nil {
		log.Printf("Unable to inventory storage from node %s in cluster %s: %s", node.Id, cluster.InstallID, err)
	} else {
		// check for a removed OSD
		for osd_uuid, osd := range node.OSDs {
			if _, exists := osds[osd_uuid]; !exists {
				if err := node.RemoveOSD(osd_uuid, osd); err != nil {
					log.Printf("Faild to remove OSD:%s from node %s in cluster %s: %s", osd_uuid, node.Id, cluster.InstallID, err)
				}
				delete(node.OSDs, osd_uuid)
			}
		}
		// check for new OSDs
		for osd_uuid, osd := range osds {
			if _, exists := node.OSDs[osd_uuid]; !exists {
				if err = node.AddOSD(osd_uuid, osd); err == nil {
					node.OSDs[osd_uuid] = osd
				} else {
					log.Printf("Faild to add OSD:%s from node %s to cluster %s: %s", osd_uuid, node.Id, cluster.InstallID, err)
				}
			}
		}

	}

	// -- update node last request field
	// -- check node certificate validity and update
	node.LatestReport = report
	cluster.storeNode(node)
	return nil
}

func (cluster *OperosCluster) AddUser(user string, groups []string) ([]byte, []byte, error) {
	return cluster.requestAndSign(user, groups)
}

func (cluster *OperosCluster) generateLuksKeyFile() ([]byte, error) {
	keym := make([]byte, 512) // 4096 bit key
	_, err := rand.Read(keym)
	if err != nil {
		return nil, err
	}
	return keym, nil
}

func operosSigner(certificate []byte, key []byte) (signer.Signer, error) {
	d := helpers.OneYear
	policy := &config.Signing{
		Profiles: map[string]*config.SigningProfile{},
		Default: &config.SigningProfile{
			Usage:        []string{"digital signature", "client auth"},
			Expiry:       d,
			ExpiryString: "8760h",
		},
	}

	parsedCa, err := helpers.ParseCertificatePEM(certificate)
	if err != nil {
		log.Printf("operosSigner: Malformed certificate %v", err)
		return nil, err
	}

	priv, err := helpers.ParsePrivateKeyPEMWithPassword(key, nil)
	if err != nil {
		log.Printf("operosSigner: Malformed private key %v", err)
		return nil, err
	}
	return local.NewSigner(priv, parsedCa, signer.DefaultSigAlgo(priv), policy)
}

func InstantiateCluster(etcd *clientv3.Client, requestTimeout time.Duration, installID string) (*OperosCluster, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)

	oc := new(OperosCluster)
	oc.etcd = etcd
	oc.etcdRequestTimeout = requestTimeout
	oc.Nodes = make(map[string]*Node)
	oc.Vars = make(map[string]string)
	oc.Secrets = make(map[string][]byte)
	resp, err := oc.etcd.Get(ctx, fmt.Sprintf("cluster/%s/", installID), clientv3.WithPrefix())
	cancel()
	if err != nil {
		log.Println("error: Unable to get cluster vars from etcd", err)
		return nil, err
	}
	oc.InstallID = installID
	var ca_key []byte
	var ca_cert []byte
	var ca_bundle []byte
	var nodeidar string

	for _, ev := range resp.Kvs {
		keyv := strings.Split(string(ev.Key), "/")
		switch keyv[2] {
		case "secret-ca-key":
			ca_key = ev.Value
			oc.Secrets[keyv[2]] = ev.Value
		case "secret-ca-cert":
			ca_cert = ev.Value
			oc.Secrets[keyv[2]] = ev.Value
		case "secret-ca-bundle":
			ca_bundle = ev.Value
			oc.Secrets[keyv[2]] = ev.Value
		case "nodeids":
			nodeidar = strings.Trim(string(ev.Value), ",")
		case "authorized-keys":
			if keyv[3] == "worker" {
				oc.WorkerAuthorizedKeys = append(oc.WorkerAuthorizedKeys, ev.Value)
			}
		case "ceph-config":
			oc.CephConfig = ev.Value
		default:
			if strings.HasPrefix(keyv[2], "secret") {
				oc.Secrets[keyv[2]] = ev.Value
			} else {
				oc.Vars[keyv[2]] = string(ev.Value)
			}
		}
	}

	if len(ca_key) == 0 {
		log.Println("error: cluster Certificate Authority Key unconfigured")
		return nil, errors.New("Certificate Authority Key unconfigured")
	}

	if len(ca_cert) == 0 {
		log.Println("error: cluster Certificate Authority Certificate unconfigured")
		return nil, errors.New("Certificate Authority Certificate unconfigured")
	}

	oc.Signer, err = operosSigner(ca_cert, ca_key)
	if err != nil {
		log.Println("error: unable to initilize signer: ", err)
		return nil, err
	}

	oc.CA_Bundle = ca_bundle
	oc.CACert = ca_cert

	if len(nodeidar) > 0 {
		nodeids := strings.Split(string(nodeidar), ",")
		for _, nodeid := range nodeids {
			node := oc.loadNode(nodeid)
			if node == nil {
				log.Println("error: unable to load node %s belonging to cluster %s", nodeid, installID)
			}
		}
	}

	return oc, nil
}

func (cluster *OperosCluster) GetCACert() (*x509.Certificate, error) {
	caPEM := cluster.CACert

	block, _ := pem.Decode(caPEM)
	if block == nil {
		return nil, errors.Errorf("could not decode certificate PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse certificate")
	}

	return cert, nil
}
