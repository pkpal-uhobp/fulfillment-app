# Shipments tests

The structure is split by feature layers:

```text
test/feature/shipments/
├── service/tests                 # ShipmentsService unit tests with fake repository
├── transport/http/tests          # HTTP handlers unit tests with fake service
└── repository/postgres/tests     # PostgreSQL integration tests (integration build tag)
```

## Run all unit tests

From the repository root:

```bash
go test ./test/feature/shipments/... -v
```

## Run only service layer

```bash
go test ./test/feature/shipments/service/tests -v
```

## Run only HTTP layer

```bash
go test ./test/feature/shipments/transport/http/tests -v
```

## Coverage

```powershell
$COVERPKG = (go list ./internal/features/shipments/... ./internal/core/domain) -join ","
go test ./test/feature/shipments/... -v -count=1 -coverpkg="$COVERPKG" -coverprofile="coverage_shipments.out"
go tool cover -func="coverage_shipments.out"
go tool cover -html="coverage_shipments.out"
```

## PostgreSQL integration tests

PostgreSQL must be running and migrations applied.

```powershell
$env:POSTGRES_HOST="localhost"
$env:POSTGRES_PORT="5433"
$env:POSTGRES_USER="postgres"
$env:POSTGRES_PASSWORD="root"
$env:POSTGRES_DB="fulfillment-app"
$env:POSTGRES_SSL_MODE="disable"
go test -tags=integration ./test/feature/shipments/repository/postgres/tests -v
```
