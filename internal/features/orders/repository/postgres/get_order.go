package orders_repository_postgres

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

func (r *OrdersRepository) GetOrder(
	ctx context.Context,
	orderID int64,
) (core_domain.OrderDetails, error) {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()

	q := r.tx.Querier(ctx)

	order, err := getOrderDetailsByID(ctx, q, orderID)
	if err != nil {
		return core_domain.OrderDetails{}, fmt.Errorf("get order: %w", err)
	}

	return order, nil
}
