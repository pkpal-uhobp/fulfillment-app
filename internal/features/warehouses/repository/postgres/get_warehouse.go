package warehouses_repository_postgres

import (
	"context"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

func (r *WarehousesRepository) GetWarehouseByID(
	ctx context.Context,
	warehouseID int64,
) (core_domain.Warehouse, error) {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()

	q := r.tx.Querier(ctx)

	const query = `
		SELECT id, name, warehouse_type, marketplace, city, address, is_active, created_at
		FROM warehouses
		WHERE id = $1;
	`

	return scanWarehouseRow(q.QueryRow(ctx, query, warehouseID))
}
