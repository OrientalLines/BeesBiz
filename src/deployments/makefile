# Variables
NAMESPACE := beesbiz-data
RELEASE_NAME := beesbiz
CHART_FILE := BeesBizData-0.1.0.tgz
KUBECTL := kubectl
HELM := helm
MINIKUBE := minikube

.PHONY: start uninstall template lint package upgrade list get_all port-forward stop-port-forward recreate pause resume dapr-init dapr-components dapr-status dapr-dashboard dapr-install pubsub postgresql dapr-install

uninstall:
	$(HELM) uninstall $(RELEASE_NAME) --namespace $(NAMESPACE)

list:
	$(HELM) list -n $(NAMESPACE)

get_all:
	$(KUBECTL) get all -n $(NAMESPACE)

upgrade:
	$(HELM) upgrade --install $(RELEASE_NAME) ./$(CHART_FILE) --namespace $(NAMESPACE)

package:
	$(HELM) package .

lint:
	$(HELM) lint .

template:
	$(HELM) template ./$(CHART_FILE)

recreate: package upgrade

dapr-status:
	dapr status -k

dapr-dashboard:
	dapr dashboard -k -p 8080

pubsub:
	$(KUBECTL) apply -f ./components/pubsub.yaml -n $(NAMESPACE);
	$(KUBECTL) rollout restart deployment -n dapr-system

postgresql:
	$(KUBECTL) apply -f ./components/postgresql.yaml -n $(NAMESPACE);
	$(KUBECTL) rollout restart deployment -n dapr-system

start:
	$(MINIKUBE) start --cpus 6 --memory 10244 --disk-size 20g
	$(MINIKUBE) addons enable storage-provisioner
	$(MINIKUBE) addons enable default-storageclass
	$(KUBECTL) create namespace $(NAMESPACE)
	$(KUBECTL) apply -f persistent-volumes.yaml -n $(NAMESPACE) --wait
	$(HELM) repo add cnpg https://cloudnative-pg.github.io/charts
	$(HELM) repo update
	$(HELM) install cnpg cnpg/cloudnative-pg --namespace cnpg-system --create-namespace
	$(KUBECTL) create secret generic postgresql-superuser -n $(NAMESPACE) \
     --from-literal=username=postgres \
     --from-literal=password=postgres

dapr-install:
	$(HELM) repo add dapr https://dapr.github.io/helm-charts/
	$(HELM) repo update
	$(HELM) install dapr dapr/dapr --namespace dapr-system --create-namespace

dapr-init:
	dapr init -k

pause:
	$(MINIKUBE) pause

resume:
	$(MINIKUBE) start

port-forward:
	@echo "Starting port forwarding for PostgreSQL, TiKV, and RabbitMQ..."
	@echo "Press Ctrl+C to stop port forwarding."
	@$(KUBECTL) port-forward svc/postgresql-rw 5432:5432 -n $(NAMESPACE) & \
	$(KUBECTL) port-forward svc/tikv 20160:20160 -n $(NAMESPACE) & \
	$(KUBECTL) port-forward svc/rabbitmq 5672:5672 15672:15672 -n $(NAMESPACE) & \
	$(KUBECTL) port-forward deployment/postgresql 3500:3500 50001:50001 -n $(NAMESPACE)
	@wait

stop-port-forward:
	@echo "Stopping all port-forwarding processes..."
	@pkill -f "$(KUBECTL) port-forward"
	@echo "Port forwarding stopped."
