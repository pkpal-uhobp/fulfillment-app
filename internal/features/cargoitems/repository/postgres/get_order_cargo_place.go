package cargoitems_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (r *CargoItemsRepository) GetOrderCargoPlaceForOrder(
	ctx context.Context,
	orderID int64,
	orderCargoPlaceID int64,
) (core_domain.OrderCargoPlace, core_domain.OrderStatus, error) {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()

	q := r.tx.Querier(ctx)

	const query = `
		SELECT
			ocp.id,
			ocp.order_id,
			ocp.cargo_place_type_id,
			ocp.quantity,
			ocp.weight_per_place_kg::float8,
			ocp.length_cm::float8,
			ocp.width_cm::float8,
			ocp.height_cm::float8,
			ocp.comment,
			ocp.created_at,
			o.status
		FROM order_cargo_places ocp
		JOIN orders o ON o.id = ocp.order_id
		WHERE ocp.id = $1 AND ocp.order_id = $2;
	`

	var cargoPlace core_domain.OrderCargoPlace
	var orderStatus string

	if err := q.QueryRow(ctx, query, orderCargoPlaceID, orderID).Scan(
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
		&orderStatus,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core_domain.OrderCargoPlace{}, "", fmt.Errorf("%w: order cargo place not found", core_errors.ErrNotFound)
		}
		return core_domain.OrderCargoPlace{}, "", fmt.Errorf("scan order cargo place: %w", err)
	}

	return cargoPlace, core_domain.OrderStatus(orderStatus), nil
}

func (r *CargoItemsRepository) CountCargoItemsByOrderCargoPlace(
	ctx context.Context,
	orderCargoPlaceID int64,
) (int, error) {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()

	q := r.tx.Querier(ctx)

	const query = `
		SELECT count(*)
		FROM cargo_items
		WHERE order_cargo_place_id = $1;
	`

	var count int
	if err := q.QueryRow(ctx, query, orderCargoPlaceID).Scan(&count); err != nil {
		return 0, fmt.Errorf("count cargo items by order cargo place: %w", err)
	}

	return count, nil
}
