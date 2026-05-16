# Pickup calendar tests

Структура разделена по слоям фичи `pickupcalendar`:

```text
test/feature/pickupcalendar/
├── service/tests                 # unit-тесты PickupCalendarService через fake repository
├── transport/http/tests          # unit-тесты HTTP handlers через fake service
└── repository/postgres/tests     # integration-тесты PostgreSQL под build tag integration
```

## Запуск обычных unit-тестов

Из корня репозитория:

```bash
go test ./test/feature/pickupcalendar/... -v
```

## Запуск только service-слоя

```bash
go test ./test/feature/pickupcalendar/service/tests -v
```

## Запуск только HTTP-слоя

```bash
go test ./test/feature/pickupcalendar/transport/http/tests -v
```

## Покрытие pickup-calendar

PowerShell:

```powershell
$COVERPKG = (go list ./internal/features/pickupcalendar/... ./internal/core/domain) -join ","
go test ./test/feature/pickupcalendar/... -v -count=1 -coverpkg="$COVERPKG" -coverprofile=coverage_pickupcalendar.out
go tool cover -func coverage_pickupcalendar.out
go tool cover -html coverage_pickupcalendar.out
```

## Запуск PostgreSQL integration-тестов

Перед запуском должна быть поднята БД и применены миграции.

```bash
$env:POSTGRES_HOST="localhost"
$env:POSTGRES_PORT="5432"
$env:POSTGRES_USER="postgres"
$env:POSTGRES_PASSWORD="root"
$env:POSTGRES_DB="fulfillment-app"
$env:POSTGRES_SSL_MODE="disable"
go test -tags=integration ./test/feature/pickupcalendar/repository/postgres/tests -v
```

Если переменные PostgreSQL не заданы, integration-тесты будут пропущены.
