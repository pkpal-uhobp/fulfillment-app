package orders_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
	core_postgres_tx "github.com/pkpal-uhobp/fulfillment-app/internal/core/repository/tx"
)

const orderColumns = `
	id,
	client_id,
	receiving_warehouse_id,
	destination_warehouse_id,
	product_type_id,
	handover_type,
	self_delivery_date::text,
	self_delivery_time_from::text,
	self_delivery_time_to::text,
	status,
	comment,
	created_at,
	updated_at
`

const cargoPlaceColumns = `
	id,
	order_id,
	cargo_place_type_id,
	quantity,
	weight_per_place_kg::float8,
	length_cm::float8,
	width_cm::float8,
	height_cm::float8,
	comment,
	created_at
`

const pickupRequestColumns = `
	id,
	order_id,
	pickup_address,
	pickup_date::text,
	pickup_time_from::text,
	pickup_time_to::text,
	contact_name,
	contact_phone,
	status,
	assigned_logist_id,
	comment,
	created_at,
	updated_at
`

func scanOrderRow(row pgx.Row) (core_domain.Order, error) {
	var order core_domain.Order
	var handoverType string
	var status string

	if err := row.Scan(
		&order.ID,
		&order.ClientID,
		&order.ReceivingWarehouseID,
		&order.DestinationWarehouseID,
		&order.ProductTypeID,
		&handoverType,
		&order.SelfDeliveryDate,
		&order.SelfDeliveryTimeFrom,
		&order.SelfDeliveryTimeTo,
		&status,
		&order.Comment,
		&order.CreatedAt,
		&order.UpdatedAt,
	); err != nil {
		return core_domain.Order{}, err
	}

	order.HandoverType = core_domain.HandoverType(handoverType)
	order.Status = core_domain.OrderStatus(status)

	return order, nil
}

func scanCargoPlaceRow(row pgx.Row) (core_domain.OrderCargoPlace, error) {
	var cargoPlace core_domain.OrderCargoPlace

	if err := row.Scan(
		&cargoPlace.ID,
		&cargoPlace.OrderID,
		&cargoPlace.CargoPlaceTypeID,
		&cargoPlace.Quantity,
		&cargoPlace.WeightPerPlaceKG,
		&cargoPlace.LengthCM,
		&cargoPlace.WidthCM,
		&cargoPlace.HeightCM,
		&cargoPlace.Comment,
		&cargoPlace.CreatedAt,
	); err != nil {
		return core_domain.OrderCargoPlace{}, err
	}

	return cargoPlace, nil
}

func scanPickupRequestRow(row pgx.Row) (core_domain.PickupRequest, error) {
	var pickup core_domain.PickupRequest

	if err := row.Scan(
		&pickup.ID,
		&pickup.OrderID,
		&pickup.PickupAddress,
		&pickup.PickupDate,
		&pickup.PickupTimeFrom,
		&pickup.PickupTimeTo,
		&pickup.ContactName,
		&pickup.ContactPhone,
		&pickup.Status,
		&pickup.AssignedLogistID,
		&pickup.Comment,
		&pickup.CreatedAt,
		&pickup.UpdatedAt,
	); err != nil {
		return core_domain.PickupRequest{}, err
	}

	return pickup, nil
}

func listCargoPlacesByOrderID(
	ctx context.Context,
	q core_postgres_tx.Querier,
	orderID int64,
) ([]core_domain.OrderCargoPlace, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM order_cargo_places
		WHERE order_id = $1
		ORDER BY id;
	`, cargoPlaceColumns)

	rows, err := q.Query(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("query order cargo places: %w", err)
	}
	defer rows.Close()

	cargoPlaces := make([]core_domain.OrderCargoPlace, 0)
	for rows.Next() {
		cargoPlace, err := scanCargoPlaceRow(rows)
		if err != nil {
			return nil, fmt.Errorf("scan order cargo place: %w", err)
		}

		cargoPlaces = append(cargoPlaces, cargoPlace)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate order cargo places: %w", err)
	}

	return cargoPlaces, nil
}

func getPickupRequestByOrderID(
	ctx context.Context,
	q core_postgres_tx.Querier,
	orderID int64,
) (*core_domain.PickupRequest, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM pickup_requests
		WHERE order_id = $1;
	`, pickupRequestColumns)

	pickup, err := scanPickupRequestRow(q.QueryRow(ctx, query, orderID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("scan pickup request: %w", err)
	}

	return &pickup, nil
}

func getOrderDetailsByID(
	ctx context.Context,
	q core_postgres_tx.Querier,
	orderID int64,
) (core_domain.OrderDetails, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM orders
		WHERE id = $1;
	`, orderColumns)

	order, err := scanOrderRow(q.QueryRow(ctx, query, orderID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core_domain.OrderDetails{}, fmt.Errorf("%w: order not found", core_errors.ErrNotFound)
		}

		return core_domain.OrderDetails{}, fmt.Errorf("scan order: %w", err)
	}

	cargoPlaces, err := listCargoPlacesByOrderID(ctx, q, order.ID)
	if err != nil {
		return core_domain.OrderDetails{}, err
	}

	pickup, err := getPickupRequestByOrderID(ctx, q, order.ID)
	if err != nil {
		return core_domain.OrderDetails{}, err
	}

	return core_domain.OrderDetails{
		Order:       order,
		CargoPlaces: cargoPlaces,
		Pickup:      pickup,
	}, nil
}
