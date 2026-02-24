# How to Create a New FNDM Module

Follow these steps to create **FNDM003** (or any new module).

---

## Step 1 — Create the folder

```bash
mkdir services/fnd/fndm003
```

## Step 2 — Copy template files

```bash
cp services/fnd/_template/service.go.tmpl  services/fnd/fndm003/service.go
cp services/fnd/_template/handler.go.tmpl  services/fnd/fndm003/handler.go
```

## Step 3 — Search & Replace

In both files, replace:

| Find | Replace |
|------|---------|
| `TEMPLATE` | `fndm003` |
| `TemplateService` | Your service description, e.g. `Fund NAV Service` |

## Step 4 — Define your proto

Create `common/proto/fnd/v1/fndm003.proto` with your RPCs.

> **Until proto is split per-module**, add your RPCs to the existing `fnd.proto`:
>
> ```protobuf
> // ─── APIFNDM003: Your Module Name ───
> rpc YourMethod(YourRequest) returns (YourResponse);
> ```

Then generate code:

```bash
npx buf generate      # or: buf generate
```

## Step 5 — Implement business logic

Edit `fndm003/service.go`:
- Add methods that operate on your domain data
- Use `shared.GetCryName()` etc. from `services/fnd/shared/`
- Keep service methods pure (no gRPC types as params)

## Step 6 — Wire handler to service

Edit `fndm003/handler.go`:
- Each handler method maps proto request → service call → proto response
- Add structured logging with `h.log.Info("MethodName", ...)`

## Step 7 — Register in server.go

Edit `services/fnd/internal/adapter/grpc/server.go`:

```go
import "go-transfer-agent/services/fnd/fndm003"

type fndServer struct {
    // ...existing fields...
    m003 *fndm003.Handler  // ← add
}

// Add delegates for m003 methods:
func (s *fndServer) YourMethod(ctx context.Context, req *fndv1.YourRequest) (*fndv1.YourResponse, error) {
    return s.m003.YourMethod(ctx, req)
}

func NewServer(log *slog.Logger) *grpc.Server {
    fndSrv := &fndServer{
        // ...existing...
        m003: fndm003.NewHandler(log),  // ← add
    }
}
```

## Step 8 — Build & test

```bash
go build ./services/fnd/cmd/fnd/
go build ./api-gateway/cmd/gateway/
```

---

## File Convention

Each module has exactly **2 required files**:

| File | Responsibility |
|------|---------------|
| `service.go` | Business logic + data store |
| `handler.go` | gRPC bridge: proto ↔ service |

Optional:
| File | When |
|------|------|
| `seed.go` | Stub/seed data for POC testing |
| `repository.go` | Future: DB repository interface |

## Rules

1. **Never import another module directly** — use `shared/` for cross-cutting concerns
2. **Package name = folder name** (lowercase, e.g., `package fndm003`)
3. **Handler type must be exported** (`Handler`, not `handler`)
4. **Service methods should not take proto types** — accept Go primitives
