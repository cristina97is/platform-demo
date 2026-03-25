CLUSTER_NAME ?= platform-demo
NAMESPACE ?= event-service
IMAGE ?= event-service:local
INGRESS_HOST ?= event-service.127.0.0.1.sslip.io

.PHONY: build load deploy restart status logs test-api load-test

build:
	cd apps/event-service && docker build -t $(IMAGE) .

load:
	kind load docker-image $(IMAGE) --name $(CLUSTER_NAME)

deploy:
	kubectl apply -f apps/event-service/k8s/namespace.yaml
	kubectl apply -f apps/event-service/k8s/secret.yaml
	kubectl apply -f apps/event-service/k8s/configmap.yaml
	kubectl apply -f apps/event-service/k8s/db-pvc.yaml
	kubectl apply -f apps/event-service/k8s/db-init-configmap.yaml
	kubectl apply -f apps/event-service/k8s/db-deployment.yaml
	kubectl apply -f apps/event-service/k8s/db-service.yaml
	kubectl apply -f apps/event-service/k8s/api-deployment.yaml
	kubectl apply -f apps/event-service/k8s/api-service.yaml
	kubectl apply -f apps/event-service/k8s/ingress.yaml
	kubectl apply -f apps/event-service/k8s/api-hpa.yaml
	kubectl apply -f apps/event-service/k8s/api-pdb.yaml

restart:
	kubectl rollout restart deployment/event-service-api -n $(NAMESPACE)
	kubectl rollout status deployment/event-service-api -n $(NAMESPACE)

status:
	kubectl get all -n $(NAMESPACE)
	kubectl get pvc -n $(NAMESPACE)
	kubectl get hpa -n $(NAMESPACE)
	kubectl get pdb -n $(NAMESPACE)
	kubectl get ingress -n $(NAMESPACE)

logs:
	kubectl logs -n $(NAMESPACE) deploy/event-service-api --tail=100

test-api:
	curl -i http://$(INGRESS_HOST)/
	echo
	curl -i http://$(INGRESS_HOST)/healthz
	echo
	curl -i http://$(INGRESS_HOST)/readyz
	echo
	curl -i http://$(INGRESS_HOST)/events
	echo

load-test:
	cd load-test && node test.js
