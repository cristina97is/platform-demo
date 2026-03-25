#!/usr/bin/env bash
set -e

CLUSTER_NAME="${CLUSTER_NAME:-platform-demo}"
NAMESPACE="${NAMESPACE:-event-service}"
HELM_RELEASE="${HELM_RELEASE:-event-service}"
HELM_CHART="${HELM_CHART:-event-service-chart}"

IMAGE_REPOSITORY="${IMAGE_REPOSITORY:-ghcr.io/cristina97is/event-service}"
IMAGE_TAG="${IMAGE_TAG:-latest}"

LOCAL_IMAGE="${LOCAL_IMAGE:-false}"
LOCAL_IMAGE_NAME="${LOCAL_IMAGE_NAME:-event-service:local}"

if [ "$LOCAL_IMAGE" = "true" ]; then
  echo "==> Using local image flow"
  docker build -t "$LOCAL_IMAGE_NAME" apps/event-service
  kind load docker-image "$LOCAL_IMAGE_NAME" --name "$CLUSTER_NAME"

  helm upgrade --install "$HELM_RELEASE" "$HELM_CHART" -n "$NAMESPACE" \
    --set api.image.repository="event-service" \
    --set api.image.tag="local" \
    --set api.image.pullPolicy="IfNotPresent"
else
  echo "==> Using registry image flow"
  echo "==> Image: ${IMAGE_REPOSITORY}:${IMAGE_TAG}"

  helm upgrade --install "$HELM_RELEASE" "$HELM_CHART" -n "$NAMESPACE" \
    --set api.image.repository="$IMAGE_REPOSITORY" \
    --set api.image.tag="$IMAGE_TAG" \
    --set api.image.pullPolicy="Always"
fi

echo "==> Wait for rollout"
kubectl rollout status deployment/event-service-api -n "$NAMESPACE"
kubectl rollout status deployment/event-service-db -n "$NAMESPACE"

echo "==> Current pods"
kubectl get pods -n "$NAMESPACE"

echo "==> Health check"
sleep 3
curl -s http://event-service.127.0.0.1.sslip.io/healthz || true
echo
curl -s http://event-service.127.0.0.1.sslip.io/readyz || true
echo

echo "==> Done"
