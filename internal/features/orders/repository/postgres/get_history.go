package orders_repository_postgres

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

func (r *OrdersRepository) ListOrderStatusHistory(
	ctx context.Context,
	orderID int64,
) ([]core_domain.OrderStatusHistory, error) {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()

	q := r.tx.Querier(ctx)

	const query = `
		SELECT
			id,
			order_id,
			old_status,
			new_status,
			changed_by,
			comment,
			changed_at
		FROM order_status_history
		WHERE order_id = $1
		ORDER BY changed_at ASC, id ASC;
	`

	rows, err := q.Query(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("query order status history: %w", err)
	}
	defer rows.Close()

	history := make([]core_domain.OrderStatusHistory, 0)
	for rows.Next() {
		var item core_domain.OrderStatusHistory
		var oldStatus *string
		var newStatus string

		if err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&oldStatus,
			&newStatus,
			&item.ChangedBy,
			&item.Comment,
			&item.ChangedAt,
		); err != nil {
			return nil, fmt.Errorf("scan order status history: %w", err)
		}

		if oldStatus != nil {
			value := core_domain.OrderStatus(*oldStatus)
			item.OldStatus = &value
		}
		item.NewStatus = core_domain.OrderStatus(newStatus)

		history = append(history, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate order status history: %w", err)
	}

	return history, nil
}
