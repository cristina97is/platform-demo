# Platform Demo — Event Service (Go + Kubernetes)

## О проекте

Event Service — backend-сервис на Go для обработки событий, развернутый в Kubernetes.

Цель проекта — продемонстрировать не только реализацию API, но и полноценный инженерный подход к разработке и эксплуатации сервиса с учетом масштабируемости, отказоустойчивости, наблюдаемости, управления конфигурацией и автоматизации CI.

Проект реализован как локальная production-like среда на базе kind (Kubernetes in Docker), что позволяет воспроизводить реальные сценарии без использования облака.

Сервис может использоваться для логирования пользовательских действий, трекинга событий, обработки транзакционных операций и сбора телеметрии.

## Архитектура

Архитектура построена по принципу разделения ответственности:

- API отвечает за обработку запросов
- база данных отвечает за хранение
- Kubernetes управляет жизненным циклом
- мониторинг вынесен в отдельный слой

Компоненты системы:

- Go API (event-service)
- PostgreSQL (stateful storage)
- Kubernetes (оркестрация)
- Ingress (входной трафик)
- Prometheus (сбор метрик)
- Grafana (визуализация)

Сценарий работы:

1. Клиент отправляет запрос в `/events`
2. API обрабатывает запрос и сохраняет данные в PostgreSQL
3. API экспортирует метрики через `/metrics`
4. Prometheus собирает метрики
5. Grafana отображает их
6. Kubernetes управляет pod’ами и их состоянием

## Структура проекта

platform-demo/
├── .github/workflows/ci.yml
├── apps/event-service/
│   ├── cmd/api/main.go
│   ├── internal/
│   ├── docker-compose.yml
│   ├── .env.example
│   ├── k8s/
│   │   └── secret.example.yaml
│   └── Dockerfile
├── event-service-chart/
├── infrastructure/
│   ├── kind/
│   └── terraform/
├── docs/
├── load-test/
└── README.md

## Технологии и выбор стека

Go выбран как основной язык благодаря высокой производительности, простоте деплоя и широкому использованию в cloud-native экосистеме.

PostgreSQL используется как надежная ACID-база данных с поддержкой сложных запросов и аналитики.

Kubernetes применяется как стандартный оркестратор контейнеров, обеспечивающий масштабирование, self-healing и декларативное управление.

Prometheus и Grafana используются для мониторинга и визуализации метрик.

Helm используется как основной инструмент для управления deployment-конфигурацией.

Terraform применяется для управления deployment layer (namespace и Helm release).

## Основные инженерные решения

Конфигурация отделена от кода через ConfigMap и Secret, что позволяет изменять параметры без пересборки приложения и не хранить чувствительные данные в коде.

Для PostgreSQL используется PersistentVolumeClaim, что обеспечивает сохранность данных при рестарте pod’ов.

Добавлены readiness и liveness probes, чтобы Kubernetes корректно управлял жизненным циклом контейнеров.

API разворачивается с двумя репликами для обеспечения отказоустойчивости и балансировки нагрузки.

Используется стратегия RollingUpdate для обновлений без downtime.

Настроен Horizontal Pod Autoscaler, который масштабирует сервис в зависимости от нагрузки.

Используется PodDisruptionBudget для защиты от одновременного удаления всех pod’ов.

Сервис экспортирует метрики через `/metrics`, что позволяет интегрироваться с Prometheus.

## CI Pipeline

GitHub Actions pipeline выполняет:

- проверку форматирования (gofmt)
- статический анализ (go vet)
- запуск тестов (go test)
- lint Dockerfile (hadolint)
- проверку Helm chart (helm lint)
- валидацию Kubernetes-манифестов (kubeconform)
- сборку Docker-образа
- security scan образа (Trivy)
- публикацию образа в GitHub Container Registry (GHCR)

Публикуются два тега:

- latest
- commit SHA

Это обеспечивает единый deployment artifact и контроль качества как кода, так и инфраструктуры.

## Deployment

Основной способ деплоя — через Helm:

helm upgrade --install event-service event-service-chart -n event-service

Helm использует Docker-образ из GHCR, что связывает CI и runtime.

Terraform используется как слой управления deployment и управляет namespace и Helm release:

cd infrastructure/terraform
terraform init
terraform apply

Для локального запуска используется:

./scripts/release-local.sh

По умолчанию используется образ из GHCR, но доступен dev-режим с локальной сборкой через переменную LOCAL_IMAGE=true.

## Работа с секретами

Реальные значения секретов не хранятся в репозитории.

Для Kubernetes используется внешний Secret, который создаётся отдельно:

kubectl create secret generic event-service-secret \
  -n event-service \
  --from-literal=DB_USER=events \
  --from-literal=DB_PASSWORD=change-me \
  --from-literal=POSTGRES_USER=events \
  --from-literal=POSTGRES_PASSWORD=change-me

В репозитории присутствует только шаблон:

apps/event-service/k8s/secret.example.yaml

Для docker-compose используется файл `.env`, который не коммитится в репозиторий:

cd apps/event-service
cp .env.example .env
docker compose up --build

## Доступ

http://event-service.127.0.0.1.sslip.io

## Нагрузочное тестирование

cd load-test
k6 run test.js

## Ограничения

Проект развёрнут в локальном окружении, поэтому:

- отсутствует полноценный auto-CD из GitHub Actions
- нет удалённого Kubernetes-кластера
- управление секретами реализовано на уровне demo (без Vault/External Secrets)

## Что демонстрирует проект

- разработку сервиса на Go
- работу с Kubernetes
- Helm deployment
- Terraform integration
- CI pipeline с quality и security проверками
- безопасную работу с конфигурацией и секретами
- observability
- масштабирование и отказоустойчивость

## Возможные улучшения

- внедрение полноценного CD (self-hosted runner или GitOps)
- использование immutable image tag вместо latest
- интеграция с внешним secret manager
- переход на удалённый кластер
- централизованный logging (Loki)

## Автор

Kristina
