# Service Template — Architecture Reference
#
# This directory is a COPY-PASTE template for juniors creating new services.
# It is NOT meant to be compiled or run directly.
#
# To create a new service:
#   1. Copy this entire _template folder → services/<YourServiceName>
#   2. Rename cmd/service/ → cmd/<yourservice>/
#   3. Implement domain, usecases, and adapters
#   4. Add proto contract in common/proto/<yourservice>/v1/
#
# Structure:
#   cmd/service/main.go           ← gRPC server entry point
#   internal/
#     domain/                     ← entities, value objects, business rules
#     usecase/                    ← application use cases (ports)
#     adapter/
#       grpc/                     ← gRPC server implementation
#       persistence/              ← GORM repository implementation
