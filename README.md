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
ta-platform/
  go.mod
  go.sum
  README.md
  docs/
    architecture.md
    api/
      postman_collection.json
      bruno/
    thai/
      README.th.md
      terminology.th.md

  common/
    proto/
      fnd/v1/
        fnd.proto
        common.proto
      customer/v1/
        customer.proto
    gen/                 # generated Go code (protoc/buf)
    middleware/
      authclaims.go      # shared claims struct + helpers (no framework deps)
    platform/
      config/            # viper/envconfig style loader (your choice)
      logger/            # zap/slog wrapper
      db/
        postgres.go      # opens *gorm.DB
      grpc/
        interceptors.go  # shared interceptors: auth, logging, request-id
      errs/
        grpcerrs.go      # map domain errors -> gRPC status

  api-gateway/
    cmd/gateway/
      main.go
    internal/
      http/
        router.go
        middleware/
          jwt.go          # hard-coded JWT middleware (MVP)
          requestid.go
        handlers/
          fnd_handlers.go # REST->gRPC mapping
          health.go
      grpcclient/
        fnd.go

  services/
    _template/            # copy this to create new service
      cmd/service/main.go
      internal/
        adapter/
          grpc/
          persistence/
        domain/
        usecase/
      migrations/
      test/

    FND/
      cmd/fnd/main.go
      internal/
        adapter/
          grpc/
            server.go
            interceptor_auth.go
            mapper.go
          persistence/
            gorm/
              db.go
              models.go
              repos.go
        domain/
          money/
            thb.go
          transfer/
            entity.go
            status.go
          approval/
            maker_checker.go
        usecase/
          submit.go
          approve.go
          query.go
      migrations/
        automigrate.go     # uses GORM AutoMigrate (MVP) citeturn1search0
      test/
        unit/
        integration/

  deploy/
    docker/
      Dockerfile.gateway
      Dockerfile.fnd
    compose/
      docker-compose.yml
    gcp/
      cloudrun.md
      gke.md
      env.example

  scripts/
    gen_proto.sh
    lint.sh
    test.sh

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
