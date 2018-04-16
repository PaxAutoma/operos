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

package installer

import (
	"fmt"

	"github.com/cloudflare/cfssl/config"
	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/helpers"
	"github.com/cloudflare/cfssl/initca"
	"github.com/cloudflare/cfssl/signer"
	"github.com/cloudflare/cfssl/signer/local"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/paxautoma/operos/components/common"
)

func CreateControllerCerts(ctx *InstallerContext) error {
	caCertBytes, caKeyBytes, err := CreateControllerCA(ctx)
	if err != nil {
		return err
	}

	ctx.ControllerCert = string(caCertBytes)
	ctx.ControllerKey = string(caKeyBytes)

	serverCert, serverKey, err := GenerateAPIServerCert(ctx, caCertBytes, caKeyBytes)
	if err != nil {
		log.Fatalf("Failed to sign API server cert: %+v", err)
	}
	ctx.ServerCert = string(serverCert)
	ctx.ServerKey = string(serverKey)

	return nil
}

func CreateControllerCA(ctx *InstallerContext) (certBytes, keyBytes []byte, errOut error) {
	req := &csr.CertificateRequest{
		KeyRequest: &csr.BasicKeyRequest{
			A: "rsa",
			S: 2048,
		},
		CN: fmt.Sprintf("%s (Controller CA)", ctx.Responses.OrgInfo.Cluster),
		Names: []csr.Name{
			{
				C:  ctx.Responses.OrgInfo.Country,
				ST: ctx.Responses.OrgInfo.Province,
				L:  ctx.Responses.OrgInfo.City,
				O:  ctx.Responses.OrgInfo.Organization,
				OU: ctx.Responses.OrgInfo.Department,
			},
		},
	}

	csrBytes, _, keyBytes, err := initca.New(req)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not create controller CA cert")
	}

	return csrBytes, keyBytes, nil
}

func GenerateAPIServerCert(ctx *InstallerContext, caCert []byte, caKey []byte) (certBytes, keyBytes []byte, errOut error) {
	csrBytes, keyBytes, err := createAPIServerCSR(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not create API server CSR")
	}

	s, err := createAPISigner(caCert, caKey)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not create controller CA signer")
	}
	req := signer.SignRequest{
		Request: string(csrBytes),
	}

	certBytes, err = s.Sign(req)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not sign API server CSR")
	}

	return certBytes, keyBytes, nil
}

func createAPISigner(caCertBytes, caKeyBytes []byte) (signer.Signer, error) {
	policy := &config.Signing{
		Default: &config.SigningProfile{
			Expiry: helpers.OneYear,
			Usage: []string{
				"digital signature",
				"key encipherment",
				"server auth",
			},
		},
	}

	caCert, err := helpers.ParseCertificatePEM(caCertBytes)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse controller CA cert")
	}

	caKey, err := helpers.ParsePrivateKeyPEM(caKeyBytes)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse controller private key")
	}

	result, err := local.NewSigner(caKey, caCert, signer.DefaultSigAlgo(caKey), policy)
	if err != nil {
		return nil, errors.Wrap(err, "could not create controller CA signer")
	}

	return result, nil
}

func createAPIServerCSR(ctx *InstallerContext) (csrBytes, keyBytes []byte, errOut error) {
	hosts := []string{
		ctx.Responses.ControllerIP,
		ctx.Responses.KubeAPIServiceIP,
		"127.0.0.1",
		"localhost",
		"kubernetes.default.svc",
	}

	if ctx.Responses.PublicHostname != "" && !common.ArrayContains(hosts, ctx.Responses.PublicHostname) {
		hosts = append(hosts, ctx.Responses.PublicHostname)
	}

	req := &csr.CertificateRequest{
		KeyRequest: &csr.BasicKeyRequest{
			A: "rsa",
			S: 2048,
		},
		Hosts: hosts,
		CN:    fmt.Sprintf("%s (Controller Server)", ctx.Responses.OrgInfo.Cluster),
		Names: []csr.Name{
			{
				C:  ctx.Responses.OrgInfo.Country,
				ST: ctx.Responses.OrgInfo.Province,
				L:  ctx.Responses.OrgInfo.City,
				O:  ctx.Responses.OrgInfo.Organization,
				OU: ctx.Responses.OrgInfo.Department,
			},
		},
	}

	csrBytes, keyBytes, err := csr.ParseRequest(req)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not parse controller server CSR")
	}

	return
}
