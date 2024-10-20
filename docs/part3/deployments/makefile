# Variables
NAMESPACE := beesbiz-data
RELEASE_NAME := beesbiz
CHART_FILE := BeesBizData-0.1.0.tgz
KUBECTL := kubectl
HELM := helm
MINIKUBE := minikube

.PHONY: start uninstall template lint package upgrade list get_all port-forward stop-port-forward recreate

uninstall:
	$(HELM) uninstall $(RELEASE_NAME) --namespace $(NAMESPACE)

list:
	$(HELM) list -n $(NAMESPACE)

get_all:
	$(KUBECTL) get all -n $(NAMESPACE)

upgrade: package
	$(HELM) upgrade --install $(RELEASE_NAME) ./$(CHART_FILE) --namespace $(NAMESPACE)

package:
	$(HELM) package .

lint:
	$(HELM) lint .

template:
	$(HELM) template ./$(CHART_FILE)

recreate: package upgrade

port-forward:
	@echo "Starting port forwarding for PostgreSQL, TiKV, and RabbitMQ..."
	@echo "Press Ctrl+C to stop port forwarding."
	@$(KUBECTL) port-forward svc/postgresql-rw 5432:5432 -n $(NAMESPACE) & \
	$(KUBECTL) port-forward svc/tikv 20160:20160 -n $(NAMESPACE) & \
	$(KUBECTL) port-forward svc/rabbitmq 5672:5672 15672:15672 -n $(NAMESPACE)
	@wait

stop-port-forward:
	@echo "Stopping all port-forwarding processes..."
	@pkill -f "$(KUBECTL) port-forward"
	@echo "Port forwarding stopped."

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
