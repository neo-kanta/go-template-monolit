# Target Architecture

This repository adopts **Clean Architecture** combined with a **Monorepo** multi-service layout. The goal is to maximize isolation between domain logic, presentation layers, and infrastructure.

## Repository Layout
- `api-gateway/`: Acts as the central REST API frontend for all underlying services. Maps HTTP requests to gRPC calls. Includes routing, authentication (JWT), and Swagger generation.
- `services/`: Contains independent Microservices.
  - `_template/`: The canonical empty service scaffold. Duplicate it when creating a new domain service.
  - `sample/`: A working baseline service implementing `HealthCheck` and `Echo` endpoints. Use it for reference.
- `common/`: Shared code meant to be used by all microservices.
  - `proto/`: Protobuf definitions representing the contracts between API Gateway and individual services.
  - `gen/`: Auto-generated Go code from protobuf using `buf generate`.
  - `platform/`: Shared technical capabilities: config loaders, loggers, generic middlewares.

## Domain Service Internal Architecture (Clean Architecture)
Inside any given service (e.g., `services/sample`), the layout strictly enforces separation of concerns:

1. **`cmd/<service>/main.go` (Infrastructure Layer)**:
   - Wires up dependencies.
   - Bootstraps the gRPC server and Database.
   - Reads environment configurations.

2. **`internal/adapter/` (Interface Adapters)**:
   - **`grpc/`**: Handlers that receive gRPC requests, call usecases, and return gRPC responses.
   - **`persistence/`**: Database repositories. Implements domain interfaces using SQL/ORMs.

3. **`internal/usecase/` (Application Layer)**:
   - Contains business flow operations.
   - Orchestrates entities and repository instructions.
   - Has no knowledge of gRPC, REST, or SQL.

4. **`internal/domain/` (Enterprise Business Rules)**:
   - Holds core entities, value objects, domain logic, and interface definitions (e.g., repository interfaces).
   - This layer must not depend on ANY other internal package.

## Constraints
- **Direction of Dependency**: `Adapter` -> `Usecase` -> `Domain`.
- Do not import `adapter` code directly into `usecase` or `domain`. Only interface implementation should connect them.
- If an ORM model differs significantly from a Domain entity, use Data Transfer Objects (DTOs) to map them within the `persistence` layer.
