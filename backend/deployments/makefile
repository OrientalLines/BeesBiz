# Variables
NAMESPACE := beesbiz-data # Default namespace
TIKV_NAMESPACE := beesbiz-tikv
MONITORING_NAMESPACE := beesbiz-monitoring
CLOUDNATIVEPG_NAMESPACE := beesbiz-data
RABBITMQ_NAMESPACE := beesbiz-rabbitmq
RUNTIME_NAMESPACE := beesbiz-runtime
TIKV_RUNTIME := beesbiz-server

RELEASE_NAME := beesbiz
CHART_FILE := BeesBizData-0.1.0.tgz
KUBECTL := kubectl
HELM := helm
MINIKUBE := minikube

.PHONY: start uninstall template lint package upgrade list get_all port-forward stop-port-forward recreate pause resume pubsub postgresql uninstall uninstall-all upgrade upgrade-all recreate recreate-all install-tikv install-monitoring install-postgresql install-rabbitmq uninstall-tikv uninstall-monitoring uninstall-postgresql uninstall-rabbitmq docker-build update-deps setup-monitoring setup-servicemonitors patch-postgresql-monitoring verify-monitoring monitoring-dashboard prometheus-dashboard

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

pubsub:
	$(KUBECTL) apply -f ./components/pubsub.yaml -n $(NAMESPACE);

docker-build:
	docker build -t $(TIKV_RUNTIME):latest .

load:
	$(MINIKUBE) image load $(TIKV_RUNTIME):latest

start:
	$(MINIKUBE) start --cpus 8 --memory 10244 --disk-size 20g
	$(MINIKUBE) addons enable storage-provisioner
	$(MINIKUBE) addons enable default-storageclass
	$(KUBECTL) create namespace $(NAMESPACE)
	$(KUBECTL) create namespace $(RABBITMQ_NAMESPACE)
	$(KUBECTL) create namespace $(MONITORING_NAMESPACE)
	$(KUBECTL) create namespace $(RUNTIME_NAMESPACE)
	# $(KUBECTL) apply -f persistent-volumes.yaml -n $(NAMESPACE) --wait
	$(HELM) repo add cnpg https://cloudnative-pg.github.io/charts
	$(HELM) repo add pingcap https://charts.pingcap.org/
	$(HELM) repo update
	$(KUBECTL) apply --server-side=true -f https://raw.githubusercontent.com/pingcap/tidb-operator/v1.6.0/manifests/crd.yaml
	$(KUBECTL) apply --server-side=true -f https://github.com/rabbitmq/cluster-operator/releases/latest/download/cluster-operator.yml
	$(HELM) install cnpg cnpg/cloudnative-pg --namespace $(CLOUDNATIVEPG_NAMESPACE) --create-namespace
	$(HELM) install tidb-operator pingcap/tidb-operator --namespace $(TIKV_NAMESPACE) --create-namespace -f ./values-tidb-operator.yaml
	# $(HELM) install tidb-operator pingcap/tidb-operator --namespace $(TIKV_NAMESPACE) --create-namespace
	$(KUBECTL) create secret generic postgresql-superuser -n $(NAMESPACE) \
     --from-literal=username=postgres \
     --from-literal=password=postgres
	$(MAKE) install-monitoring
	$(MAKE) setup-servicemonitors
	$(MAKE) patch-postgresql-monitoring

pause:
	$(MINIKUBE) pause

resume:
	$(MINIKUBE) start

install-monitoring:
	$(HELM) repo add prometheus-community https://prometheus-community.github.io/helm-charts
	$(HELM) repo add grafana https://grafana.github.io/helm-charts
	$(HELM) repo update
	$(HELM) dependency update
	$(HELM) upgrade --install monitoring prometheus-community/kube-prometheus-stack \
		--namespace $(MONITORING_NAMESPACE) \
		--create-namespace \
		-f ./values.yaml
	@echo "Waiting for Prometheus operator to be ready..."
	sleep 30
	$(MAKE) setup-servicemonitors
	$(MAKE) patch-postgresql-monitoring

