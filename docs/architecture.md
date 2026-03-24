# Architecture

## Overview

```text
Developer
   |
   v
Docker image (event-service:local)
   |
   v
kind Kubernetes cluster
   |
   +--> Ingress NGINX
   |       |
   |       v
   |   Service: event-service-api (ClusterIP)
   |       |
   |       v
   |   Pod: event-service-api
   |       |
   |       v
   |   Service: event-service-db
   |       |
   |       v
   |   Pod: PostgreSQL
   |
   +--> Prometheus
   |       |
   |       v
   |   Scrapes metrics from event-service
   |
   +--> Grafana
           |
           v
       Visualizes Prometheus metrics

User
  |
  v
http://event-service.127.0.0.1.sslip.io
