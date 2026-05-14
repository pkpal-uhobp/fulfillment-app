package orders_service

import core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"

func mapOrderDetailsToDTO(details core_domain.OrderDetails) OrderDTO {
	cargoPlaces := make([]CargoPlaceDTO, 0, len(details.CargoPlaces))
	for _, cargoPlace := range details.CargoPlaces {
		cargoPlaces = append(cargoPlaces, CargoPlaceDTO{
			ID:               cargoPlace.ID,
			OrderID:          cargoPlace.OrderID,
			CargoPlaceTypeID: cargoPlace.CargoPlaceTypeID,
			Quantity:         cargoPlace.Quantity,
			WeightPerPlaceKG: cargoPlace.WeightPerPlaceKG,
			LengthCM:         cargoPlace.LengthCM,
			WidthCM:          cargoPlace.WidthCM,
			HeightCM:         cargoPlace.HeightCM,
			Comment:          cargoPlace.Comment,
			CreatedAt:        cargoPlace.CreatedAt,
		})
	}

	var pickup *PickupDTO
	if details.Pickup != nil {
		pickup = &PickupDTO{
			ID:               details.Pickup.ID,
			OrderID:          details.Pickup.OrderID,
			PickupAddress:    details.Pickup.PickupAddress,
			PickupDate:       details.Pickup.PickupDate,
			PickupTimeFrom:   details.Pickup.PickupTimeFrom,
			PickupTimeTo:     details.Pickup.PickupTimeTo,
			ContactName:      details.Pickup.ContactName,
			ContactPhone:     details.Pickup.ContactPhone,
			Status:           details.Pickup.Status,
			AssignedLogistID: details.Pickup.AssignedLogistID,
			Comment:          details.Pickup.Comment,
			CreatedAt:        details.Pickup.CreatedAt,
			UpdatedAt:        details.Pickup.UpdatedAt,
		}
	}

	return OrderDTO{
		ID:                     details.Order.ID,
		ClientID:               details.Order.ClientID,
		ReceivingWarehouseID:   details.Order.ReceivingWarehouseID,
		DestinationWarehouseID: details.Order.DestinationWarehouseID,
		ProductTypeID:          details.Order.ProductTypeID,
		HandoverType:           details.Order.HandoverType.String(),
		SelfDeliveryDate:       details.Order.SelfDeliveryDate,
		SelfDeliveryTimeFrom:   details.Order.SelfDeliveryTimeFrom,
		SelfDeliveryTimeTo:     details.Order.SelfDeliveryTimeTo,
		Status:                 details.Order.Status.String(),
		Comment:                details.Order.Comment,
		CargoPlaces:            cargoPlaces,
		Pickup:                 pickup,
		CreatedAt:              details.Order.CreatedAt,
		UpdatedAt:              details.Order.UpdatedAt,
	}
}

func mapOrderDetailsListToDTO(details []core_domain.OrderDetails) []OrderDTO {
	orders := make([]OrderDTO, 0, len(details))
	for _, order := range details {
		orders = append(orders, mapOrderDetailsToDTO(order))
	}

	return orders
}
