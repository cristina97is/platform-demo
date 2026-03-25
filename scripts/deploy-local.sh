#!/usr/bin/env bash
set -e

CLUSTER_NAME="${CLUSTER_NAME:-platform-demo}"
NAMESPACE="${NAMESPACE:-event-service}"
IMAGE="${IMAGE:-event-service:local}"

echo "==> Building image"
docker build -t "$IMAGE" apps/event-service

echo "==> Loading image into kind cluster: $CLUSTER_NAME"
kind load docker-image "$IMAGE" --name "$CLUSTER_NAME"

echo "==> Applying manifests"
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

echo "==> Restarting API deployment"
kubectl rollout restart deployment/event-service-api -n "$NAMESPACE"
kubectl rollout status deployment/event-service-api -n "$NAMESPACE"

echo "==> Current status"
kubectl get all -n "$NAMESPACE"
kubectl get ingress -n "$NAMESPACE"
