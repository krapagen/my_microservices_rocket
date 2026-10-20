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
    order[Order Service\nHTTP :8080\nKafka producer/consumer]
    inventory[Inventory Service\ngRPC :50051]
    payment[Payment Service\ngRPC :50052]
    iam[IAM Service\ngRPC :50053]
    kafka[(Kafka)]
    assembly[Assembly Service\nKafka consumer/producer]

    client --> order
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
| **Order**     | HTTP API + Kafka producer/consumer + gRPC client (inventory, payment, IAM) | HTTP :8080 / gRPC :50051, :50052, :50053 |
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
├── platform/       # Platform-level utilities
├── deploy/
│   └── compose/    # Docker Compose for core infra and service stacks
├── migrations/     # SQL migrations (goose)
├── go.work         # Go workspace
└── Taskfile.yaml   # All dev tasks
```

## ⚙️ Tech Stack

- **Language:** Go 1.27
- **Transport:** gRPC (`protoc-gen-go`, `protoc-gen-go-grpc`) · REST (OpenAPI via `ogen`)
- **Messaging:** Apache Kafka (KRaft mode, no ZooKeeper)
- **Database:** PostgreSQL · Migrations via `goose`
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
task up-all
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
| `core`        | Kafka (KRaft) + Kafka UI |
| `observability` | OTel Collector + Jaeger + Prometheus + Grafana + Elasticsearch + Kibana + `es-init` (создаёт data stream `otel-logs`) |
| `order`       | Order service + PostgreSQL |
| `inventory`   | Inventory service + PostgreSQL |
| `iam`         | IAM service + PostgreSQL + Redis |

Kafka UI is available at **http://localhost:8088** (default port from `core.env`).

`payment` and `assembly` are run from source, not via Docker Compose.

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
- **Кастомные спаны** `order.Create` / `order.Pay` с атрибутами заказа и запиcью ошибок
- **Логи в ECS-формате** с `trace.id`/`span.id` — корреляция лог ↔ трейс в Kibana/Jaeger

> Для работы логов нужен data stream `otel-logs` — его создаёт one-shot контейнер `es-init` (запускается автоматически через `task up-core`).

## 📦 Go Workspace

The project uses a [Go workspace](https://go.dev/blog/get-familiar-with-workspaces) (`go.work`) to manage multiple modules locally without publishing them:

```
order · inventory · payment · assembly · iam · shared · platform
```

## 📄 License

MIT © [krapagen](https://github.com/krapagen)
