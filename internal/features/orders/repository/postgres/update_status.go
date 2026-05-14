package orders_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (r *OrdersRepository) UpdateOrderStatus(
	ctx context.Context,
	orderID int64,
	status string,
	changedBy int64,
	comment *string,
) error {
	if err := r.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		q := r.tx.Querier(ctx)

		const selectQuery = `
			SELECT status
			FROM orders
			WHERE id = $1
			FOR UPDATE;
		`

		var oldStatus string
		if err := q.QueryRow(ctx, selectQuery, orderID).Scan(&oldStatus); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("%w: order not found", core_errors.ErrNotFound)
			}

			return fmt.Errorf("select order status: %w", err)
		}

		const updateOrderQuery = `
			UPDATE orders
			SET status = $2
			WHERE id = $1;
		`

		if _, err := q.Exec(ctx, updateOrderQuery, orderID, status); err != nil {
			return fmt.Errorf("update order status: %w", err)
		}

		const historyQuery = `
			INSERT INTO order_status_history (
				order_id,
				old_status,
				new_status,
				changed_by,
				comment
			)
			VALUES ($1, $2, $3, $4, $5);
		`

		if _, err := q.Exec(
			ctx,
			historyQuery,
			orderID,
			oldStatus,
			status,
			changedBy,
			comment,
		); err != nil {
			return fmt.Errorf("create order status history: %w", err)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
