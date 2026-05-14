# Warehouses tests by layers

Структура архива:

```text
test/feature/warehouses/
├── service/tests
├── transport/http/tests
└── repository/postgres/tests
```

Запуск всех unit/HTTP тестов:

```bash
go test ./test/feature/warehouses/... -v
```

Запуск только service-тестов:

```bash
go test ./test/feature/warehouses/service/tests -v
```

Запуск только HTTP-тестов:

```bash
go test ./test/feature/warehouses/transport/http/tests -v
```

Интеграционные PostgreSQL-тесты лежат под build tag `integration`, поэтому обычный `go test ./...` их не запускает:

```bash
$env:POSTGRES_HOST="localhost"
$env:POSTGRES_PORT="5432"
$env:POSTGRES_USER="postgres"
$env:POSTGRES_PASSWORD="root"
$env:POSTGRES_DB="fulfillment-app"
$env:POSTGRES_SSL_MODE="disable"
go test -tags=integration ./test/feature/warehouses/repository/postgres/tests -v
```

Для реальных PostgreSQL-тестов нужно поднять тестовую БД, применить миграции и передать `WAREHOUSES_TEST_DATABASE_URL`.
