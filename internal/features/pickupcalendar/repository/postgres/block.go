package pickupcalendar_repository_postgres

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (r *PickupCalendarRepository) CreateBlock(
	ctx context.Context,
	block core_domain.PickupCalendarBlock,
) (core_domain.PickupCalendarBlock, error) {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()
	q := r.tx.Querier(ctx)

	const query = `
		INSERT INTO pickup_calendar_blocks (
			warehouse_id,
			blocked_date,
			reason,
			created_by
		)
		VALUES ($1, $2::date, $3, $4)
		RETURNING id, warehouse_id, blocked_date, reason, created_by, created_at;
	`
	created, err := scanBlockRow(q.QueryRow(
		ctx,
		query,
		block.WarehouseID,
		block.BlockedDate,
		block.Reason,
		block.CreatedBy,
	))
	if err != nil {
		if isUniqueViolation(err) {
			return core_domain.PickupCalendarBlock{}, fmt.Errorf("%w: pickup date already blocked", core_errors.ErrConflict)
		}
		if isForeignKeyViolation(err) {
			return core_domain.PickupCalendarBlock{}, fmt.Errorf("%w: warehouse or creator not found", core_errors.ErrNotFound)
		}
		return core_domain.PickupCalendarBlock{}, fmt.Errorf("create pickup calendar block: %w", err)
	}
	return created, nil
}

func (r *PickupCalendarRepository) DeleteBlock(
	ctx context.Context,
	blockID int64,
) error {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()
	q := r.tx.Querier(ctx)

	const query = `DELETE FROM pickup_calendar_blocks WHERE id = $1;`
	tag, err := q.Exec(ctx, query, blockID)
	if err != nil {
		return fmt.Errorf("delete pickup calendar block: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%w: pickup calendar block not found", core_errors.ErrNotFound)
	}
	return nil
}
