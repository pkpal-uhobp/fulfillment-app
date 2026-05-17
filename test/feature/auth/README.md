# Auth tests

The structure is split by feature layers:

```text
test/feature/auth/
├── service/tests                 # AuthService unit tests with fake repository
├── transport/http/tests          # HTTP handlers unit tests with fake service
└── repository/postgres/tests     # PostgreSQL integration tests (integration build tag)
```

## Run all unit tests

From the repository root:

```bash
go test ./test/feature/auth/... -v
```

## Run only service layer

```bash
go test ./test/feature/auth/service/tests -v
```

## Run only HTTP layer

```bash
go test ./test/feature/auth/transport/http/tests -v
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
go test -tags=integration ./test/feature/auth/repository/postgres/tests -v
```

If PostgreSQL variables are not set, integration tests are skipped.
