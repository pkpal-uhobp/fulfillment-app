package warehouses_repository_postgres

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

func (r *WarehousesRepository) PatchWarehouse(
	ctx context.Context,
	warehouseID int64,
	patch core_domain.WarehousePatch,
) (core_domain.Warehouse, error) {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()

	q := r.tx.Querier(ctx)

	var name string
	var warehouseType string
	var marketplace string
	var city string
	var address string
	var isActive bool

	if patch.Name != nil {
		name = *patch.Name
	}

	if patch.WarehouseType != nil {
		warehouseType = patch.WarehouseType.String()
	}

	if patch.Marketplace != nil {
		marketplace = *patch.Marketplace
	}

	if patch.City != nil {
		city = *patch.City
	}

	if patch.Address != nil {
		address = *patch.Address
	}

	if patch.IsActive != nil {
		isActive = *patch.IsActive
	}

	const query = `
		UPDATE warehouses
		SET name = CASE WHEN $2 THEN $3 ELSE name END,
			warehouse_type = CASE WHEN $4 THEN $5 ELSE warehouse_type END,
			marketplace = CASE WHEN $6 THEN $7 ELSE marketplace END,
			city = CASE WHEN $8 THEN $9 ELSE city END,
			address = CASE WHEN $10 THEN $11 ELSE address END,
			is_active = CASE WHEN $12 THEN $13 ELSE is_active END
		WHERE id = $1
		RETURNING id, name, warehouse_type, marketplace, city, address, is_active, created_at;
	`

	warehouse, err := scanWarehouseRow(q.QueryRow(
		ctx,
		query,
		warehouseID,

		patch.Name != nil,
		name,

		patch.WarehouseType != nil,
		warehouseType,

		patch.Marketplace != nil,
		marketplace,

		patch.City != nil,
		city,

		patch.Address != nil,
		address,

		patch.IsActive != nil,
		isActive,
	))
	if err != nil {
		return core_domain.Warehouse{}, fmt.Errorf("patch warehouse: %w", err)
	}

	return warehouse, nil
}
