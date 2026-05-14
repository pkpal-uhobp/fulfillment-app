package orders_repository_postgres

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_postgres_tx "github.com/pkpal-uhobp/fulfillment-app/internal/core/repository/tx"
)

const (
	defaultRepositoryOrdersPage  = 1
	defaultRepositoryOrdersLimit = 20
)

func (r *OrdersRepository) ListOrders(
	ctx context.Context,
	filter core_domain.OrderFilter,
) ([]core_domain.OrderDetails, error) {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()

	q := r.tx.Querier(ctx)

	page := filter.Page
	if page <= 0 {
		page = defaultRepositoryOrdersPage
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = defaultRepositoryOrdersLimit
	}
	offset := (page - 1) * limit

	query := fmt.Sprintf(`
		SELECT %s
		FROM orders
		WHERE ($1::bigint IS NULL OR client_id = $1)
		  AND ($2::varchar = '' OR status = $2)
		  AND ($3::varchar = '' OR handover_type = $3)
		  AND (
			$4::bigint IS NULL
			OR receiving_warehouse_id = $4
			OR destination_warehouse_id = $4
		  )
		  AND ($5::bigint IS NULL OR receiving_warehouse_id = $5)
		  AND ($6::bigint IS NULL OR destination_warehouse_id = $6)
		ORDER BY id DESC
		LIMIT $7 OFFSET $8;
	`, orderColumns)

	rows, err := q.Query(
		ctx,
		query,
		filter.ClientID,
		filter.Status,
		filter.HandoverType,
		filter.WarehouseID,
		filter.ReceivingWarehouseID,
		filter.DestinationWarehouseID,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("query orders: %w", err)
	}
	defer rows.Close()

	orders := make([]core_domain.OrderDetails, 0)
	orderIDs := make([]int64, 0)

	for rows.Next() {
		order, err := scanOrderRow(rows)
		if err != nil {
			return nil, fmt.Errorf("scan order: %w", err)
		}

		orders = append(orders, core_domain.OrderDetails{Order: order})
		orderIDs = append(orderIDs, order.ID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate orders: %w", err)
	}

	cargoPlacesByOrderID, err := listCargoPlacesByOrderIDs(ctx, q, orderIDs)
	if err != nil {
		return nil, err
	}

	pickupByOrderID, err := listPickupRequestsByOrderIDs(ctx, q, orderIDs)
	if err != nil {
		return nil, err
	}

	for i := range orders {
		orderID := orders[i].Order.ID
		orders[i].CargoPlaces = cargoPlacesByOrderID[orderID]
		orders[i].Pickup = pickupByOrderID[orderID]
	}

	return orders, nil
}

func listCargoPlacesByOrderIDs(
	ctx context.Context,
	q core_postgres_tx.Querier,
	orderIDs []int64,
) (map[int64][]core_domain.OrderCargoPlace, error) {
	result := make(map[int64][]core_domain.OrderCargoPlace, len(orderIDs))
	if len(orderIDs) == 0 {
		return result, nil
	}

	query := fmt.Sprintf(`
		SELECT %s
		FROM order_cargo_places
		WHERE order_id = ANY($1::bigint[])
		ORDER BY order_id, id;
	`, cargoPlaceColumns)

	rows, err := q.Query(ctx, query, orderIDs)
	if err != nil {
		return nil, fmt.Errorf("query order cargo places batch: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		cargoPlace, err := scanCargoPlaceRow(rows)
		if err != nil {
			return nil, fmt.Errorf("scan order cargo place batch: %w", err)
		}
		result[cargoPlace.OrderID] = append(result[cargoPlace.OrderID], cargoPlace)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate order cargo places batch: %w", err)
	}

	return result, nil
}

func listPickupRequestsByOrderIDs(
	ctx context.Context,
	q core_postgres_tx.Querier,
	orderIDs []int64,
) (map[int64]*core_domain.PickupRequest, error) {
	result := make(map[int64]*core_domain.PickupRequest, len(orderIDs))
	if len(orderIDs) == 0 {
		return result, nil
	}

	query := fmt.Sprintf(`
		SELECT %s
		FROM pickup_requests
		WHERE order_id = ANY($1::bigint[])
		ORDER BY order_id, id;
	`, pickupRequestColumns)

	rows, err := q.Query(ctx, query, orderIDs)
	if err != nil {
		return nil, fmt.Errorf("query pickup requests batch: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		pickup, err := scanPickupRequestRow(rows)
		if err != nil {
			return nil, fmt.Errorf("scan pickup request batch: %w", err)
		}
		pickupCopy := pickup
		result[pickup.OrderID] = &pickupCopy
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate pickup requests batch: %w", err)
	}

	return result, nil
}
