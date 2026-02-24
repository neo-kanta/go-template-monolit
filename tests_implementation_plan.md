# Test Implementation Plan: FND Modules (FNDM002 - FNDM013)

This plan outlines the strategy for writing Unit, Logic, and Integration tests for the Foundational Modules using GORM.

## 1. Testing Strategy & Challenges
The implemented modules (`FNDM002` through `FNDM013`) heavily rely on GORM for database interactions. Furthermore, the GORM models utilize PostgreSQL-specific features (e.g., `gen_random_uuid()` as default values).
Because of this, using in-memory SQLite for testing will cause migration failures. 

We have two primary options for Integration/Logic Testing:
1. **Mocking (sqlmock)**: Intercept SQL queries. (Pros: Fast, no DB needed. Cons: Immensely tedious to maintain for complex 4-table join queries and nested structs, doesn't test real DB behavior).
2. **Real PostgreSQL Database**: Connect to a real database purely for testing. (Pros: 100% accurate integration testing. Cons: Requires an active DB).
    *   **Option 2A**: Require a local `TEST_DB_DSN` environment variable to connect to an existing local Postgres instance.
    *   **Option 2B**: Introduce `testcontainers-go` to automatically spin up an ephemeral Dockerized PostgreSQL instance for the duration of the tests. (Highly recommended for CI/CD).

## 2. Proposed Changes

### Setup Test Infrastructure
*   Create a shared `services/fnd/testutil` package.
*   Implement `SetupTestDB()` that connects to the test database, runs `AutoMigrate` for the required tables, and provides a clean `*gorm.DB` instance.
*   Implement `SeedData(db, models...)` helper to easily inject test fixtures before each test suite.

### Module Tests (For all FNDM002 - FNDM013)
*   **Create `service_test.go`** (Integration & Logic):
    *   Initialize the test DB.
    *   Seed complex relational data based on the spec sheets (Master + repeating details like Fees, Groups, Agents).
    *   Call the `service.go` methods.
    *   **Logic Verification**: Assert that the highly nested Protobuf Response structure (e.g., `FundFeeSubListResDTO`) is correctly aggregated and mapped from the flat relational rows. Provide coverage for custom logic like FNDM004's string splitting logic.
*   **Create `handler_test.go`** (Unit):
    *   Instantiate the gRPC Handler with the Service.
    *   Verify request routing and standard input/output mapping via the gRPC endpoint.

## User Review Required

> [!IMPORTANT]
> **Database Testing Preference**
> How would you like to handle the PostgreSQL dependency for these tests?
> 1. Should I write the tests to connect to your existing local Postgres instance via a connection string (e.g., `localhost:5432`)?
> 2. Should I add the `github.com/testcontainers/testcontainers-go/modules/postgres` dependency so the tests automatically spin up their own isolated PostgreSQL Docker container?
> 3. Should I try to mock everything with `sqlmock` instead (not recommended due to complexity)?
