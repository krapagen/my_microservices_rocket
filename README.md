<!-- Coverage badge is auto-updated via GitHub Actions -->
![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/krapagen/52a65d3bc94b4dfc010de299283c7f78/raw/coverage.json)
![Go Version](https://img.shields.io/badge/Go-1.27-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-green)

<h1 align="center">🚀 My Microservices Rocket</h1>

<p align="center">
  A Go microservices platform with HTTP, gRPC, Kafka, PostgreSQL, and clean architecture.
</p>

---

## 📐 Architecture

<div align="center">

```mermaid
flowchart TB
    client[Client]
    nginx[Nginx LB :8080]
    order[Order Service ×3\nHTTP :8080\nKafka producer/consumer]
    ratelimit[(Redis\nrate limiter)]
    inventory[Inventory Service\ngRPC :50051]
    payment[Payment Service\ngRPC :50052]
    iam[IAM Service\ngRPC :50053]
    kafka[(Kafka)]
    assembly[Assembly Service\nKafka consumer/producer]

    client --> nginx
    nginx --> order
    order -->|GCRA check| ratelimit
    order -->|gRPC| inventory
    order -->|gRPC| payment
    order -->|gRPC| iam
    inventory -->|gRPC| iam
    order -->|order.paid| kafka
    kafka -->|order.paid| assembly
    assembly -->|assembly.ship-assembled| kafka
    kafka -->|assembly.ship-assembled| order
```

</div>

## 🔄 Event Flow

1. `Order` handles payment and publishes `order.paid` to Kafka.
2. `Assembly` consumes `order.paid`, assembles the order, and publishes `assembly.ship-assembled`.
3. `Order` consumes `assembly.ship-assembled`, commits inventory, and marks the order as assembled.

| Service       | Role                         | Protocol / Port |
|---------------|------------------------------|-----------------|
| **Nginx**     | Load balancer перед репликами Order | HTTP :8080      |
| **Order**     | HTTP API (3 replicas) + rate limiter + Kafka producer/consumer + gRPC client (inventory, payment, IAM) | HTTP :8080 (internal) |
| **Assembly**  | Kafka consumer/producer      | Kafka only      |
| **Inventory** | gRPC service                  | gRPC :50051     |
| **Payment**   | gRPC service                  | gRPC :50052     |
| **IAM**       | Authentication & session management | gRPC :50053     |

## 🗂️ Project Structure

```
my_microservices_rocket/
├── order/          # Order service — HTTP API + Kafka integration
├── inventory/      # Inventory service — gRPC
├── payment/        # Payment service — gRPC
├── assembly/       # Assembly service — Kafka consumer/producer
├── iam/            # IAM service — gRPC authentication & sessions
├── shared/         # Shared code: .proto definitions, OpenAPI specs, generated clients
│   ├── proto/      # Protobuf sources (buf managed)
│   └── api/        # OpenAPI specs
├── platform/       # Platform-level utilities (metrics, tracing, ratelimit, redis, closer)
├── deploy/
│   ├── compose/    # Docker Compose stacks (core, observability, services, per-service DBs)
│   ├── nginx/      # Nginx load balancer config
│   ├── otel/       # OTel Collector config
│   ├── grafana/    # Grafana provisioning (datasources + dashboards)
│   └── dockerfiles/# Per-service Dockerfiles
├── migrations/     # SQL migrations (goose)
├── go.work         # Go workspace
└── Taskfile.yaml   # All dev tasks
```

## ⚙️ Tech Stack

- **Language:** Go 1.27
- **Transport:** gRPC (`protoc-gen-go`, `protoc-gen-go-grpc`) · REST (OpenAPI via `ogen`)
- **Messaging:** Apache Kafka (KRaft mode, no ZooKeeper)
- **Database:** PostgreSQL · Migrations via `goose`
- **Cache / Rate limit:** Redis (`redis_rate`, GCRA) — общий лимит для реплик Order + сессии IAM
- **Load balancing:** Nginx перед репликами Order (`--scale order-service=3`)
- **Testing:** `testify` · `mockery` v3 · `testcontainers` · race detector
- **Code generation:** `buf` (proto) · `ogen` (OpenAPI)
- **Linting / Formatting:** `golangci-lint` · `gofumpt` · `gci`
- **Task runner:** `task` (Taskfile)
- **Observability:** OpenTelemetry (OTLP gRPC) · OTel Collector · Jaeger (трейсы) · Prometheus + Grafana (метрики) · Elasticsearch + Kibana (логи)

## 🚀 Getting Started

### Prerequisites

- Go 1.27+
- Docker & Docker Compose
- [`task`](https://taskfile.dev) CLI

### 1. Install all dev tools

```bash
task setup
```

### 2. Start infrastructure (Kafka + Observability)

```bash
task up-core
```

Поднимает Kafka и весь observability-стек: OTel Collector, Jaeger, Prometheus, Grafana, Elasticsearch, Kibana.

### 3. Run services

```bash
# Each in a separate terminal
task run:inventory   # gRPC  :50051
task run:payment     # gRPC  :50052
task run:order       # HTTP  :8080
task run:assembly    # Kafka consumer/producer
task run:iam         # gRPC  :50053
```

Or start everything with Docker Compose:

```bash
task up-all            # инфраструктура + все сервисы; Order ×3 за Nginx :8080
task up-all-with-db    # то же + per-service БД и миграции
```

### 4. Stop everything

```bash
task down-all
```

## 🧪 Testing

```bash
# Unit tests (with race detector)
task test

# Unit tests + coverage report
task test:coverage

# Open HTML coverage report
task coverage:html

# API integration tests
task test:api

# E2E tests (Postgres + Kafka via testcontainers)
task test:e2e
```

> Coverage threshold is enforced at **95%** on business logic packages.

## 🔧 Code Generation

```bash
# Generate Go code from .proto files
task proto:gen

# Generate Go code from OpenAPI specs
task ogen:gen

# Generate all (proto + openapi)
task gen

# Regenerate mocks
task mocks:gen
```

## 🗄️ Database Migrations

```bash
# Order service
task migrate:order:up
task migrate:order:down
task migrate:order:status

# Inventory service
task migrate:inventory:up
task migrate:inventory:down
task migrate:inventory:status

# IAM service
task migrate:iam:up
task migrate:iam:down
task migrate:iam:status
```

## 🧹 Code Quality

```bash
# Format code (gofumpt + gci import sorting)
task format

# Run linter
task lint

# Update all dependencies
task deps:update
```

## 🐳 Docker Compose Services

| Compose stack | Contents |
|---------------|----------|
| `core`        | Kafka (KRaft) + Kafka UI + `redis-ratelimit` |
| `observability` | OTel Collector + Jaeger + Prometheus + Grafana + Elasticsearch + Kibana + `es-init` (создаёт data stream `otel-logs`) |
| `services`    | Все микросервисы + Nginx LB (`up-all` поднимает Order ×3) |
| `order`       | PostgreSQL для Order |
| `inventory`   | PostgreSQL для Inventory |
| `iam`         | PostgreSQL + Redis (сессии) для IAM |

Kafka UI is available at **http://localhost:8088** (default port from `core.env`).

## 📊 Observability

Все сервисы экспортируют телеметрию по OTLP gRPC в OTel Collector (`localhost:4317`), который маршрутизирует:

| Сигнал | Хранилище | UI |
|--------|-----------|-----|
| Трейсы | Jaeger (OTLP) | http://localhost:16686 |
| Метрики | Prometheus (Remote Write) | http://localhost:9090 / Grafana http://localhost:3000 (admin/admin) |
| Логи | Elasticsearch (data stream `otel-logs`, ECS mapping) | Kibana http://localhost:5601 |

Что подключено:

- **`tracing.InitTracer`** в `app.go` каждого сервиса — TracerProvider + W3C TraceContext propagator для сквозных трейсов
- **`otelgrpc`** stats handlers на gRPC-серверах и клиентах — автоматические span'ы на каждый вызов
- **`otelhttp`** middleware на HTTP-сервере Order — спаны + метрики `http_server_request_duration_seconds`
- **Бизнес-метрики OrderService**: `orders_created_total`, `orders_paid_total`, `orders_revenue_total` (сумма в копейках)
- **`rate_limit_requests_total`** — решения лимитера (`result`/`path`/`instance`), считает и отказы 429, которые до `otelhttp` не доходят
- **`service.instance.id`** = hostname в resource метрик — каждая реплика пишет свою серию (иначе реплики сливаются в одну и `rate()`/`irate()` врут)
- **Кастомные спаны** `order.Create` / `order.Pay` с атрибутами заказа и запиcью ошибок
- **Логи в ECS-формате** с `trace.id`/`span.id` — корреляция лог ↔ трейс в Kibana/Jaeger
- **Дашборды Grafana** (provisioned из `deploy/grafana/dashboards`): Order Service и Rate Limiter (`http://localhost:3000/d/rate-limiter` — всего/пропущено/отклонено в req/s)

> Для работы логов нужен data stream `otel-logs` — его создаёт one-shot контейнер `es-init` (запускается автоматически через `task up-core`).

## 🚦 Rate Limiting

HTTP API Order защищён распределённым rate limiter'ом (`platform/pkg/ratelimit`, Redis + `redis_rate`, алгоритм GCRA):

- **Лимит общий для всех реплик** — состояние в `redis-ratelimit`, ключ — request path
- **`rate: 100` req/s, `burst: 200`** — в `order/config.*.yaml` (`rate_limit:`)
- Превышение → `429 Too Many Requests` (до `otelhttp`/auth — поэтому 429 не попадает в `http_server_*`)
- **Fail-open**: недоступный Redis → пропуск + `result="fail_open"` в метрике

Нагрузочный тест: `task load:http` (vegeta, 50 → 500 RPS за 2.5 мин) — линии «пропущено ≈100/с / отклонено» на дашборде Rate Limiter.

## 📦 Go Workspace

The project uses a [Go workspace](https://go.dev/blog/get-familiar-with-workspaces) (`go.work`) to manage multiple modules locally without publishing them:

```
order · inventory · payment · assembly · iam · shared · platform
```

## 📄 License

MIT © [krapagen](https://github.com/krapagen)
