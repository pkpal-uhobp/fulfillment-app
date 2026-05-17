include .env
export

export PROJECT_ROOT=${CURDIR}

.PHONY: env-up env-down env-cleanup env-port-forward env-port-close \
        migrate-create migrate-up migrate-down migrate-action \
        fulfillment-app-run backend-up backend-down \
        frontend-install frontend-run frontend-dev-up frontend-dev-recreate frontend-down \
        prod-up prod-down logs logs-frontend test coverage seed

env-up:
	docker compose up -d fulfillment-app-postgres
	@powershell -NoProfile -ExecutionPolicy Bypass -Command "for ($$i = 0; $$i -lt 30; $$i++) { docker exec fulfillment-env-postgres pg_isready -U $(POSTGRES_USER) -d $(POSTGRES_DB); if ($$LASTEXITCODE -eq 0) { Write-Host 'postgres is ready'; exit 0 }; Start-Sleep -Seconds 1 }; Write-Host 'postgres is not ready'; exit 1"

env-down:
	docker compose down

env-cleanup:
	@powershell -NoProfile -ExecutionPolicy Bypass -Command "$$ans = Read-Host 'Clean all volume? Warning loss data. [y/N]'; if ($$ans -eq 'y') { docker compose down; if (Test-Path 'out/pgdata') { Remove-Item -Recurse -Force 'out/pgdata' }; Write-Host 'done' } else { Write-Host 'cancel operation' }"

env-port-forward:
	docker compose --profile dev up -d port-forwarder

env-port-close:
	docker compose stop port-forwarder

migrate-create:
	@if "$(seq)"=="" (echo not seq && exit /b 1)
	docker compose --profile tools run --rm fulfillment-app-migrate create -ext sql -dir /migrations -seq "$(seq)"

migrate-up:
	make migrate-action action=up

migrate-down:
	make migrate-action action=down

migrate-action:
	@if "$(action)"=="" (echo not action && exit /b 1)
	docker compose --profile tools run --rm fulfillment-app-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@fulfillment-app-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

seed:
	docker cp .\scripts\seed_test_data.sql fulfillment-env-postgres:/tmp/seed_test_data.sql
	docker exec -i fulfillment-env-postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -v ON_ERROR_STOP=1 -f /tmp/seed_test_data.sql

fulfillment-app-run:
	set "LOGGER_FOLDER=%CD%\out\logs" && set "POSTGRES_HOST=127.0.0.1" && set "POSTGRES_PORT=5433" && go mod tidy && go run .\cmd\fulfillment\main.go

backend-up:
	docker compose --profile app up -d --build fulfillment-app-backend

backend-down:
	docker compose stop fulfillment-app-backend

frontend-install:
	cd frontend && npm install

frontend-run:
	cd frontend && npm run dev

frontend-dev-up:
	docker compose --profile frontend-dev up -d --force-recreate fulfillment-app-frontend-dev

frontend-dev-recreate:
	docker compose --profile frontend-dev up -d --force-recreate fulfillment-app-frontend-dev

frontend-down:
	docker compose stop fulfillment-app-frontend-dev fulfillment-app-frontend

prod-up:
	docker compose --profile prod up -d --build

prod-down:
	docker compose --profile prod down

logs:
	docker compose logs -f

logs-frontend:
	docker compose logs -f fulfillment-app-frontend-dev

test:
	go test ./test/feature/... -v -count=1

coverage:
	@powershell -NoProfile -ExecutionPolicy Bypass -Command "$$COVERPKG = (go list ./internal/features/...) -join ','; go test ./test/feature/... -v -count=1 -coverpkg=$$COVERPKG -coverprofile='coverage_all.out'; go tool cover -func='coverage_all.out'"
