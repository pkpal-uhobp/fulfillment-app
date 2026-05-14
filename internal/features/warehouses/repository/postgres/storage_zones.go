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

func (r *WarehousesRepository) ListStorageZones(
	ctx context.Context,
	warehouseID int64,
) ([]core_domain.StorageZone, error) {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()

	q := r.tx.Querier(ctx)

	const query = `
		SELECT id, warehouse_id, name, description, is_active
		FROM storage_zones
		WHERE ($1::bigint = 0 OR warehouse_id = $1)
		ORDER BY id;
	`

	rows, err := q.Query(ctx, query, warehouseID)
	if err != nil {
		return nil, fmt.Errorf("list storage zones: %w", err)
	}
	defer rows.Close()

	zones := make([]core_domain.StorageZone, 0)
	for rows.Next() {
		zone, err := scanStorageZoneRow(rows)
		if err != nil {
			return nil, err
		}

		zones = append(zones, zone)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate storage zones: %w", err)
	}

	return zones, nil
}

func (r *WarehousesRepository) CreateStorageZone(
	ctx context.Context,
	zone core_domain.StorageZone,
) (core_domain.StorageZone, error) {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()

	q := r.tx.Querier(ctx)

	const query = `
		INSERT INTO storage_zones (
			warehouse_id,
			name,
			description
		)
		VALUES ($1, $2, $3)
		RETURNING id, warehouse_id, name, description, is_active;
	`

	created, err := scanStorageZoneRow(q.QueryRow(
		ctx,
		query,
		zone.WarehouseID,
		zone.Name,
		zone.Description,
	))
	if err != nil {
		if isUniqueViolation(err) {
			return core_domain.StorageZone{}, fmt.Errorf(
				"%w: storage zone already exists",
				core_errors.ErrConflict,
			)
		}
		if isForeignKeyViolation(err) {
			return core_domain.StorageZone{}, fmt.Errorf(
				"%w: warehouse not found",
				core_errors.ErrNotFound,
			)
		}

		return core_domain.StorageZone{}, fmt.Errorf("create storage zone: %w", err)
	}

	return created, nil
}

func (r *WarehousesRepository) PatchStorageZone(
	ctx context.Context,
	zoneID int64,
	patch core_domain.StorageZonePatch,
) (core_domain.StorageZone, error) {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()

	q := r.tx.Querier(ctx)

	const query = `
		UPDATE storage_zones
		SET name = COALESCE($2, name),
			description = COALESCE($3, description),
			is_active = COALESCE($4, is_active)
		WHERE id = $1
		RETURNING id, warehouse_id, name, description, is_active;
	`

	zone, err := scanStorageZoneRow(q.QueryRow(
		ctx,
		query,
		zoneID,
		patch.Name,
		patch.Description,
		patch.IsActive,
	))
	if err != nil {
		if isUniqueViolation(err) {
			return core_domain.StorageZone{}, fmt.Errorf(
				"%w: storage zone already exists",
				core_errors.ErrConflict,
			)
		}

		return core_domain.StorageZone{}, fmt.Errorf("patch storage zone: %w", err)
	}

	return zone, nil
}

func scanStorageZoneRow(row pgx.Row) (core_domain.StorageZone, error) {
	var (
		zone        core_domain.StorageZone
		description sql.NullString
	)

	err := row.Scan(
		&zone.ID,
		&zone.WarehouseID,
		&zone.Name,
		&description,
		&zone.IsActive,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core_domain.StorageZone{}, fmt.Errorf(
				"%w: storage zone not found",
				core_errors.ErrNotFound,
			)
		}

		return core_domain.StorageZone{}, fmt.Errorf("scan storage zone: %w", err)
	}

	if description.Valid {
		zone.Description = &description.String
	}

	return zone, nil
}
