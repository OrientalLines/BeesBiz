RUNTIME_NAMESPACE := beesbiz-runtime
SERVER_RUNTIME := beesbiz-server

MINIKUBE := minikube

.PHONY: all build load

build:
	docker build -t $(SERVER_RUNTIME):latest .

load:
	$(MINIKUBE) image load $(SERVER_RUNTIME):latest
