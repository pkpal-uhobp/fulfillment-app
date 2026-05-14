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

	var warehouseType *string
	if patch.WarehouseType != nil {
		value := patch.WarehouseType.String()
		warehouseType = &value
	}

	const query = `
		UPDATE warehouses
		SET name = COALESCE($2, name),
			warehouse_type = COALESCE($3, warehouse_type),
			marketplace = COALESCE($4, marketplace),
			city = COALESCE($5, city),
			address = COALESCE($6, address),
			is_active = COALESCE($7, is_active)
		WHERE id = $1
		RETURNING id, name, warehouse_type, marketplace, city, address, is_active, created_at;
	`

	warehouse, err := scanWarehouseRow(q.QueryRow(
		ctx,
		query,
		warehouseID,
		patch.Name,
		warehouseType,
		patch.Marketplace,
		patch.City,
		patch.Address,
		patch.IsActive,
	))
	if err != nil {
		return core_domain.Warehouse{}, fmt.Errorf("patch warehouse: %w", err)
	}

	return warehouse, nil
}
