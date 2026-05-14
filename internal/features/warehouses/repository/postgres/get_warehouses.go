package warehouses_repository_postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (r *WarehousesRepository) ListWarehouses(
	ctx context.Context,
	filter core_domain.WarehouseFilter,
) ([]core_domain.Warehouse, error) {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()

	q := r.tx.Querier(ctx)

	const query = `
		SELECT id, name, warehouse_type, marketplace, city, address, is_active, created_at
		FROM warehouses
		WHERE ($1 = '' OR warehouse_type = $1)
		  AND ($2 = '' OR lower(coalesce(marketplace, '')) = lower($2))
		  AND ($3 = '' OR lower(city) = lower($3))
		ORDER BY id;
	`

	rows, err := q.Query(
		ctx,
		query,
		filter.WarehouseType,
		filter.Marketplace,
		filter.City,
	)
	if err != nil {
		return nil, fmt.Errorf("list warehouses: %w", err)
	}
	defer rows.Close()

	warehouses := make([]core_domain.Warehouse, 0)
	for rows.Next() {
		warehouse, err := scanWarehouseRow(rows)
		if err != nil {
			return nil, err
		}

		warehouses = append(warehouses, warehouse)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate warehouses: %w", err)
	}

	return warehouses, nil
}

func scanWarehouseRow(row pgx.Row) (core_domain.Warehouse, error) {
	var (
		warehouse     core_domain.Warehouse
		warehouseType string
		marketplace   sql.NullString
	)

	err := row.Scan(
		&warehouse.ID,
		&warehouse.Name,
		&warehouseType,
		&marketplace,
		&warehouse.City,
		&warehouse.Address,
		&warehouse.IsActive,
		&warehouse.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core_domain.Warehouse{}, fmt.Errorf(
				"%w: warehouse not found",
				core_errors.ErrNotFound,
			)
		}

		return core_domain.Warehouse{}, fmt.Errorf("scan warehouse: %w", err)
	}

	warehouse.WarehouseType = core_domain.WarehouseType(warehouseType)
	if marketplace.Valid {
		warehouse.Marketplace = &marketplace.String
	}

	return warehouse, nil
}
