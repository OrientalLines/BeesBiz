# Variables
NAMESPACE := beesbiz-data # Default namespace
TIKV_NAMESPACE := beesbiz-tikv
MONITORING_NAMESPACE := beesbiz-monitoring
CLOUDNATIVEPG_NAMESPACE := beesbiz-data
RABBITMQ_NAMESPACE := beesbiz-rabbitmq

RELEASE_NAME := beesbiz
CHART_FILE := BeesBizData-0.1.0.tgz
KUBECTL := kubectl
HELM := helm
MINIKUBE := minikube

.PHONY: start uninstall template lint package upgrade list get_all port-forward stop-port-forward recreate pause resume dapr-init dapr-components dapr-status dapr-dashboard dapr-install pubsub postgresql dapr-install uninstall uninstall-all upgrade upgrade-all recreate recreate-all install-tikv install-monitoring install-postgresql install-rabbitmq uninstall-tikv uninstall-monitoring uninstall-postgresql uninstall-rabbitmq

uninstall:
	$(HELM) uninstall $(RELEASE_NAME) --namespace $(NAMESPACE)

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

start:
	$(MINIKUBE) start --cpus 8 --memory 10244 --disk-size 20g
	$(MINIKUBE) addons enable storage-provisioner
	$(MINIKUBE) addons enable default-storageclass
	$(KUBECTL) create namespace $(NAMESPACE)
	$(KUBECTL) create namespace $(RABBITMQ_NAMESPACE)
	$(KUBECTL) create namespace $(MONITORING_NAMESPACE)
	# $(KUBECTL) apply -f persistent-volumes.yaml -n $(NAMESPACE) --wait
	$(HELM) repo add cnpg https://cloudnative-pg.github.io/charts
	$(HELM) repo add pingcap https://charts.pingcap.org/
	$(HELM) repo update
	$(KUBECTL) apply --server-side=true -f https://raw.githubusercontent.com/pingcap/tidb-operator/v1.6.0/manifests/crd.yaml
	$(HELM) install cnpg cnpg/cloudnative-pg --namespace $(CLOUDNATIVEPG_NAMESPACE) --create-namespace
	# $(HELM) install tidb-operator pingcap/tidb-operator --namespace $(TIKV_NAMESPACE) --create-namespace -f ./values-tidb-operator.yaml
	$(HELM) install tidb-operator pingcap/tidb-operator --namespace $(TIKV_NAMESPACE) --create-namespace
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
