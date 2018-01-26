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

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"github.com/paxautoma/operos/components/common/gatekeeper"
)

var gatekeeperAddress = flag.String("gatekeeper", "gatekeeper.paxautoma.com:57345", "address of the Gatekeeper server (host:port)")
var noGatekeeperTLS = flag.Bool("no-gatekeeper-tls", false, "do not use TLS when contacting Gatekeeper")
var installID = flag.String("install", "", "install ID")

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "Please specify file to upload")
		os.Exit(1)
	}

	if *installID == "" {
		fmt.Fprintln(os.Stderr, "Please specify the install ID")
		os.Exit(1)
	}

	fileName := flag.Arg(0)

	opts := []grpc.DialOption{
		grpc.WithTimeout(10 * time.Second),
	}

	if *noGatekeeperTLS {
		opts = append(opts, grpc.WithInsecure())
	} else {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	}

	conn, err := grpc.Dial(*gatekeeperAddress, opts...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error dialing (%s)\n", grpc.ErrorDesc(err))
		os.Exit(2)
	}
	defer conn.Close()

	client := gatekeeper.NewGatekeeperClient(conn)

	auth, err := client.AuthorizeDiagUpload(context.Background(), &gatekeeper.AuthorizeDiagUploadReq{
		InstallId: *installID,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to request authorization from Pax Automa:\n%s\n", err.Error())
		os.Exit(2)
	}

	fmt.Printf(strings.Join([]string{
		"curl",
		"-f",
		"-o /dev/null",
		`-F "acl=bucket-owner-full-control"`,
		fmt.Sprintf(`-F "policy=%s"`, auth.GetPolicy()),
		fmt.Sprintf(`-F "key=%s/%s"`, *installID, path.Base(fileName)),
		fmt.Sprintf(`-F "x-amz-date=%sT000000Z"`, auth.GetDate()),
		fmt.Sprintf(`-F "x-amz-credential=%s"`, auth.GetCredential()),
		`-F "x-amz-algorithm=AWS4-HMAC-SHA256"`,
		fmt.Sprintf(`-F "x-amz-signature=%s"`, auth.GetSignature()),
		fmt.Sprintf(`-F "file=@%s"`, fileName),
		auth.GetUrl(),
	}, " "))
}
