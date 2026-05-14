package warehouses_repository_postgres

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

func (r *WarehousesRepository) CreateWarehouse(
	ctx context.Context,
	warehouse core_domain.Warehouse,
) (core_domain.Warehouse, error) {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()

	q := r.tx.Querier(ctx)

	const query = `
		INSERT INTO warehouses (
			name,
			warehouse_type,
			marketplace,
			city,
			address
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, warehouse_type, marketplace, city, address, is_active, created_at;
	`

	created, err := scanWarehouseRow(q.QueryRow(
		ctx,
		query,
		warehouse.Name,
		warehouse.WarehouseType.String(),
		warehouse.Marketplace,
		warehouse.City,
		warehouse.Address,
	))
	if err != nil {
		return core_domain.Warehouse{}, fmt.Errorf("create warehouse: %w", err)
	}

	return created, nil
}
