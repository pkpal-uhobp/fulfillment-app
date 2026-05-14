package warehouses_repository_postgres

import (
	"context"
	"fmt"

	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (r *WarehousesRepository) ActivateStorageZone(ctx context.Context, zoneID int64) error {
	return r.setStorageZoneActive(ctx, zoneID, true)
}

func (r *WarehousesRepository) DeactivateStorageZone(ctx context.Context, zoneID int64) error {
	return r.setStorageZoneActive(ctx, zoneID, false)
}

func (r *WarehousesRepository) setStorageZoneActive(ctx context.Context, zoneID int64, isActive bool) error {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()

	q := r.tx.Querier(ctx)

	const query = `
		UPDATE storage_zones
		SET is_active = $2
		WHERE id = $1;
	`

	tag, err := q.Exec(ctx, query, zoneID, isActive)
	if err != nil {
		return fmt.Errorf("set storage zone active: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%w: storage zone not found", core_errors.ErrNotFound)
	}

	return nil
}
