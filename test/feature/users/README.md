# Users tests

Структура разделена по слоям фичи `users`:

```text
test/feature/users/
├── service/tests                 # unit-тесты UsersService через fake repository
├── transport/http/tests          # unit-тесты HTTP handlers через fake service
└── repository/postgres/tests     # integration-тесты PostgreSQL под build tag integration
```

## Запуск обычных unit-тестов

Из корня репозитория:

```bash
go test ./test/feature/users/... -v
```

## Запуск только service-слоя

```bash
go test ./test/feature/users/service/tests -v
```

## Запуск только HTTP-слоя

```bash
go test ./test/feature/users/transport/http/tests -v
```

## Покрытие фичи users

PowerShell:

```powershell
$COVERPKG = (go list ./internal/features/users/... ./internal/core/domain) -join ","
go test ./test/feature/users/... -v -count=1 -coverpkg="$COVERPKG" -coverprofile=coverage_users.out
go tool cover -func="coverage_users.out"
go tool cover -html="coverage_users.out"
```

## Endpoints

```text
GET    /api/v1/users
POST   /api/v1/users
PATCH  /api/v1/users/{id}
PATCH  /api/v1/users/{id}/block
DELETE /api/v1/users/{id}
```

Роль не указана в endpoint-е. Доступ ограничивается через role middleware: все маршруты доступны только `admin`.
