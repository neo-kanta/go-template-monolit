# Go Transfer Agent Platform 🏦

A Go monorepo for a **Transfer Agent** platform built with Clean Architecture, Gin API Gateway, gRPC services, and PostgreSQL.

## Current Status

| Component | Status | Description |
|---|---|---|
| **API Gateway** | ✅ Working | Gin + Swagger UI + JWT auth |
| **FND Service** | ✅ Health only | gRPC health check endpoint |
| **Common Platform** | ✅ Scaffolded | Shared config, logger |

## Quick Start

### Prerequisites

- **Go 1.25+**
- **Docker + Docker Compose** (for containerized runs)
- **swag CLI** — `go install github.com/swaggo/swag/cmd/swag@latest`

### Run API Gateway (local)

```powershell
# Generate Swagger docs
swag init -g ./api-gateway/cmd/gateway/main.go -o ./api-gateway/docs

# Start gateway
go run ./api-gateway/cmd/gateway/main.go
```

Open **http://localhost:8080/swagger/index.html**

### Run FND Service (local)

```powershell
go run ./services/FND/cmd/fnd/main.go
```

Listens on `0.0.0.0:50051` (gRPC).

### Test FND Service (gRPC)

**Option 1: Use the helper script** (No extra tools needed)

```powershell
go run ./scripts/test_grpc_health.go
```

**Option 2: Use grpcurl**

```powershell
# Install
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Test
grpcurl -plaintext localhost:50051 list
grpcurl -plaintext localhost:50051 grpc.health.v1.Health/Check
```

### Run via Docker Compose

```powershell
docker compose -f deploy/compose/docker-compose.yml up --build
```

## Project Layout

```
go-transfer-agent/
├── api-gateway/          # Gin HTTP gateway (REST + Swagger + JWT)
├── services/
│   ├── FND/              # Fund service (gRPC) — MVP target
│   └── _template/        # Copy-paste scaffold for new services
├── common/
│   ├── proto/            # Protobuf contracts
│   ├── gen/              # Generated Go code (protoc)
│   └── platform/         # Shared wiring (config, logger, db, grpc)
├── deploy/
│   ├── docker/           # Dockerfiles
│   └── compose/          # Docker Compose
├── docs/                 # Architecture docs
└── scripts/              # Dev scripts (swagger gen, proto gen, tests)
```

## Architecture

**Clean Architecture + Hexagonal (Ports & Adapters)**

- **Edge:** API Gateway — auth, validation, Swagger, request routing
- **Core:** gRPC services — business logic, domain rules
- **Infra:** GORM/Postgres persistence, messaging (future)

See [docs/architecture.md](docs/architecture.md) for the full system diagram.

## Test Users (POC)

| Username | Password | Role |
|---|---|---|
| `admin` | `admin123` | admin |
| `user` | `user123` | user |

## Environment Variables

### API Gateway

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8080` | HTTP listen port |
| `JWT_SECRET` | `super-secret-poc-key-change-me` | HMAC HS256 secret |
| `JWT_EXPIRY_MINUTES` | `60` | Token TTL |
| `GIN_MODE` | `debug` | `release` for production |

### FND Service

| Variable | Default | Description |
|---|---|---|
| `GRPC_ADDR` | `0.0.0.0:50051` | gRPC listen address |
| `DB_DSN` | _(empty)_ | PostgreSQL connection string |
| `LOG_LEVEL` | `info` | `debug` / `info` / `warn` / `error` |
