# Auth tests

Структура разделена по слоям фичи `auth`:

```text
test/feature/auth/
├── service/tests                 # unit-тесты AuthService через fake repository
├── transport/http/tests          # unit-тесты HTTP handlers через fake service
└── repository/postgres/tests     # integration-тесты PostgreSQL под build tag integration
```

## Запуск обычных unit-тестов

Из корня репозитория:

```bash
go test ./test/feature/auth/... -v
```

## Запуск только service-слоя

```bash
go test ./test/feature/auth/service/tests -v
```

## Запуск только HTTP-слоя

```bash
go test ./test/feature/auth/transport/http/tests -v
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
go test -tags=integration ./test/feature/auth/repository/postgres/tests -v
```

Если переменные PostgreSQL не заданы, integration-тесты будут пропущены.
