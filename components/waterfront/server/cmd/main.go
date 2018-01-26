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
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/http/pprof"
	"net/url"
	"path"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/paxautoma/operos/components/waterfront/server/pkg/kube"
	"github.com/paxautoma/operos/components/waterfront/server/pkg/teamster"
	"github.com/paxautoma/operos/components/waterfront/server/pkg/waterfront"
)

func main() {
	logger := logrus.StandardLogger()
	log.SetOutput(logger.Writer())

	listenAddr := flag.String("listen-addr", ":2780", "address:port to listen for HTTP/JSON")
	listenGrpc := flag.String("listen-grpc", ":2781", "address:port to listen for gRPC")
	teamsterHTTPAddr := flag.String("teamster-http", "localhost:2680", "teamster HTTP endpoint")
	teamsterAddr := flag.String("teamster", "localhost:2681", "teamster endpoint")
	kubeURL := flag.String("kube-url", "http://localhost:8080", "kubernetes API server endpoint")
	kubeConfig := flag.String("kubeconfig", "", "kubeconfig file; if not set, service account is used")
	clientDir := flag.String("clientdir", "client", "directory containing the client files")
	dashboardURL := flag.String("dashboard-url", "http://kubernetes-dashboard.kube-system", "URL of the kube-dashboard for proxying")
	promURL := flag.String("prometheus-url", "http://prometheus.operos:9090/api/v1", "URL of Prometheus API for proxying")
	sessionKey := flag.String("session-key", "this is a totally secret key", "cookie session key")
	debugAddr := flag.String("debug-addr", "", "enable debug server on this address")

	flag.Parse()

	kubeDashboardURL, err := url.Parse(*dashboardURL)
	if err != nil {
		log.Fatalf("invalid dashboard URL")
	}

	prometheusURL, err := url.Parse(*promURL)
	if err != nil {
		log.Fatalf("invalid Prometheus URL")
	}

	teamsterClient, err := teamster.NewTeamsterClientFromAddr(*teamsterAddr)
	if err != nil {
		log.Fatalf("failed to instantiate teamster client: %v", err)
	}

	kubeClient, err := kube.NewKubeClient(*kubeURL, *kubeConfig)
	if err != nil {
		log.Fatalf("failed to create kube client: %v", err)
	}

	waterfrontAPI := waterfront.NewWaterfrontAPI(teamsterClient, kubeClient)

	lis, err := net.Listen("tcp", *listenGrpc)
	if err != nil {
		log.Fatalf("failed to listen on gRPC port: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logger)),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)
	waterfront.RegisterWaterfrontServer(grpcServer, waterfrontAPI)

	grpcMux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBackoffMaxDelay(10 * time.Second),
	}
	if err := waterfront.RegisterWaterfrontHandlerFromEndpoint(context.Background(), grpcMux, *listenGrpc, opts); err != nil {
		log.Fatalf("error registering gRPC handler: %s", err.Error())
	}

	go func() {
		log.Printf("Listening for gRPC on %s", *listenGrpc)
		grpcServer.Serve(lis)
	}()

	auth := waterfront.AuthSessionMiddleware(
		waterfront.SessionStore(sessions.NewCookieStore([]byte(*sessionKey))))
	cors := handlers.CORS(
		handlers.AllowCredentials(),
		handlers.AllowedOrigins([]string{"http://localhost:10000"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)

	prometheusProxy := httputil.NewSingleHostReverseProxy(prometheusURL)
	prometheusProxy.ModifyResponse = func(res *http.Response) error {
		// Remove the Prometheus CORS stuff since we're gonna apply our own
		res.Header.Del("Access-Control-Allow-Origin")
		res.Header.Del("Access-Control-Allow-Credentials")
		return nil
	}

	apiRouter := mux.NewRouter()
	// Our single HTTP URL for generating client cert tarballs
	apiRouter.Path("/api/v1/clientcert").Handler(waterfront.MakeGenClientCertHandler(*teamsterHTTPAddr))
	// Proxy to Prometheus
	apiRouter.PathPrefix("/api/v1/metrics/").Handler(http.StripPrefix("/api/v1/metrics/", cors(prometheusProxy)))
	// grpc-proxy API
	apiRouter.PathPrefix("/api/").Handler(cors(http.StripPrefix("/api", grpcMux)))

	mainRouter := mux.NewRouter()
	// Auth endpoints
	mainRouter.Handle("/api/v1/login", cors(auth.GetLoginHandler()))
	mainRouter.Handle("/api/v1/logout", cors(auth.GetLogoutHandler()))
	// API
	mainRouter.PathPrefix("/api/").Handler(auth.GetAuthHandler(apiRouter))

	// Proxy to the kube-dashboard
	mainRouter.PathPrefix("/kube-dashboard/").Handler(auth.GetAuthHandler(http.StripPrefix("/kube-dashboard", httputil.NewSingleHostReverseProxy(kubeDashboardURL))))

	// Static files
	mainRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static", http.FileServer(http.Dir(path.Join(*clientDir, "static")))))
	mainRouter.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		return r.URL.Path == "/app.js" || r.URL.Path == "/app.js.map"
	}).Handler(http.FileServer(http.Dir(*clientDir)))

	// The UI SPA uses dynamic routes so everything else just forwards to index.html
	mainRouter.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join(*clientDir, "index.html"))
	})

	handler := handlers.LoggingHandler(logger.Writer(), handlers.RecoveryHandler()(mainRouter))

	if *debugAddr != "" {
		go func() {
			debugH := http.NewServeMux()
			debugH.HandleFunc("/debug/pprof/", pprof.Index)
			debugH.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
			debugH.HandleFunc("/debug/pprof/profile", pprof.Profile)
			debugH.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
			debugH.HandleFunc("/debug/pprof/trace", pprof.Trace)

			log.Printf("Debug server listening on on %s", *debugAddr)
			http.ListenAndServe(*debugAddr, debugH)
		}()
	}

	log.Printf("Listening for HTTP on %s", *listenAddr)
	http.ListenAndServe(*listenAddr, handler)
}
