package orders_repository_postgres

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

func (r *OrdersRepository) CreateOrder(
	ctx context.Context,
	order core_domain.Order,
	cargoPlaces []core_domain.OrderCargoPlace,
	pickup *core_domain.PickupRequest,
	changedBy int64,
) (core_domain.OrderDetails, error) {
	var result core_domain.OrderDetails

	if err := r.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		q := r.tx.Querier(ctx)

		query := fmt.Sprintf(`
			INSERT INTO orders (
				client_id,
				receiving_warehouse_id,
				destination_warehouse_id,
				product_type_id,
				handover_type,
				self_delivery_date,
				self_delivery_time_from,
				self_delivery_time_to,
				status,
				comment
			)
			VALUES ($1, $2, $3, $4, $5, $6::date, $7::time, $8::time, $9, $10)
			RETURNING %s;
		`, orderColumns)

		createdOrder, err := scanOrderRow(q.QueryRow(
			ctx,
			query,
			order.ClientID,
			order.ReceivingWarehouseID,
			order.DestinationWarehouseID,
			order.ProductTypeID,
			order.HandoverType.String(),
			order.SelfDeliveryDate,
			order.SelfDeliveryTimeFrom,
			order.SelfDeliveryTimeTo,
			order.Status.String(),
			order.Comment,
		))
		if err != nil {
			return fmt.Errorf("create order: %w", err)
		}

		createdCargoPlaces := make([]core_domain.OrderCargoPlace, 0, len(cargoPlaces))
		for _, cargoPlace := range cargoPlaces {
			cargoPlaceQuery := fmt.Sprintf(`
				INSERT INTO order_cargo_places (
					order_id,
					cargo_place_type_id,
					quantity,
					weight_per_place_kg,
					length_cm,
					width_cm,
					height_cm,
					comment
				)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
				RETURNING %s;
			`, cargoPlaceColumns)

			createdCargoPlace, err := scanCargoPlaceRow(q.QueryRow(
				ctx,
				cargoPlaceQuery,
				createdOrder.ID,
				cargoPlace.CargoPlaceTypeID,
				cargoPlace.Quantity,
				cargoPlace.WeightPerPlaceKG,
				cargoPlace.LengthCM,
				cargoPlace.WidthCM,
				cargoPlace.HeightCM,
				cargoPlace.Comment,
			))
			if err != nil {
				return fmt.Errorf("create order cargo place: %w", err)
			}

			createdCargoPlaces = append(createdCargoPlaces, createdCargoPlace)
		}

		var createdPickup *core_domain.PickupRequest
		if pickup != nil {
			pickupQuery := fmt.Sprintf(`
				INSERT INTO pickup_requests (
					order_id,
					pickup_address,
					pickup_date,
					pickup_time_from,
					pickup_time_to,
					contact_name,
					contact_phone,
					comment
				)
				VALUES ($1, $2, $3::date, $4::time, $5::time, $6, $7, $8)
				RETURNING %s;
			`, pickupRequestColumns)

			value, err := scanPickupRequestRow(q.QueryRow(
				ctx,
				pickupQuery,
				createdOrder.ID,
				pickup.PickupAddress,
				pickup.PickupDate,
				pickup.PickupTimeFrom,
				pickup.PickupTimeTo,
				pickup.ContactName,
				pickup.ContactPhone,
				pickup.Comment,
			))
			if err != nil {
				return fmt.Errorf("create pickup request: %w", err)
			}

			createdPickup = &value
		}

		const historyQuery = `
			INSERT INTO order_status_history (
				order_id,
				old_status,
				new_status,
				changed_by,
				comment
			)
			VALUES ($1, NULL, $2, $3, $4);
		`

		if _, err := q.Exec(
			ctx,
			historyQuery,
			createdOrder.ID,
			createdOrder.Status.String(),
			changedBy,
			createdOrder.Comment,
		); err != nil {
			return fmt.Errorf("create order status history: %w", err)
		}

		result = core_domain.OrderDetails{
			Order:       createdOrder,
			CargoPlaces: createdCargoPlaces,
			Pickup:      createdPickup,
		}

		return nil
	}); err != nil {
		return core_domain.OrderDetails{}, err
	}

	return result, nil
}
