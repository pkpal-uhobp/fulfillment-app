package warehouses_repository_postgres

import (
	"context"
	"fmt"

	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (r *WarehousesRepository) DeactivateWarehouse(
	ctx context.Context,
	warehouseID int64,
) error {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()

	q := r.tx.Querier(ctx)

	const deactivateWarehouseQuery = `
		UPDATE warehouses
		SET is_active = false
		WHERE id = $1;
	`

	tag, err := q.Exec(ctx, deactivateWarehouseQuery, warehouseID)
	if err != nil {
		return fmt.Errorf("deactivate warehouse: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%w: warehouse not found", core_errors.ErrNotFound)
	}

	const deactivateStorageZonesQuery = `
		UPDATE storage_zones
		SET is_active = false
		WHERE warehouse_id = $1;
	`

	if _, err := q.Exec(ctx, deactivateStorageZonesQuery, warehouseID); err != nil {
		return fmt.Errorf("deactivate warehouse storage zones: %w", err)
	}

	const deactivateGatesQuery = `
		UPDATE gates
		SET is_active = false
		WHERE warehouse_id = $1;
	`

	if _, err := q.Exec(ctx, deactivateGatesQuery, warehouseID); err != nil {
		return fmt.Errorf("deactivate warehouse gates: %w", err)
	}

	return nil
}
