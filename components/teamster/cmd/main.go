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
	"flag"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/gorilla/handlers"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/paxautoma/operos/components/teamster/pkg/cluster"
	"github.com/paxautoma/operos/components/teamster/pkg/teamster"
)

func main() {
	logger := logrus.StandardLogger()
	log.SetOutput(logger.Writer())

	installID := flag.String("install-id", "", "the Install ID of the Operos Cluster")
	listenAddr := flag.String("listen-addr", ":2680", "the address:port that teamster should bind to for HTTP/1")
	listenGrpc := flag.String("listen-grpc", ":2681", "the address:port that teamster should bind to for gRPC")
	etcdCluster := flag.String("etcd-cluster", "localhost:2379", "the hostname:port of the etcd cluster to connect to")
	shadowFile := flag.String("shadow-file", "/etc/shadow", "name of the shadow file to use to obtain root password")
	rootAccount := flag.String("root", "root", "user name of the user whose password hash will be sent to worker nodes")

	flag.Parse()

	if *installID == "" {
		log.Fatal("error: You must specify the Operos Install ID")
	}

	log.Printf("cluster: %s", *installID)
	log.Printf("etcd cluster: %s", *etcdCluster)

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{*etcdCluster},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		log.Fatalf("error: Unable to connect to Etcd cluster: %s", err)
	}

	oc, err := cluster.InstantiateCluster(client, 5*time.Second, *installID)

	if err != nil {
		log.Fatalf("error: Unable to instantiate operos cluster: %s", err)
	}

	api := teamster.NewTeamsterAPI(oc, *shadowFile, *rootAccount)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logger)),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)
	api.RegisterGRPCService(grpcServer)

	lis, err := net.Listen("tcp", *listenGrpc)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		grpcServer.Serve(lis)
	}()

	var handler http.Handler
	handler = handlers.RecoveryHandler(handlers.RecoveryLogger(logger))(api.GetHttpHandler())
	handler = handlers.LoggingHandler(logger.Writer(), handler)
	log.Fatal(http.ListenAndServe(*listenAddr, handler))
}
