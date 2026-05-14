package warehouses_repository_postgres

import (
	"context"
	"fmt"

	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (r *WarehousesRepository) ActivateGate(ctx context.Context, gateID int64) error {
	return r.setGateActive(ctx, gateID, true)
}

func (r *WarehousesRepository) DeactivateGate(ctx context.Context, gateID int64) error {
	return r.setGateActive(ctx, gateID, false)
}

func (r *WarehousesRepository) setGateActive(ctx context.Context, gateID int64, isActive bool) error {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()

	q := r.tx.Querier(ctx)

	const query = `
		UPDATE gates
		SET is_active = $2
		WHERE id = $1;
	`

	tag, err := q.Exec(ctx, query, gateID, isActive)
	if err != nil {
		return fmt.Errorf("set gate active: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%w: gate not found", core_errors.ErrNotFound)
	}

	return nil
}
