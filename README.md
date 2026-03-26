# Go Clean Architecture Template
This is a starter repository for Go microservices following Clean Architecture.

## Architecture Structure
- `api-gateway/`: HTTP Gateway with Swagger UI, proxies requests to gRPC services.
- `services/`: Contains back-end gRPC services.
  - `sample/`: A runnable baseline template service implementation.
  - `_template/`: A canonical empty service scaffold. Use this to create new domain services.
- `common/`: Shared resources across all services.
  - `proto/`: Protobuf definitions.
  - `gen/`: Auto-generated gRPC code.
  - `platform/`: Shared libraries for logging, configuration, middleware, etc.

## How to create a new service
1. Copy `services/_template` to `services/my_new_service`.
2. Add your `.proto` file in `common/proto/my_new_service/v1/`.
3. Run `buf generate` or `scripts/gen_proto.ps1`.
4. Implement the service logically using Domain, Usecase, and Adapter patterns in `services/my_new_service/internal`.
5. Wire up the HTTP routing in `api-gateway` to proxy to your new service.

## Running Locally
- Ensure you have Go 1.22+ and Buf installed.
- Generate protos: `.\scripts\gen_proto.ps1`
- Build / Run API Gateway: `go run api-gateway/cmd/gateway/main.go`
- Build / Run Sample Service: `go run services/sample/cmd/sample/main.go`
