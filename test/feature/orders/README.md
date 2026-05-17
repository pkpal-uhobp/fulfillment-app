# Orders tests

The structure is split by feature layers:

```text
test/feature/orders/
├── service/tests                 # OrdersService unit tests with fake repository
├── transport/http/tests          # HTTP handlers unit tests with fake service
└── repository/postgres/tests     # PostgreSQL integration tests (integration build tag)
```

## Run all unit tests

From the repository root:

```bash
go test ./test/feature/orders/... -v
```

## Run only service layer

```bash
go test ./test/feature/orders/service/tests -v
```

## Run only HTTP layer

```bash
go test ./test/feature/orders/transport/http/tests -v
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
go test -tags=integration ./test/feature/orders/repository/postgres/tests -v
```

If PostgreSQL variables are not set, integration tests are skipped.

## Service layer coverage

Because tests are outside the source tree, use `-coverpkg` for accurate coverage:

```bash
go test ./test/feature/orders/service/tests \
  -v \
  -count=1 \
  -coverpkg=./internal/features/orders/service \
  -coverprofile=coverage_orders_service.out

go tool cover -func=coverage_orders_service.out
go tool cover -html=coverage_orders_service.out
```

## Full orders feature coverage

```bash
go test ./test/feature/orders/... \
  -v \
  -count=1 \
  -coverpkg=./internal/features/orders/... \
  -coverprofile=coverage_orders.out

go tool cover -func=coverage_orders.out
go tool cover -html=coverage_orders.out
```
