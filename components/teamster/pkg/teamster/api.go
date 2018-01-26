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
	"bufio"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	crypt "github.com/amoghe/go-crypt"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/paxautoma/operos/components/prospector"
	"github.com/paxautoma/operos/components/teamster/pkg/cluster"
	"github.com/paxautoma/operos/components/teamster/pkg/identity"
	"github.com/paxautoma/operos/components/teamster/pkg/tarball"
)

type TeamsterAPI struct {
	cluster     *cluster.OperosCluster
	shadowFile  string
	rootAccount string
}

func NewTeamsterAPI(c *cluster.OperosCluster, shadowFile, rootAccount string) *TeamsterAPI {
	return &TeamsterAPI{c, shadowFile, rootAccount}
}

func (t *TeamsterAPI) GetHttpHandler() http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	router.
		Methods("POST").
		Path("/whoami").
		Name("whoami").
		Handler(http.HandlerFunc(t.Whoami))
	router.
		Methods("GET").
		Path("/clientcert").
		Name("clientcert").
		Handler(http.HandlerFunc(t.GenClientCert))

	return router
}

func (t *TeamsterAPI) RegisterGRPCService(grpcServer *grpc.Server) {
	RegisterTeamsterServer(grpcServer, t)
}

func (t *TeamsterAPI) Whoami(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Println(err)
		return
	}

	if err := r.Body.Close(); err != nil {
		log.Println(err)
		return
	}

	report := new(prospector.Report)

	if err := json.Unmarshal(body, &report); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnprocessableEntity)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Println(err)
		}
		return
	}

	uuid, err := report.System.GetUUID()
	if err != nil {
		log.Println(err)
		return
	}

	uuidString := uuid.ToString()

	log.Printf("UUID for node is %s", uuidString)

	node, exists := t.cluster.Nodes[uuidString]

	if !exists {
		log.Printf("%s does not exist: node: %p", uuidString, node)
		node, err = t.cluster.AddNode(&uuid, uuidString, report)
		if err != nil {
			return
		}
	} else {
		err = t.cluster.UpdateNode(node, &uuid, uuidString, report)
		if err != nil {
			return
		}
	}

	ctx := identity.WorkerContext{
		Node:        node,
		Cluster:     t.cluster,
		ShadowFile:  t.shadowFile,
		RootAccount: t.rootAccount,
	}
	tarball.SendTarball(identity.WorkerManifest, &ctx, w, "worker-credentials.tar.gz")
}

func (t *TeamsterAPI) GenClientCert(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	user := query.Get("user")
	groups := query["group"]

	if user == "" || len(groups) == 0 {
		http.Error(w, "request should include 'user' and 'group' arguments", http.StatusBadRequest)
		return
	}

	c, p, err := t.cluster.AddUser(user, groups)
	if err != nil {
		panic(errors.Wrap(err, "failed to create user credentials"))
	}

	ip, err := getAPIServerIP(t.cluster.Vars["CONTROLLER_PRIVATE_IF"])
	if err != nil {
		panic(errors.Wrap(err, "failed to obtain controller private IP"))
	}

	ctx := identity.ClientContext{
		Cert:      c,
		Key:       p,
		Bundle:    t.cluster.CA_Bundle,
		InstallID: t.cluster.InstallID,
		User:      user,
	}
	if ip != "" {
		ctx.ServerURL = fmt.Sprintf("https://%s:%s", ip, t.cluster.Vars["OPEROS_KUBE_API_SECURE_PORT"])
	}

	tarball.SendTarball(identity.ClientManifest, ctx, w, "operos-credentials.tar.gz")
}

func (t *TeamsterAPI) ListNodes(ctx context.Context, req *Empty) (*ListNodesResponse, error) {
	respNodes := make([]*NodeSummary, len(t.cluster.Nodes))
	idx := 0
	for uuid := range t.cluster.Nodes {
		respNodes[idx] = &NodeSummary{Uuid: uuid}
		idx++
	}
	return &ListNodesResponse{Nodes: respNodes}, nil
}

