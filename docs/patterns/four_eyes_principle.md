# 4-Eyes Principle (Maker-Checker) — Implementation Pattern

> **Purpose**: This document defines the exact pattern for implementing the 4-Eyes Principle (Maker-Checker workflow) in the Go Transfer Agent FND microservice. Follow this pattern precisely when adding approval workflows to any FND module.

---

## Architecture Overview

```
┌──────────────┐    JWT Token    ┌─────────────────┐    gRPC ctx     ┌──────────────────┐
│   Client     │───────────────►│ AuthInterceptor │──────────────►│  Service Layer   │
│ (Postman/    │                │ (extracts Sub,  │                │  (SaveFundInfo,  │
│  curl/Bruno) │                │  Role from JWT) │                │   ApproveFundInfo)│
└──────────────┘                └─────────────────┘                └──────────────────┘
                                        │                                   │
                                        │ shared.NewAuthContext(            │ shared.GetUsernameFromCtx(ctx)
                                        │   ctx, sub, role)                 │ shared.GetRoleFromCtx(ctx)
                                        ▼                                   ▼
                                ┌─────────────────┐                ┌──────────────────┐
                                │ shared/authctx.go│                │  DB: AuditFields │
                                │ (context helpers)│                │  + MakerChecker  │
                                └─────────────────┘                │    Fields        │
                                                                   └──────────────────┘
```

## Status Flow (Simplified from Legacy C# FlowAPI)

```
DRAFT ──► PENDING_APPROVAL ──► APPROVED
                            └──► REJECTED
```

Legacy C# equivalents:
| Go Status            | C# FlowStatus     | C# Code |
|----------------------|--------------------|---------|
| `DRAFT`              | `Initiated`        | 100     |
| `PENDING_APPROVAL`   | `Pending`          | 200     |
| `APPROVED`           | `Approved`         | 300     |
| `REJECTED`           | `Rejected`         | 400     |

---

## File Structure

```
common/platform/model/base.go          # MakerCheckerFields, ApprovalStatus constants
services/fnd/shared/authctx.go         # Context helpers (GetUsernameFromCtx, GetRoleFromCtx, NewAuthContext)
services/fnd/internal/adapter/grpc/
  auth_interceptor.go                  # JWT → shared.NewAuthContext(ctx, sub, role)
services/fnd/{module}/
  db/models.go                         # AuditFields embeds models.MakerCheckerFields
  service.go                           # SaveXxx(ctx, req) + ApproveXxx(ctx, req)
  handler.go                           # Passes ctx to service
```

---

## Step-by-Step Implementation Pattern

### 1. Embed MakerCheckerFields in DB Models

In `services/fnd/{module}/db/models.go`, add `models.MakerCheckerFields` to the `AuditFields` struct:

```go
import models "go-transfer-agent/common/platform/model"

type AuditFields struct {
    // ... existing audit fields ...
    models.MakerCheckerFields  // ← Add this line
}
```

This adds the following columns to every table that uses `AuditFields`:
- `MakerID` (varchar(50), not null, default '')
- `CheckerID` (varchar(50), nullable)
- `Status` (varchar(20), not null, default 'DRAFT')
- `Remark` (text, nullable)

### 2. Modify Create/Save to Auto-Set Maker

In `services/fnd/{module}/service.go`:

```go
import (
    models "go-transfer-agent/common/platform/model"
    "go-transfer-agent/services/fnd/shared"
)

// SaveXxx creates a new record with 4-Eyes fields auto-set.
func (s *Service) SaveXxx(ctx context.Context, req *proto.SaveXxxRequest) *proto.SaveXxxResponse {
    // 1. Extract MakerID from JWT context
    makerID := shared.GetUsernameFromCtx(ctx)
    if makerID == "" {
        makerID = "system" // Fallback for tests
    }

    // 2. Map proto to DB model
    record := mapToDBModel(req)

    // 3. Set 4-Eyes fields
    record.MakerID = makerID
    record.Status = models.StatusPendingApproval

    // 4. Persist
    if err := s.db.Create(&record).Error; err != nil {
        return errorResponse(err)
    }
    return successResponse()
}
```

### 3. Implement Approve/Reject Endpoint

