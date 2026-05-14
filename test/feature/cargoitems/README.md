# Cargo items feature tests

Тесты новой фичи `cargoitems` разделены по тем же слоям, что и у остальных фич проекта:

```text
cargoitems/
├── repository/
├── service/
├── transport/
└── README.md
```

## Что проверяется

### repository/postgres
Интеграционные тесты репозитория. Проверяют реальные SQL-запросы к PostgreSQL:

- создание `cargo_items`;
- запись в `cargo_status_history`;
- получение одного грузового места;
- получение списка с фильтрами;
- проверку принадлежности грузового места клиенту;
- проверку принадлежности зоны хранения и гейта складу заявки;
- назначение зоны хранения;
- назначение гейта;
- изменение статуса;
- конфликт при повторном `qr_code`.

Запуск:

```bash
go test -tags=integration ./test/feature/cargoitems/repository/postgres/...
```

Перед запуском нужна тестовая БД PostgreSQL с применёнными миграциями и переменные окружения:

```bash
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=fulfillment_test
POSTGRES_SSL_MODE=disable
```

### service
Unit-тесты service-слоя через fake repository. Проверяют бизнес-правила:

- роли пользователей;
- генерацию QR;
- запрет принять больше грузовых мест, чем заявлено;
- запрет работы с отменённой заявкой;
- переходы статусов;
- обязательность зоны перед `stored`;
- обязательность гейта перед `ready_to_ship`/`shipped`;
- ограничения клиента на просмотр только своих грузовых мест.

Запуск:

```bash
go test ./test/feature/cargoitems/service/...
```

### transport/http
Unit-тесты HTTP-слоя через fake service и `httptest`.
Проверяют routes и основные handlers.

Запуск:

```bash
go test ./test/feature/cargoitems/transport/http/tests -v
```

## Запуск всех unit-тестов cargoitems

```bash
go test ./test/feature/cargoitems/service/tests -v
```

## Запуск всех тестов cargoitems вместе с repository integration

```bash
$env:POSTGRES_HOST="localhost"
$env:POSTGRES_PORT="5433"
$env:POSTGRES_USER="postgres"
$env:POSTGRES_PASSWORD="root"
$env:POSTGRES_DB="fulfillment-app"
$env:POSTGRES_SSL_MODE="disable"
go test -tags=integration ./test/feature/cargoitems/repository/postgres/tests -v
```
