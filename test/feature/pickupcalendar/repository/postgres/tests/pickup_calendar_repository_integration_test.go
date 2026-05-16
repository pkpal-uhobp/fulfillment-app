//go:build integration

package pickupcalendar_repository_postgres_tests

import "testing"

func TestPickupCalendarRepositoryIntegrationRequiresDatabase(t *testing.T) {
	t.Skip("integration skeleton: поднять PostgreSQL, применить migrations/000011_pickup_calendar и добавить реальные repository-тесты")
}
