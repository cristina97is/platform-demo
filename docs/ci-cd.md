# CI/CD

## Goal

The goal of the delivery flow is to make application updates repeatable and predictable.

## Current workflow

At the moment the project uses a semi-manual delivery process:

1. Application code is updated locally.
2. A Docker image is built locally.
3. The image is loaded into the local `kind` cluster.
4. Kubernetes manifests are applied to deploy or update the application.
5. Service availability is verified through:
   - `/healthz`
   - `/readyz`
   - Ingress URL

## Example delivery flow

```text
Code change
   |
   v
Docker build
   |
   v
Local image: event-service:local
   |
   v
kind load docker-image
   |
   v
kubectl apply
   |
   v
Kubernetes Deployment update
   |
   v
Health check verification
