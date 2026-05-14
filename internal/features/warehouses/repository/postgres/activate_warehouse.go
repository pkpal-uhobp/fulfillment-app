package warehouses_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (r *WarehousesRepository) ActivateWarehouse(ctx context.Context, warehouseID int64) error {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()

	q := r.tx.Querier(ctx)

	const query = `
		WITH updated_warehouse AS (
			UPDATE warehouses
			SET is_active = true
			WHERE id = $1
			RETURNING id
		),
		updated_storage_zones AS (
			UPDATE storage_zones
			SET is_active = true
			WHERE warehouse_id IN (SELECT id FROM updated_warehouse)
			RETURNING id
		),
		updated_gates AS (
			UPDATE gates
			SET is_active = true
			WHERE warehouse_id IN (SELECT id FROM updated_warehouse)
			RETURNING id
		)
		SELECT id FROM updated_warehouse;
	`

	var id int64
	if err := q.QueryRow(ctx, query, warehouseID).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%w: warehouse not found", core_errors.ErrNotFound)
		}

		return fmt.Errorf("activate warehouse: %w", err)
	}

	return nil
}
