.PHONY: minikube uninstall template lint package upgrade list get_all

uninstall:
	helm uninstall beesbiz --namespace beesbiz-data

list:
	helm list -n beesbiz-data

get_all:
	kubectl get all -n beesbiz-data

upgrade:
	helm upgrade --install beesbiz ./BeesBizData-0.1.0.tgz --namespace beesbiz-data

package:
	helm package .

lint:
	helm lint .

template:
	helm template ./BeesBizData-0.1.0.tgz

minikube:
	minikube stop
	minikube delete
	minikube start --cpus 4 --memory 10244 --disk-size 20g
	minikube addons enable storage-provisioner
	minikube addons enable default-storageclass
