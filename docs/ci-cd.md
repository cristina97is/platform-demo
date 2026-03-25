# CI/CD Flow

## Overview

The project uses a mixed delivery model adapted for a local Kubernetes environment (kind).

### CI

GitHub Actions pipeline performs:

- gofmt check
- go vet
- go test
- Docker image build
- Docker image push to GitHub Container Registry (GHCR)

Published images:

- ghcr.io/cristina97is/event-service:latest
- ghcr.io/cristina97is/event-service:<commit-sha>

### Deployment

The main deployment path uses Helm and GHCR image:

```bash
helm upgrade --install event-service event-service-chart -n event-service