uninstall-monitoring:
	$(HELM) uninstall monitoring --namespace $(MONITORING_NAMESPACE)

update-deps:
	$(HELM) repo add prometheus-community https://prometheus-community.github.io/helm-charts
	$(HELM) repo add grafana https://grafana.github.io/helm-charts
	$(HELM) repo add cnpg https://cloudnative-pg.github.io/charts
	$(HELM) repo add pingcap https://charts.pingcap.org/
	$(HELM) repo update
	$(HELM) dependency update

setup-users: setup-rabbitmq-user setup-postgresql-user

setup-postgresql-user:
	$(KUBECTL) exec -it -n $(CLOUDNATIVEPG_NAMESPACE) $$($(KUBECTL) get pods -n $(CLOUDNATIVEPG_NAMESPACE) | grep postgresql | awk '{print $$1}') -- psql -U postgres -c "DO \$$\$$ BEGIN IF NOT EXISTS (SELECT FROM pg_user WHERE usename = 'user') THEN CREATE USER \"user\" WITH PASSWORD 'postgres'; GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO \"user\"; ALTER USER \"user\" CREATEDB; GRANT ALL PRIVILEGES ON DATABASE postgres TO \"user\"; END IF; END \$$\$$;"

setup-rabbitmq-user:
	$(KUBECTL) exec -it -n $(RABBITMQ_NAMESPACE) $$($(KUBECTL) get pods -n $(RABBITMQ_NAMESPACE) | grep rabbitmq-cluster | awk '{print $$1}') -- /bin/bash -c "rabbitmqctl delete_user guest || true && rabbitmqctl add_user guest guest && rabbitmqctl set_user_tags guest administrator && rabbitmqctl set_permissions -p / guest '.*' '.*' '.*'"

setup-monitoring: setup-rabbit-monitoring setup-postgresql-monitoring setup-grafana

setup-rabbit-monitoring:
	sh rabbitmq-monitoring-setup.sh

setup-postgresql-monitoring:
	sh postgresql-monitoring-setup.sh

setup-grafana:
	$(KUBECTL) exec -it -n $(MONITORING_NAMESPACE) $$($(KUBECTL) get pods -n $(MONITORING_NAMESPACE) | grep grafana | awk '{print $$1}') -- grafana cli admin reset-admin-password admin

setup-servicemonitors:
	@echo "Applying monitoring configurations..."
	$(HELM) upgrade --install $(RELEASE_NAME) ./$(CHART_FILE) --namespace $(NAMESPACE)

patch-postgresql-monitoring:
	kubectl patch cluster postgresql -n $(CLOUDNATIVEPG_NAMESPACE) --type merge -p '{"spec":{"monitoring":{"enablePodMonitor":true,"customQueriesConfigMap":[{"name":"postgres-metrics-config","key":"queries.yaml"}]}}}'

verify-monitoring:
	@echo "Verifying ServiceMonitors..."
	$(KUBECTL) get servicemonitor -n $(MONITORING_NAMESPACE)
	@echo "\nVerifying PostgreSQL metrics..."
	$(KUBECTL) port-forward -n $(CLOUDNATIVEPG_NAMESPACE) postgresql-1 9187:9187 & \
	sleep 5 && curl http://localhost:9187/metrics && kill %%
	@echo "\nVerifying RabbitMQ metrics..."
	$(KUBECTL) port-forward svc/rabbitmq-cluster -n $(RABBITMQ_NAMESPACE) 15692:15692 & \
	sleep 5 && curl http://localhost:15692/metrics && kill %%

monitoring-dashboard:
	$(KUBECTL) port-forward svc/monitoring-grafana -n $(MONITORING_NAMESPACE) 3000:80

prometheus-dashboard:
	$(KUBECTL) port-forward svc/monitoring-kube-prometheus-prometheus -n $(MONITORING_NAMESPACE) 9090:9090
