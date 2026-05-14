include .env
export

export PROJECT_ROOT=${CURDIR}

env-up:
	docker compose up -d fulfillment-app-postgres
	@powershell -NoProfile -ExecutionPolicy Bypass -Command "for ($$i = 0; $$i -lt 30; $$i++) { docker exec fulfillment-env-postgres pg_isready -U $(POSTGRES_USER) -d $(POSTGRES_DB); if ($$LASTEXITCODE -eq 0) { Write-Host 'postgres is ready'; exit 0 }; Start-Sleep -Seconds 1 }; Write-Host 'postgres is not ready'; exit 1"

env-down:
	docker compose down fulfillment-app-postgres

env-cleanup:
	@powershell -NoProfile -ExecutionPolicy Bypass -Command "$$ans = Read-Host 'Clean all volume? Warning loss data. [y/N]'; if ($$ans -eq 'y') { docker compose down; if (Test-Path 'out/pgdata') { Remove-Item -Recurse -Force 'out/pgdata' }; Write-Host 'done' } else { Write-Host 'cancel operation' }"

env-port-forward:
	docker compose up -d port-forwarder

env-port-close:
	docker compose down port-forwarder

migrate-create:
	@if "$(seq)"=="" (echo not seq && exit /b 1)
	docker compose run --rm fulfillment-app-migrate create -ext sql -dir /migrations -seq "$(seq)"

migrate-up:
	make migrate-action action=up

migrate-down:
	make migrate-action action=down

migrate-action:
	@if "$(action)"=="" (echo not action && exit /b 1)
	docker compose run --rm fulfillment-app-migrate \
    		-path /migrations \
    		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@fulfillment-app-postgres:5432/${POSTGRES_DB}?sslmode=disable \
    		"${action}"

fulfillment-app-run:
	set "LOGGER_FOLDER=%CD%\out\logs" && set "POSTGRES_HOST=127.0.0.1" && go mod tidy && go run .\cmd\fulfillment\main.go