//go:build integration

package warehouses_repository_postgres_tests

import (
	"os"
	"testing"
)

func TestWarehousesRepositoryIntegrationSetup(t *testing.T) {
	if os.Getenv("WAREHOUSES_TEST_DATABASE_URL") == "" {
		t.Skip("set WAREHOUSES_TEST_DATABASE_URL and run migrations before enabling warehouses repository integration tests")
	}

	t.Skip("add real PostgreSQL repository tests here when a dedicated test database is configured")
}
