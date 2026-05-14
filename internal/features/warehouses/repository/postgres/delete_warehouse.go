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

	const query = `
		UPDATE warehouses
		SET is_active = false
		WHERE id = $1;
	`

	tag, err := q.Exec(ctx, query, warehouseID)
	if err != nil {
		return fmt.Errorf("deactivate warehouse: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%w: warehouse not found", core_errors.ErrNotFound)
	}

	return nil
}
