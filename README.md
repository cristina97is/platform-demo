# Platform Demo — Event Service (Go + Kubernetes)

## О проекте

Event Service — backend-сервис на Go для обработки событий (events), развернутый в Kubernetes.

Проект демонстрирует полный цикл разработки и эксплуатации сервиса:

* разработка API
* контейнеризация
* деплой в Kubernetes
* масштабирование
* мониторинг
* CI pipeline
* управление деплоем через Helm и Terraform

Проект реализован как локальная production-like среда (kind), без использования облака.

---

## Архитектура проекта

```
platform-demo/
├── .github/workflows/ci.yml
├── apps/event-service/
├── event-service-chart/
├── infrastructure/
│   ├── kind/
│   ├── prometheus/
│   └── terraform/
├── scripts/
├── load-test/
└── README.md
```

---

## Компоненты системы

* Go API (event-service)
* PostgreSQL (хранение данных)
* Kubernetes (оркестрация)
* Helm (управление деплоем)
* Terraform (управление deployment layer)
* Prometheus (метрики)
* Grafana (визуализация)

---

## Основной функционал API

* `POST /events` — создание события
* `GET /events` — получение списка событий
* `GET /healthz` — проверка состояния
* `GET /readyz` — проверка готовности
* `GET /metrics` — метрики Prometheus

---

## CI/CD

### CI (GitHub Actions)

Pipeline выполняет:

* проверку форматирования (gofmt)
* статический анализ (go vet)
* тесты
* сборку Docker image
* публикацию в GitHub Container Registry

Image публикуется в:

```
ghcr.io/cristina97is/event-service:latest
```

---

### Deployment

Деплой выполняется через Helm:

```bash
helm upgrade --install event-service event-service-chart -n event-service
```

Helm использует образ из registry:

```
ghcr.io/cristina97is/event-service:latest
```

---

### Local CD flow

Для локальной среды используется скрипт:

```bash
./scripts/release-local.sh
```

Поддерживаются два режима:

**1. Registry (основной)**
использует image из GHCR

**2. Local dev**

```bash
LOCAL_IMAGE=true ./scripts/release-local.sh
```

использует локально собранный образ через kind

---

## Terraform

Terraform используется для управления deployment layer:

* создание namespace
* установка Helm release

Запуск:

```bash
cd infrastructure/terraform
terraform init
terraform apply
```

Важно:

* state-файлы не хранятся в репозитории
* Terraform используется только для локального Kubernetes

---

## Запуск

### Kubernetes

```bash
./scripts/release-local.sh
```

---

## Доступ

```
http://event-service.127.0.0.1.sslip.io
```

---

## Нагрузочное тестирование

```bash
cd load-test
k6 run test.js
```

---

## Обоснование архитектурных решений

### Единый артефакт (CI → Registry → Deployment)

CI собирает Docker image и публикует его в GHCR.
Helm использует этот же image при деплое.

Результат:

* одинаковый образ во всех средах
* воспроизводимость
* отсутствие рассинхронизации

---

### Helm

Helm используется для шаблонизации Kubernetes ресурсов.

Результат:

* централизованная конфигурация
* удобное обновление
* повторное использование

---

### Terraform

Terraform управляет namespace и Helm release.

Результат:

* декларативное описание деплоя
* воспроизводимость окружения
* подготовка к cloud IaC

---

### Kubernetes

Используется как основной runtime.

Результат:

* автоматическое восстановление
* управление контейнерами
* масштабируемость

---

### Probes

* liveness — проверка процесса
* readiness — проверка готовности

Результат:

* корректный rollout
* стабильность

---

### HPA

Масштабирование по CPU.

Результат:

* адаптация к нагрузке
* эффективное использование ресурсов

---

### Persistent Storage

PostgreSQL использует PVC.

Результат:

* данные сохраняются при перезапуске

---

### Метрики

Сервис экспортирует `/metrics`.

Результат:

* наблюдаемость
* возможность анализа

---

### CI

CI проверяет код и собирает image.

Результат:

* стабильные сборки
* контроль качества

---

## Ограничения

* нет автоматического CD в удалённый кластер
* используется локальный кластер (kind)
* Terraform покрывает только deployment layer

---

## Возможные улучшения

* GitOps (ArgoCD)
* versioned Docker tags
* централизованный logging
* деплой в облако

---

## Автор

Kristina

