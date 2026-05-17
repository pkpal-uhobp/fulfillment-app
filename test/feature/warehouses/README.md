# Warehouses tests

Structure:

```text
test/feature/warehouses/
├── service/tests
├── transport/http/tests
└── repository/postgres/tests
```

## Run all unit/HTTP tests

```bash
go test ./test/feature/warehouses/... -v
```

## Run only service tests

```bash
go test ./test/feature/warehouses/service/tests -v
```

## Run only HTTP tests

```bash
go test ./test/feature/warehouses/transport/http/tests -v
```

## PostgreSQL integration tests

Integration tests are under the `integration` build tag and are not run by `go test ./...`.

```powershell
$env:POSTGRES_HOST="localhost"
$env:POSTGRES_PORT="5433"
$env:POSTGRES_USER="postgres"
$env:POSTGRES_PASSWORD="root"
$env:POSTGRES_DB="fulfillment-app"
$env:POSTGRES_SSL_MODE="disable"
go test -tags=integration ./test/feature/warehouses/repository/postgres/tests -v
```

Note: repository integration tests are currently a placeholder. To enable real tests, set `WAREHOUSES_TEST_DATABASE_URL` and apply migrations for a dedicated test database.