func (t *TeamsterAPI) GetNodeHardware(ctx context.Context, req *GetNodeHardwareRequest) (*GetNodeHardwareResponse, error) {
	if node, ok := t.cluster.Nodes[req.Uuid]; ok {
		//XXX: disregard error, bad form
		hinfo, _ := json.Marshal(node.LatestReport.System)
		return &GetNodeHardwareResponse{HardwareInfo: string(hinfo)}, nil
	}
	return nil, grpc.Errorf(codes.NotFound, "node not found")
}

func getAPIServerIP(ifname string) (string, error) {
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		return "", err
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return "", err
	}

	var ipAddr net.IP
	for _, addr := range addrs {
		if ip, ok := addr.(*net.IPNet); ok && ip.IP.To4() != nil {
			ipAddr = ip.IP
		}
	}

	if ipAddr == nil {
		return "", nil
	}

	return ipAddr.String(), nil
}

func (t *TeamsterAPI) GetCACertExpiry(context.Context, *Empty) (*GetCACertExpiryResponse, error) {
	cert, err := t.cluster.GetCACert()
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse CA certificate")
	}

	return &GetCACertExpiryResponse{
		ExpiryUnix: cert.NotAfter.Unix(),
	}, nil
}

const saltChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789./"

func (t *TeamsterAPI) SetRootPassword(ctx context.Context, req *SetRootPasswordRequest) (*Empty, error) {
	// This method manually salts and encrypts the password using crypt(3)
	// Another option is to shell out to something like chpasswd. The problem
	// with that approach is that all of the shadow tools (including chpasswd)
	// hardcode the path "/etc/shadow". When updating the passwords, they write
	// to a temporary file, then copy it over /etc/shadow, like we do here.
	//
	// When running in a Docker container, if we want to be able to change the
	// host's files, we have to mount the entire /etc directory - if only
	// /etc/shadow was mounted, then replacing that file would be an attempt
	// to replace the mountpoint, which will return an error.
	//
	// The host's /etc should probably be mounted to something other than /etc
	// in the container to prevent any misconfiguration issues. This means that
	// changing the password must write /etc-host/shadow, not /etc/shadow,
	// rendering chpassword useless. Thank you for reading my essay.

	if strings.TrimSpace(req.Password) == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "password must not be blank")
	}

	// Generate salt first
	salt := make([]byte, 16)
	for i := 0; i < len(salt); i++ {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(saltChars)-1)))
		if err != nil {
			return nil, errors.Wrap(err, "failed generating salt")
		}

		salt[i] = saltChars[idx.Int64()]
	}

	// crypt(3) the password
	hashedPass, err := crypt.Crypt(req.Password, fmt.Sprintf("$6$%s$", salt))
	if err != nil {
		return nil, errors.Wrap(err, "failed to hash password")
	}

	// Write a new version of the shadow file
	shadow, err := os.Open(t.shadowFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open shadow file")
	}
	defer shadow.Close()

	tmpFile, err := ioutil.TempFile(path.Dir(t.shadowFile), "shadow")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create temporary file")
	}

	if err := tmpFile.Chmod(0640); err != nil {
		return nil, errors.Wrap(err, "failed to set temporary file permissions")
	}

	// Deliberately ignore errors here, this is just in case there's an
	// early return part way through this method.
	defer func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}()

	scanner := bufio.NewScanner(shadow)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, errors.Wrap(err, "failed to read shadow file")
		}

		line := strings.TrimSpace(scanner.Text())
		parts := strings.Split(line, ":")

		if parts[0] == t.rootAccount {
			parts[1] = hashedPass
			parts[2] = fmt.Sprintf("%.0f", time.Since(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)).Hours()/24)
		}

		fmt.Fprintln(tmpFile, strings.Join(parts, ":"))
	}

	// Replace the existing shadow file
	if err := os.Rename(tmpFile.Name(), t.shadowFile); err != nil {
		return nil, errors.Wrap(err, "failed to rename temporary shadow file")
	}

	return &Empty{}, nil
}
