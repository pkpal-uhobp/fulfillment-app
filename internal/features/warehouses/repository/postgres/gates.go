package warehouses_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (r *WarehousesRepository) ListGates(
	ctx context.Context,
	warehouseID int64,
) ([]core_domain.Gate, error) {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()

	q := r.tx.Querier(ctx)

	const query = `
		SELECT id, warehouse_id, name, is_active
		FROM gates
		WHERE ($1::bigint = 0 OR warehouse_id = $1)
		ORDER BY id;
	`

	rows, err := q.Query(ctx, query, warehouseID)
	if err != nil {
		return nil, fmt.Errorf("list gates: %w", err)
	}
	defer rows.Close()

	gates := make([]core_domain.Gate, 0)
	for rows.Next() {
		gate, err := scanGateRow(rows)
		if err != nil {
			return nil, err
		}

		gates = append(gates, gate)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate gates: %w", err)
	}

	return gates, nil
}

func (r *WarehousesRepository) CreateGate(
	ctx context.Context,
	gate core_domain.Gate,
) (core_domain.Gate, error) {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()

	q := r.tx.Querier(ctx)

	const query = `
		INSERT INTO gates (
			warehouse_id,
			name
		)
		VALUES ($1, $2)
		RETURNING id, warehouse_id, name, is_active;
	`

	created, err := scanGateRow(q.QueryRow(
		ctx,
		query,
		gate.WarehouseID,
		gate.Name,
	))
	if err != nil {
		if isUniqueViolation(err) {
			return core_domain.Gate{}, fmt.Errorf(
				"%w: gate already exists",
				core_errors.ErrConflict,
			)
		}
		if isForeignKeyViolation(err) {
			return core_domain.Gate{}, fmt.Errorf(
				"%w: warehouse not found",
				core_errors.ErrNotFound,
			)
		}

		return core_domain.Gate{}, fmt.Errorf("create gate: %w", err)
	}

	return created, nil
}

func (r *WarehousesRepository) PatchGate(
	ctx context.Context,
	gateID int64,
	patch core_domain.GatePatch,
) (core_domain.Gate, error) {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()

	q := r.tx.Querier(ctx)

	const query = `
		UPDATE gates
		SET name = COALESCE($2, name),
			is_active = COALESCE($3, is_active)
		WHERE id = $1
		RETURNING id, warehouse_id, name, is_active;
	`

	gate, err := scanGateRow(q.QueryRow(
		ctx,
		query,
		gateID,
		patch.Name,
		patch.IsActive,
	))
	if err != nil {
		if isUniqueViolation(err) {
			return core_domain.Gate{}, fmt.Errorf(
				"%w: gate already exists",
				core_errors.ErrConflict,
			)
		}

		return core_domain.Gate{}, fmt.Errorf("patch gate: %w", err)
	}

	return gate, nil
}

func scanGateRow(row pgx.Row) (core_domain.Gate, error) {
	var gate core_domain.Gate

	err := row.Scan(
		&gate.ID,
		&gate.WarehouseID,
		&gate.Name,
		&gate.IsActive,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core_domain.Gate{}, fmt.Errorf(
				"%w: gate not found",
				core_errors.ErrNotFound,
			)
		}

		return core_domain.Gate{}, fmt.Errorf("scan gate: %w", err)
	}

	return gate, nil
}
