package pickupcalendar_repository_postgres

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (r *PickupCalendarRepository) SetCapacity(
	ctx context.Context,
	capacity core_domain.PickupCapacity,
) (core_domain.PickupCapacity, error) {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()
	q := r.tx.Querier(ctx)

	const query = `
		INSERT INTO pickup_capacity (
			warehouse_id,
			pickup_date,
			max_orders,
			current_orders,
			is_closed
		)
		VALUES ($1, $2::date, $3, $4, $5)
		ON CONFLICT (warehouse_id, pickup_date)
		DO UPDATE SET
			max_orders = EXCLUDED.max_orders,
			current_orders = EXCLUDED.current_orders,
			is_closed = EXCLUDED.is_closed
		RETURNING id, warehouse_id, pickup_date, max_orders, current_orders, is_closed;
	`
	updated, err := scanCapacityRow(q.QueryRow(
		ctx,
		query,
		capacity.WarehouseID,
		capacity.PickupDate,
		capacity.MaxOrders,
		capacity.CurrentOrders,
		capacity.IsClosed,
	))
	if err != nil {
		if isForeignKeyViolation(err) {
			return core_domain.PickupCapacity{}, fmt.Errorf("%w: warehouse not found", core_errors.ErrNotFound)
		}
		if isCheckViolation(err) {
			return core_domain.PickupCapacity{}, fmt.Errorf("%w: invalid pickup capacity values", core_errors.ErrInvalidArgument)
		}
		return core_domain.PickupCapacity{}, fmt.Errorf("set pickup capacity: %w", err)
	}
	return updated, nil
}
