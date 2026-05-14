//go:build integration

package cargoitems_repository_postgres_tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

func TestListCargoItemsByQRCodeForScan(t *testing.T) {
	repo, pool := newIntegrationRepository(t)
	fixture := seedCargoItemsFixture(t, pool)
	ctx := context.Background()
	qrCode := "QR-SCAN-" + uuid.NewString()
	created := createRepositoryCargoItem(t, repo, fixture, qrCode)

	items, err := repo.ListCargoItems(ctx, core_domain.CargoItemFilter{
		QRCode: qrCode,
		Page:   1,
		Limit:  1,
	})
	if err != nil {
		t.Fatalf("ListCargoItems by QR returned error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}
	if items[0].ID != created.ID {
		t.Fatalf("items[0].ID = %d, want %d", items[0].ID, created.ID)
	}

	items, err = repo.ListCargoItems(ctx, core_domain.CargoItemFilter{
		QRCode: "QR-UNKNOWN-" + uuid.NewString(),
		Page:   1,
		Limit:  1,
	})
	if err != nil {
		t.Fatalf("ListCargoItems by missing QR returned error: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("len(items) = %d, want 0", len(items))
	}
}
