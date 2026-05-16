# Shipments tests

Структура разделена по слоям фичи `shipments`:

```text
test/feature/shipments/
├── service/tests                 # unit-тесты ShipmentsService через fake repository
├── transport/http/tests          # unit-тесты HTTP handlers через fake service
└── repository/postgres/tests     # integration-тесты PostgreSQL под build tag integration
```

## Запуск обычных unit-тестов

Из корня репозитория:

```bash
go test ./test/feature/shipments/... -v
```

## Запуск только service-слоя

```bash
go test ./test/feature/shipments/service/tests -v
```

## Запуск только HTTP-слоя

```bash
go test ./test/feature/shipments/transport/http/tests -v
```

## Проверка покрытия

```powershell
$COVERPKG = (go list ./internal/features/shipments/... ./internal/core/domain) -join ","
go test ./test/feature/shipments/... -v -count=1 -coverpkg="$COVERPKG" -coverprofile="coverage_shipments.out"
go tool cover -func="coverage_shipments.out"
go tool cover -html="coverage_shipments.out"
```

## PostgreSQL integration-тесты

Перед запуском должна быть поднята БД и применены миграции.

```powershell
$env:POSTGRES_HOST="localhost"
$env:POSTGRES_PORT="5432"
$env:POSTGRES_USER="postgres"
$env:POSTGRES_PASSWORD="root"
$env:POSTGRES_DB="fulfillment-app"
$env:POSTGRES_SSL_MODE="disable"
go test -tags=integration ./test/feature/shipments/repository/postgres/tests -v
```
