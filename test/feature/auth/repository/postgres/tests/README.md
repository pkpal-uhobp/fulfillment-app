# Repository/Postgres tests

Файл `auth_repository_integration_test.go` помечен build tag `integration`:

```go
//go:build integration
```

Поэтому он не запускается обычной командой `go test ./...` и не мешает unit-тестам.
Для запуска нужна тестовая PostgreSQL с примененными миграциями `000001_users` и `000002_issued_tokens` минимум.
