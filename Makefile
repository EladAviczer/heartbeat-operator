.PHONY: load-test
load-test:
	@echo "Fetching Kubernetes context for load test..."
	@kubectl create clusterrolebinding default-admin --clusterrole=cluster-admin --serviceaccount=default:default >/dev/null 2>&1 || true ; \
	API_URL=$$(kubectl config view --minify -o jsonpath='{.clusters[0].cluster.server}') ; \
	TOKEN=$$(kubectl create token default) ; \
	echo "Running K6 load test against API Server: $$API_URL" ; \
	KUBERNETES_API_URL=$$API_URL KUBERNETES_TOKEN=$$TOKEN k6 run load-testing/scenarios/control-plane.js

.PHONY: clean-load-test
clean-load-test:
	@./load-testing/cleanup.sh
