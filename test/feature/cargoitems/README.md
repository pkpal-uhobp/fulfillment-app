# Cargo items feature tests

Tests are split by the same layers as other features:

```text
test/feature/cargoitems/
├── repository/
├── service/
├── transport/
└── README.md
```

## What is covered

### repository/postgres
Integration tests for PostgreSQL queries:

- create `cargo_items`
- insert into `cargo_status_history`
- fetch single cargo item
- fetch list with filters
- validate cargo item ownership by client
- validate storage zone and gate belong to order warehouse
- assign storage zone
- assign gate
- update status
- duplicate `qr_code` conflict

Run:

```bash
go test -tags=integration ./test/feature/cargoitems/repository/postgres/...
```

PostgreSQL must be running and migrations applied. Use these env vars:

```powershell
$env:POSTGRES_HOST="localhost"
$env:POSTGRES_PORT="5433"
$env:POSTGRES_USER="postgres"
$env:POSTGRES_PASSWORD="root"
$env:POSTGRES_DB="fulfillment-app"
$env:POSTGRES_SSL_MODE="disable"
```

### service
Unit tests for the service layer with a fake repository. Covers business rules:

- user roles
- QR generation
- cannot accept more cargo items than ordered
- cannot operate on cancelled orders
- status transitions
- storage zone required before `stored`
- gate required before `ready_to_ship`/`shipped`
- client access only to their own items

Run:

```bash
go test ./test/feature/cargoitems/service/...
```

### transport/http
Unit tests for HTTP handlers with a fake service and `httptest`.

Run:

```bash
go test ./test/feature/cargoitems/transport/http/tests -v
```

## Run all cargoitems unit tests

```bash
go test ./test/feature/cargoitems/service/tests -v
```

## Run all cargoitems tests including repository integration

```powershell
$env:POSTGRES_HOST="localhost"
$env:POSTGRES_PORT="5433"
$env:POSTGRES_USER="postgres"
$env:POSTGRES_PASSWORD="root"
$env:POSTGRES_DB="fulfillment-app"
$env:POSTGRES_SSL_MODE="disable"
go test -tags=integration ./test/feature/cargoitems/repository/postgres/tests -v
```
