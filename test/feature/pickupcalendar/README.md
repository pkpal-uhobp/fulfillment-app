# Pickup calendar tests

The structure is split by feature layers:

```text
test/feature/pickupcalendar/
├── service/tests                 # PickupCalendarService unit tests with fake repository
├── transport/http/tests          # HTTP handlers unit tests with fake service
└── repository/postgres/tests     # PostgreSQL integration tests (integration build tag)
```

## Run all unit tests

From the repository root:

```bash
go test ./test/feature/pickupcalendar/... -v
```

## Run only service layer

```bash
go test ./test/feature/pickupcalendar/service/tests -v
```

## Run only HTTP layer

```bash
go test ./test/feature/pickupcalendar/transport/http/tests -v
```

## Coverage for pickup calendar

PowerShell:

```powershell
$COVERPKG = (go list ./internal/features/pickupcalendar/... ./internal/core/domain) -join ","
go test ./test/feature/pickupcalendar/... -v -count=1 -coverpkg="$COVERPKG" -coverprofile=coverage_pickupcalendar.out
go tool cover -func coverage_pickupcalendar.out
go tool cover -html coverage_pickupcalendar.out
```

## Run PostgreSQL integration tests

PostgreSQL must be running and migrations applied.

```powershell
$env:POSTGRES_HOST="localhost"
$env:POSTGRES_PORT="5433"
$env:POSTGRES_USER="postgres"
$env:POSTGRES_PASSWORD="root"
$env:POSTGRES_DB="fulfillment-app"
$env:POSTGRES_SSL_MODE="disable"
go test -tags=integration ./test/feature/pickupcalendar/repository/postgres/tests -v
```

If PostgreSQL variables are not set, integration tests are skipped.