```go
func (s *Service) ApproveXxx(ctx context.Context, req *proto.ApproveXxxRequest) *proto.SaveXxxResponse {
    // 1. Extract CheckerID from JWT context (or request)
    checkerID := req.CheckerId
    if checkerID == "" {
        checkerID = shared.GetUsernameFromCtx(ctx)
    }
    if checkerID == "" {
        return validationError("CheckerID is required")
    }

    // 2. Load the record
    var record DBModel
    if err := s.db.Where(...).First(&record).Error; err != nil {
        return notFoundError()
    }

    // 3. Validate status is PENDING_APPROVAL
    if record.Status != models.StatusPendingApproval {
        return invalidStatusError(record.Status)
    }

    // 4. Enforce 4-Eyes: MakerID != CheckerID
    if record.MakerID == checkerID {
        return fourEyesViolationError(record.MakerID)
    }

    // 5. Transition status
    newStatus := models.StatusRejected
    if req.IsApproved {
        newStatus = models.StatusApproved
    }

    // 6. Update
    updates := map[string]interface{}{
        "CheckerID": checkerID,
        "Status":    newStatus,
        "Remark":    req.Remark,
    }
    s.db.Model(&DBModel{}).Where(...).Updates(updates)
    return successResponse()
}
```

### 4. Wire Handler to Pass Context

In `services/fnd/{module}/handler.go`:

```go
func (h *Handler) SaveXxx(ctx context.Context, req *proto.SaveXxxRequest) (*proto.SaveXxxResponse, error) {
    return h.svc.SaveXxx(ctx, req), nil  // ← Pass ctx through
}

func (h *Handler) ApproveXxx(ctx context.Context, req *proto.ApproveXxxRequest) (*proto.SaveXxxResponse, error) {
    return h.svc.ApproveXxx(ctx, req), nil
}
```

### 5. Add Proto Definitions

In `common/proto/fnd/v1/{module}.proto`:

```protobuf
// Add to your DTO message:
message YourDTO {
    // ... existing fields ...
    string maker_id = N;
    string checker_id = N+1;
    string approval_status = N+2;
}

// Add approval request:
message ApproveXxxRequest {
    string sys_co_id = 1;
    string primary_key = 2;
    string checker_id = 3;
    bool is_approved = 4;
    string remark = 5;
}

// Add RPC:
rpc ApproveXxx(ApproveXxxRequest) returns (SaveXxxResponse) {
    option (google.api.http) = {
        post: "/api/fndmXXX/approve"
        body: "*"
    };
}
```

### 6. Update Tests

```go
func (suite *ServiceTestSuite) TestSaveXxx_AutoSetsMakerAndPending() {
    svc := NewService(testutil.SetupTestDB(suite.T()), testutil.NewTestLogger())
    
    // Use context.Background() — MakerID will default to "system"
    res := svc.SaveXxx(context.Background(), &proto.SaveXxxRequest{...})
    assert.True(suite.T(), res.Success)
    
    // Verify MakerID and Status were auto-set in DB
    // ...
}

func (suite *ServiceTestSuite) TestApproveXxx_FourEyesViolation() {
    svc := NewService(...)
    
    // Create with maker context
    makerCtx := shared.NewAuthContext(context.Background(), "junior_officer", "maker")
    svc.SaveXxx(makerCtx, &proto.SaveXxxRequest{...})
    
    // Try to approve as same user — should fail
    res := svc.ApproveXxx(makerCtx, &proto.ApproveXxxRequest{...})
    assert.False(suite.T(), res.Success)
    assert.Equal(suite.T(), "FOUR_EYES_VIOLATION", res.ReturnCode)
}
```

---

## Key Constraints

1. **MakerID != CheckerID**: The person who creates a record cannot approve it.
2. **Status must be PENDING_APPROVAL**: Only records in `PENDING_APPROVAL` status can be approved or rejected.
3. **JWT Context**: The `MakerID` and `CheckerID` are automatically extracted from the authenticated JWT token's `sub` claim.
4. **No Multi-Step**: This is a simplified single-step approval. No `FlowStep` audit trail table is needed for the MVP.
5. **THB Only**: All currency fields must be `"THB"` (enforced by `THBValidationInterceptor`).

## Important: Avoid Import Cycles

The context helpers (`GetUsernameFromCtx`, `GetRoleFromCtx`, `NewAuthContext`) are in `services/fnd/shared/authctx.go` — **NOT** in `internal/adapter/grpc/`. This is critical because `internal/adapter/grpc/server.go` imports `fndm001`, `fndm002`, etc., so those packages cannot import back into `internal/adapter/grpc/`.

```
✅ Correct:  fndm001 → shared (authctx.go)
✅ Correct:  internal/adapter/grpc → shared (authctx.go)
❌ WRONG:    fndm001 → internal/adapter/grpc (IMPORT CYCLE!)
```
