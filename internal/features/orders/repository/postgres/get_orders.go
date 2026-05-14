package orders_repository_postgres

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

func (r *OrdersRepository) ListOrders(
	ctx context.Context,
	filter core_domain.OrderFilter,
) ([]core_domain.OrderDetails, error) {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()

	q := r.tx.Querier(ctx)

	query := fmt.Sprintf(`
		SELECT %s
		FROM orders
		WHERE ($1::bigint IS NULL OR client_id = $1)
			AND ($2::varchar = '' OR status = $2)
			AND ($3::varchar = '' OR handover_type = $3)
		ORDER BY id DESC;
	`, orderColumns)

	rows, err := q.Query(
		ctx,
		query,
		filter.ClientID,
		filter.Status,
		filter.HandoverType,
	)
	if err != nil {
		return nil, fmt.Errorf("query orders: %w", err)
	}
	defer rows.Close()

	orders := make([]core_domain.OrderDetails, 0)
	for rows.Next() {
		order, err := scanOrderRow(rows)
		if err != nil {
			return nil, fmt.Errorf("scan order: %w", err)
		}

		cargoPlaces, err := listCargoPlacesByOrderID(ctx, q, order.ID)
		if err != nil {
			return nil, err
		}

		pickup, err := getPickupRequestByOrderID(ctx, q, order.ID)
		if err != nil {
			return nil, err
		}

		orders = append(orders, core_domain.OrderDetails{
			Order:       order,
			CargoPlaces: cargoPlaces,
			Pickup:      pickup,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate orders: %w", err)
	}

	return orders, nil
}
