# Users tests

The structure is split by feature layers:

```text
test/feature/users/
├── service/tests                 # UsersService unit tests with fake repository
├── transport/http/tests          # HTTP handlers unit tests with fake service
└── repository/postgres/tests     # PostgreSQL integration tests (integration build tag)
```

## Run all unit tests

From the repository root:

```bash
go test ./test/feature/users/... -v
```

## Run only service layer

```bash
go test ./test/feature/users/service/tests -v
```

## Run only HTTP layer

```bash
go test ./test/feature/users/transport/http/tests -v
```

## Coverage for users feature

PowerShell:

```powershell
$COVERPKG = (go list ./internal/features/users/... ./internal/core/domain) -join ","
go test ./test/feature/users/... -v -count=1 -coverpkg="$COVERPKG" -coverprofile=coverage_users.out
go tool cover -func="coverage_users.out"
go tool cover -html="coverage_users.out"
```

## Endpoints

```text
GET    /api/v1/users
POST   /api/v1/users
PATCH  /api/v1/users/{id}
PATCH  /api/v1/users/{id}/block
DELETE /api/v1/users/{id}
```

All users endpoints require `Admin` role (enforced by role middleware).
