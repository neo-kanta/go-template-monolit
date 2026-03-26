# Go Clean Architecture Template

A starter repository for building robust Go microservices following Clean Architecture principles. It provides a modular, extensible monorepo foundation with an API gateway, gRPC backend services, and a structured set of shared common libraries.

## 🏗️ Architecture Overview

The workspace is organized into discrete components:

- **`api-gateway/`**: An HTTP gateway built on [Gin](https://gin-gonic.com/). It serves as the primary entry point for external clients, proxies REST requests to the underlying gRPC services, and provides a Swagger UI for API documentation.
- **`services/`**: The directory dedicated to backend gRPC microservices.
  - **`_template/`**: A canonical and strictly typed empty service scaffold. Use this blueprint to create new domain services to ensure architecture consistency in the repository.
  - **`sample/`**: A runnable baseline template service that demonstrates a complete implementation of the architectural pattern.
- **`common/`**: Shared resources strictly accessible across all components and services.
  - **`proto/`**: Protobuf definitions mapping the contract between gateways and microservices.
  - **`gen/`**: Auto-generated gRPC Go bindings based on the `.proto` files.
  - **`platform/`**: Standalone shared libraries containing utilities like logging, configuration loaders, connection poolers, and middleware.
- **`deploy/`**: Deployment configurations including local staging with Docker Compose (`deploy/compose/docker-compose.yml`).
- **`scripts/`**: Automation scripts to accelerate local development scenarios (e.g. `gen_proto.ps1`, `gen_swagger.ps1`).

## 🚀 Getting Started

### Prerequisites

- Go 1.25+
- [Buf](https://buf.build/docs/installation) CLI (for clean protobuf generation)
- Docker & Docker Compose (optional, for running with containers)

### Running Locally

1. **Generate Protobuf and Swagger Code:**
   Ensure your generated gRPC files and Swagger API documentations are up to date:
   ```powershell
   .\scripts\gen_proto.ps1
   .\scripts\gen_swagger.ps1
   ```

2. **Run Services Manually:**
   Launch the API Gateway and arbitrary downstream services natively:
   ```bash
   # Terminal 1: Run the HTTP API Gateway
   go run api-gateway/cmd/gateway/main.go

   # Terminal 2: Run the Template/Sample service
   go run services/sample/cmd/sample/main.go
   ```

3. **Run using Docker Compose:**
   Alternatively, run the entire ecosystem implicitly via docker setup:
   ```bash
   cd deploy/compose
   docker-compose up --build -d
   ```

## 🛠 How to Create a New Service

Creating a new robust microservice in this monorepo demands only a few explicit steps:

1. **Scaffold the Service:** Look at `services/_template`, safely copy the directory, and rename it to your service context (e.g., `services/my_service`).
2. **Define the API Contract:** Develop a brand new `.proto` file in `common/proto/my_service/v1/`.
3. **Generate gRPC Bindings:** Run `buf generate` or execute the script `.\scripts\gen_proto.ps1` to re-generate the underlying Go definitions into `common/gen/`.
4. **Implement the Clean Layers:** Fill out your business logic within the boundaries inside `services/my_service/internal`:
   - **Domain:** Pure business rules and enterprise entities.
   - **Usecase:** Application-specific workflows connecting to repositories.
   - **Adapter/Delivery:** External bindings representing the database, cache, or external SDKs.
5. **Route in the Gateway:** Inject your new GRPC connections inside `api-gateway` logic to securely proxy HTTP requests.

## 📚 Tech Stack

- **Target Language:** Go (1.25+)
- **Web Framework:** [Gin](https://gin-gonic.com/)
- **RPC Framework:** [gRPC](https://grpc.io/)
- **Protocol Buffers Management:** [Buf](https://buf.build/)
- **ORM Configuration:** [GORM](https://gorm.io/)
- **API Documentation:** [Swagger (Swaggo)](https://github.com/swaggo/swag) (alongside [grpc-gateway](https://grpc-ecosystem.github.io/grpc-gateway/))

## 📄 License

Review the `LICENSE` file for broad distribution details.
