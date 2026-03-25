# Platform Demo — Event Service (Go + Kubernetes)

## О проекте

Event Service — backend-сервис на Go для обработки и хранения событий, развернутый в Kubernetes.

Проект демонстрирует базовый платформенный контур вокруг сервиса:

- разработка приложения
- контейнеризация
- локальный запуск
- деплой в Kubernetes
- управление конфигурацией
- observability
- CI pipeline
- deployment через Helm
- управление deployment через Terraform

Система разворачивается в локальном Kubernetes-кластере на базе kind.

---

## Функциональность

Сервис предоставляет HTTP API:

- `GET /healthz` — проверка доступности процесса
- `GET /readyz` — проверка доступности базы данных
- `GET /events` — получение списка событий
- `POST /events` — создание события
- `GET /metrics` — метрики Prometheus

Пример payload:

```json
{
  "user_id": 1,
  "type": "bet",
  "amount": 100
}
```

---

## Архитектура

Система состоит из следующих компонентов:

- API-сервис (Go)
- PostgreSQL
- Kubernetes
- Ingress
- Prometheus
- Grafana

### Поток запроса

1. Клиент → Ingress  
2. Ingress → Service  
3. Service → Pod (API)  
4. API → PostgreSQL  
5. API → `/metrics`  
6. Prometheus → сбор метрик  
7. Grafana → визуализация  

---

## Структура проекта

```text
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
```

---

## Технологии и обоснование выбора

### Go

Используется для реализации API.

Причины выбора:

- низкие накладные расходы рантайма
- простой деплой (один бинарник)
- стандартная библиотека покрывает HTTP-сервер и работу с сетью
- распространён в cloud-native инструментах

Применение в проекте:

- HTTP сервер (`net/http`)
- обработка запросов
- экспорт метрик
- конфигурация через env

---

### PostgreSQL

Используется как основное хранилище.

Причины выбора:

- транзакционность (ACID)
- предсказуемое поведение
- стандартное решение для backend-сервисов

Применение:

- хранение событий
- инициализация через SQL-migration
- доступ через pgx

---

### Docker

Используется для упаковки приложения.

Причины выбора:

- единый формат артефакта
- воспроизводимость окружения
- одинаковый запуск в CI и Kubernetes

Применение:

- multi-stage build
- минимальный runtime-образ

---

### Docker Compose

Используется для локального запуска.

Причины выбора:

- быстрый запуск всех компонентов
- изоляция окружения разработки
- отсутствие зависимости от Kubernetes

Применение:

- API
- PostgreSQL
- Prometheus
- Grafana

---

### Kubernetes

Используется как runtime-платформа.

Причины выбора:

- управление контейнерами
- автоматическое восстановление
- масштабирование
- декларативная конфигурация

Применение:

- Deployment
- Service
- Ingress
- ConfigMap
- Secret
- PVC
- HPA
- PDB

---

### kind

Используется как локальный кластер Kubernetes.

Причины выбора:

- запуск в Docker
- отсутствие зависимости от облака
- возможность воспроизвести k8s-сценарии локально

---

### Helm

Используется для управления deployment.

Причины выбора:

- шаблонизация Kubernetes manifests
- централизованная конфигурация
- управление релизами

Применение:

- chart `event-service-chart`
- управление image, replicas, ingress, HPA

---

### Terraform

Используется для управления deployment layer.

Причины выбора:

- декларативное описание инфраструктуры
- воспроизводимость состояния

Применение:

- создание namespace
- управление Helm release

---

### Prometheus

Используется для сбора метрик.

Причины выбора:

- стандарт для Kubernetes
- pull-модель
- интеграция с Go

Применение:

- scrape `/metrics`
- алерты

---

### Grafana

Используется для визуализации.

Причины выбора:

- стандартный UI для метрик
- удобные дашборды

---

### GitHub Actions

Используется как CI.

Причины выбора:

- интеграция с GitHub
- простая настройка
- прозрачность pipeline

---

### GHCR

Используется как container registry.

Причины выбора:

- хранение Docker-образов
- интеграция с GitHub Actions

---

### k6

Используется для нагрузочного тестирования.

Причины выбора:

- простой синтаксис
- подходит для HTTP-нагрузки

---

## Реализация

### Конфигурация

Используются переменные окружения:

- `PORT`
- `DB_HOST`
- `DB_PORT`
- `DB_NAME`
- `DB_USER`
- `DB_PASSWORD`
- `DB_SSLMODE`

В Kubernetes:

- ConfigMap — для конфигурации
- Secret — для чувствительных данных

---

### Хранение данных

- PostgreSQL
- PersistentVolumeClaim
- сохранение данных между рестартами

---

### Health checks

- `healthz` — проверка процесса
- `readyz` — проверка подключения к БД

Используется в liveness и readiness probes.

---

### Масштабирование

- 2 replicas API
- Horizontal Pod Autoscaler (2–5 pod)

---

### Обновления

- RollingUpdate
- постепенная замена pod’ов

---

### Надёжность

- PodDisruptionBudget
- защита от полной потери сервиса

---

### Метрики

- `/metrics`
- Prometheus scrape
- базовые показатели:
  - количество запросов
  - latency

---

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

Теги:

- `latest`
- `<commit-sha>`

---

## Deployment

### Helm

```bash
helm upgrade --install event-service event-service-chart -n event-service
```

Используется образ из GHCR.

---

### Terraform

```bash
cd infrastructure/terraform
terraform init
terraform apply
```

Управляет:

- namespace
- Helm release

---

## Локальный запуск

### Docker Compose

```bash
cd apps/event-service
cp .env.example .env
docker compose up --build
```

---

### Kubernetes

```bash
./scripts/release-local.sh
```

Dev режим:

```bash
LOCAL_IMAGE=true ./scripts/release-local.sh
```

---

## Секреты

Секреты не хранятся в Git.

Создаются вручную:

```bash
kubectl create secret generic event-service-secret \
  -n event-service \
  --from-literal=DB_USER=events \
  --from-literal=DB_PASSWORD=change-me \
  --from-literal=POSTGRES_USER=events \
  --from-literal=POSTGRES_PASSWORD=change-me
```

---

## Доступ

```text
http://event-service.127.0.0.1.sslip.io
```

---

## Нагрузочное тестирование

```bash
cd load-test
k6 run test.js
```

---

## Ограничения

- нет автоматического CD
- локальный кластер
- нет внешнего secret manager

---

## Возможные улучшения

- CD pipeline
- immutable image tags
- external secrets
- удалённый кластер
- централизованный logging

---

## Автор

Kristina

