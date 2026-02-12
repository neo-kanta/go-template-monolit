# Architecture — Go Transfer Agent Platform (POC → MVP) 🏦

This document describes the **current POC** architecture and how it evolves into the **MVP** based on requirements:

- Gin API Gateway (REST + Swagger)
- gRPC core services (FND first)
- PostgreSQL (single DB, multi-schema)
- GORM migrations/ORM
- Maker–Checker (4-eyes principle)
- Thai localization (THB precision, Thai date formatting)
- Docker Compose local + Google Cloud deploy

---

## 1) Architecture Style

**Clean Architecture + Hexagonal (Ports & Adapters)**

- **Edge / Adapter:** `api-gateway` (Gin + Swagger, validation, auth)
- **Core / Use-cases:** `services/*` (gRPC servers, business rules)
- **Infrastructure adapters:** persistence (GORM/Postgres), messaging (future), external APIs (future)
- **Contracts:** Protobuf in `common/proto/*`

Key goal: juniors can create a new service by copying `_template` and implementing:

- domain
- usecases
- adapters (grpc + persistence)

---

## 2) System Overview (Target MVP)

```txt
Clients (Swagger/Bruno/Postman/curl)
        |
        | REST/HTTP
        v
+-----------------------+
| api-gateway (Gin)     |
| - Auth (JWT)          |
| - Validation          |
| - Swagger UI          |
| - Request mapping     |
+-----------------------+
        |
        | gRPC
        v
+-----------------------+       +-----------------------+
| services/FND (gRPC)   |       | services/Customer      |
| - Txn + balances      |       | - Customer master      |
| - Maker/Checker       |       | - KYC flags (future)   |
+-----------------------+       +-----------------------+
        |
        | GORM
        v
+-----------------------------------------------+
| PostgreSQL (single DB, multiple schemas)       |
| - engine.*  (customer, users, audit)           |
| - fnd.*     (fund, transactions, approvals)    |
+-----------------------------------------------+
```
