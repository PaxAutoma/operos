TEAMSTER_HTTP=192.168.33.10:2680
TEAMSTER=192.168.33.10:2681
KUBE=localhost:8001
PROMETHEUS=localhost:9090

# Assumes you have:
# - a local Operos controller at 192.168.33.10
# - running "kubectl kube-proxy", listening on localhost:8001
# - running "kubectl port-forward <prometheus-pod> 9090"
# use KUBE, TEAMSTER, PROMETHEUS to override above 
.PHONY: dev
dev:
	GRPC_GO_LOG_SEVERITY_LEVEL=info \
	GRPC_GO_LOG_VERBOSITY_LEVEL=10 \
	go run -v cmd/main.go \
		--kubeconfig=../../operos/kubeconfig \
		--kube-url="http://$(KUBE)" \
		--teamster $(TEAMSTER) \
		--teamster-http $(TEAMSTER_HTTP) \
		--prometheus-url "http://$(PROMETHEUS)/api/v1" \
		--clientdir "../client/src" \
		--debug-addr 127.0.0.1:2782
