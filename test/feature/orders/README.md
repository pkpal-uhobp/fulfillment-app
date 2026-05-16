# Orders tests

Структура разделена по слоям фичи `orders`:

```text
test/feature/orders/
├── service/tests                 # unit-тесты OrdersService через fake repository
├── transport/http/tests          # unit-тесты HTTP handlers через fake service
└── repository/postgres/tests     # integration-тесты PostgreSQL под build tag integration
```

## Запуск обычных unit-тестов

Из корня репозитория:

```bash
go test ./test/feature/orders/... -v
```

## Запуск только service-слоя

```bash
go test ./test/feature/orders/service/tests -v
```

## Запуск только HTTP-слоя

```bash
go test ./test/feature/orders/transport/http/tests -v
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
go test -tags=integration ./test/feature/orders/repository/postgres/tests -v
```

Если переменные PostgreSQL не заданы, integration-тесты будут пропущены.

## Проверка покрытия service-слоя

Так как тесты лежат отдельно от исходного кода, для корректного процента покрытия используй `-coverpkg`:

```bash
go test ./test/feature/orders/service/tests \
  -v \
  -count=1 \
  -coverpkg=./internal/features/orders/service \
  -coverprofile=coverage_orders_service.out

go tool cover -func=coverage_orders_service.out
go tool cover -html=coverage_orders_service.out
```

## Проверка покрытия всей фичи orders

```bash
go test ./test/feature/orders/... \
  -v \
  -count=1 \
  -coverpkg=./internal/features/orders/... \
  -coverprofile=coverage_orders.out

go tool cover -func=coverage_orders.out
go tool cover -html=coverage_orders.out
```
