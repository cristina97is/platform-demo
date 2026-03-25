#!/bin/bash
set -e

NAMESPACE="event-service"
RELEASE="event-service"
CHART="./event-service-chart"

IMAGE_TAG=${IMAGE_TAG:-latest}

if [ "$LOCAL_IMAGE" = "true" ]; then
  IMAGE_NAME="event-service:local"

  docker build -t $IMAGE_NAME ./apps/event-service
  kind load docker-image $IMAGE_NAME

  helm upgrade --install $RELEASE $CHART \
    -n $NAMESPACE \
    --create-namespace \
    --set api.image.repository=event-service \
    --set api.image.tag=local \
    --set api.image.pullPolicy=IfNotPresent

else
  helm upgrade --install $RELEASE $CHART \
    -n $NAMESPACE \
    --create-namespace \
    --set api.image.repository=ghcr.io/cristina97is/event-service \
    --set api.image.tag=$IMAGE_TAG \
    --set api.image.pullPolicy=Always
fi
