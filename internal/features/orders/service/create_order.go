package orders_service

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (s *OrdersService) CreateOrder(
	ctx context.Context,
	actorID int64,
	actorRole string,
	input CreateOrderInput,
) (OrderDTO, error) {
	if actorID <= 0 {
		return OrderDTO{}, fmt.Errorf("%w: invalid actor id", core_errors.ErrInvalidArgument)
	}

	if actorRole != core_domain.RoleClient.String() && actorRole != core_domain.RoleAdmin.String() {
		return OrderDTO{}, fmt.Errorf("%w: only client or admin can create order", core_errors.ErrForbidden)
	}

	handoverType, err := validateHandoverType(input.HandoverType)
	if err != nil {
		return OrderDTO{}, err
	}

	if input.ReceivingWarehouseID <= 0 {
		return OrderDTO{}, fmt.Errorf("%w: invalid receiving warehouse id", core_errors.ErrInvalidArgument)
	}

	if input.DestinationWarehouseID <= 0 {
		return OrderDTO{}, fmt.Errorf("%w: invalid destination warehouse id", core_errors.ErrInvalidArgument)
	}

	if input.ProductTypeID <= 0 {
		return OrderDTO{}, fmt.Errorf("%w: invalid product type id", core_errors.ErrInvalidArgument)
	}

	if len(input.CargoPlaces) == 0 {
		return OrderDTO{}, fmt.Errorf("%w: cargo places required", core_errors.ErrInvalidArgument)
	}

	for _, cargoPlace := range input.CargoPlaces {
		if cargoPlace.CargoPlaceTypeID <= 0 {
			return OrderDTO{}, fmt.Errorf("%w: invalid cargo place type id", core_errors.ErrInvalidArgument)
		}

		if cargoPlace.Quantity <= 0 {
			return OrderDTO{}, fmt.Errorf("%w: cargo place quantity must be positive", core_errors.ErrInvalidArgument)
		}

		if cargoPlace.WeightPerPlaceKG != nil && *cargoPlace.WeightPerPlaceKG <= 0 {
			return OrderDTO{}, fmt.Errorf("%w: weight per place must be positive", core_errors.ErrInvalidArgument)
		}

		if cargoPlace.LengthCM != nil && *cargoPlace.LengthCM <= 0 {
			return OrderDTO{}, fmt.Errorf("%w: length must be positive", core_errors.ErrInvalidArgument)
		}

		if cargoPlace.WidthCM != nil && *cargoPlace.WidthCM <= 0 {
			return OrderDTO{}, fmt.Errorf("%w: width must be positive", core_errors.ErrInvalidArgument)
		}

		if cargoPlace.HeightCM != nil && *cargoPlace.HeightCM <= 0 {
			return OrderDTO{}, fmt.Errorf("%w: height must be positive", core_errors.ErrInvalidArgument)
		}
	}

	var pickup *core_domain.PickupRequest

	switch handoverType {
	case core_domain.HandoverTypePickup:
		if input.SelfDeliveryDate != nil ||
			input.SelfDeliveryTimeFrom != nil ||
			input.SelfDeliveryTimeTo != nil {
			return OrderDTO{}, fmt.Errorf("%w: self delivery schedule is not allowed for pickup order", core_errors.ErrInvalidArgument)
		}

		if input.Pickup == nil {
			return OrderDTO{}, fmt.Errorf("%w: pickup data required", core_errors.ErrInvalidArgument)
		}

		if input.Pickup.PickupAddress == "" {
			return OrderDTO{}, fmt.Errorf("%w: pickup address required", core_errors.ErrInvalidArgument)
		}

		if err := validateDate(input.Pickup.PickupDate, "pickup_date"); err != nil {
			return OrderDTO{}, err
		}

		if err := validateTimeRange(input.Pickup.PickupTimeFrom, input.Pickup.PickupTimeTo, "pickup"); err != nil {
			return OrderDTO{}, err
		}

		pickup = &core_domain.PickupRequest{
			PickupAddress:  input.Pickup.PickupAddress,
			PickupDate:     input.Pickup.PickupDate,
			PickupTimeFrom: input.Pickup.PickupTimeFrom,
			PickupTimeTo:   input.Pickup.PickupTimeTo,
			ContactName:    input.Pickup.ContactName,
			ContactPhone:   input.Pickup.ContactPhone,
			Comment:        input.Pickup.Comment,
		}

	case core_domain.HandoverTypeSelfDelivery:
		if input.Pickup != nil {
			return OrderDTO{}, fmt.Errorf("%w: pickup data is not allowed for self delivery order", core_errors.ErrInvalidArgument)
		}

		if input.SelfDeliveryDate == nil || *input.SelfDeliveryDate == "" {
			return OrderDTO{}, fmt.Errorf("%w: self delivery date required", core_errors.ErrInvalidArgument)
		}

		if err := validateDate(*input.SelfDeliveryDate, "self_delivery_date"); err != nil {
			return OrderDTO{}, err
		}

		if err := validateTimeRange(input.SelfDeliveryTimeFrom, input.SelfDeliveryTimeTo, "self_delivery"); err != nil {
			return OrderDTO{}, err
		}
	}

	cargoPlaces := make([]core_domain.OrderCargoPlace, 0, len(input.CargoPlaces))
	for _, cargoPlace := range input.CargoPlaces {
		cargoPlaces = append(cargoPlaces, core_domain.OrderCargoPlace{
			CargoPlaceTypeID: cargoPlace.CargoPlaceTypeID,
			Quantity:         cargoPlace.Quantity,
			WeightPerPlaceKG: cargoPlace.WeightPerPlaceKG,
			LengthCM:         cargoPlace.LengthCM,
			WidthCM:          cargoPlace.WidthCM,
			HeightCM:         cargoPlace.HeightCM,
			Comment:          cargoPlace.Comment,
		})
	}

	order := core_domain.Order{
		ClientID:               actorID,
		ReceivingWarehouseID:   input.ReceivingWarehouseID,
		DestinationWarehouseID: input.DestinationWarehouseID,
		ProductTypeID:          input.ProductTypeID,
		HandoverType:           handoverType,
		SelfDeliveryDate:       input.SelfDeliveryDate,
		SelfDeliveryTimeFrom:   input.SelfDeliveryTimeFrom,
		SelfDeliveryTimeTo:     input.SelfDeliveryTimeTo,
		Status:                 core_domain.OrderStatusCreated,
		Comment:                input.Comment,
	}

	details, err := s.repo.CreateOrder(ctx, order, cargoPlaces, pickup, actorID)
	if err != nil {
		return OrderDTO{}, fmt.Errorf("create order: %w", err)
	}

	return mapOrderDetailsToDTO(details), nil
}
