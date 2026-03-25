# Platform Demo — Event Service (Go + Kubernetes)

## О проекте

Event Service — backend-сервис на Go для обработки и хранения событий, развернутый в Kubernetes.

Проект демонстрирует полный контур вокруг сервиса: разработка приложения, контейнеризация, локальный запуск, деплой в Kubernetes, управление конфигурацией, observability, CI pipeline, deployment через Helm и управление deployment через Terraform.

Система разворачивается в локальном Kubernetes-кластере (kind).

## Функциональность

Сервис предоставляет HTTP API:

- GET /healthz — проверка процесса
- GET /readyz — проверка подключения к базе
- GET /events — список событий
- POST /events — создание события
- GET /metrics — метрики Prometheus

## Архитектура

Компоненты системы:

- Go API
- PostgreSQL
- Kubernetes
- Ingress
- Prometheus
- Grafana

Поток запроса:

1. Клиент → Ingress  
2. Ingress → Service  
3. Service → Pod  
4. API → PostgreSQL  
5. API → /metrics  
6. Prometheus → сбор метрик  

## Структура проекта

platform-demo/
├── .github/workflows/ci.yml
├── apps/event-service/
├── event-service-chart/
├── infrastructure/
│   ├── kind/
│   └── terraform/
├── scripts/
│   └── release-local.sh
├── load-test/
├── docs/
└── README.md

## Технологии

Go — API  
PostgreSQL — хранение данных  
Docker — контейнеризация  
Docker Compose — локальная среда  
Kubernetes — orchestration  
Helm — deployment  
Terraform — управление deployment layer  
Prometheus — метрики  
Grafana — визуализация  
GitHub Actions — CI  
GHCR — container registry  
k6 — нагрузочное тестирование  

## CI Pipeline

Pipeline выполняет:

- gofmt
- go vet
- go test
- hadolint
- helm lint
- kubeconform
- docker build
- Trivy scan
- push в GHCR

Публикуются два тега:

- latest
- <commit-sha>

latest используется только как fallback. Основной deployment выполняется по immutable tag (commit SHA).

## Deployment

Основной способ:

IMAGE_TAG=<commit-sha> ./scripts/release-local.sh

Пример:

IMAGE_TAG=db93f41 ./scripts/release-local.sh

Deployment использует образ:

ghcr.io/cristina97is/event-service:<commit-sha>

Это обеспечивает воспроизводимость, контроль версии и возможность отката.

Dev режим (только для локальной разработки):

LOCAL_IMAGE=true ./scripts/release-local.sh

Terraform:

cd infrastructure/terraform  
terraform init  
terraform apply -var="image_tag=<commit-sha>"

Terraform управляет namespace и Helm release.

## Локальный запуск

cd apps/event-service  
cp .env.example .env  
docker compose up --build  

## Секреты

Секреты не хранятся в репозитории.

Docker Compose использует файл .env:

cp .env.example .env

Kubernetes secret создаётся отдельно:

kubectl create secret generic event-service-secret \
  -n event-service \
  --from-literal=DB_USER=events \
  --from-literal=DB_PASSWORD=change-me \
  --from-literal=POSTGRES_USER=events \
  --from-literal=POSTGRES_PASSWORD=change-me

В репозитории хранятся только шаблоны:

- .env.example  
- secret.example.yaml  

## Доступ

http://event-service.127.0.0.1.sslip.io

## Нагрузочное тестирование

cd load-test  
k6 run test.js  

## Ограничения

- нет автоматического CD  
- локальный кластер  
- нет external secrets  

## Возможные улучшения

- автоматический deploy из CI  
- deployment только по immutable tag  
- external secrets  
- удалённый кластер  
- централизованный logging  

## Автор

Kristina
