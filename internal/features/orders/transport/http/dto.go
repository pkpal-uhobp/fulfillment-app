package orders_transport_http

import orders_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/orders/service"

type CreateOrderRequest struct {
	ReceivingWarehouseID   int64                     `json:"receiving_warehouse_id" validate:"required"`
	DestinationWarehouseID int64                     `json:"destination_warehouse_id" validate:"required"`
	ProductTypeID          int64                     `json:"product_type_id" validate:"required"`
	HandoverType           string                    `json:"handover_type" validate:"required"`
	SelfDeliveryDate       *string                   `json:"self_delivery_date,omitempty"`
	SelfDeliveryTimeFrom   *string                   `json:"self_delivery_time_from,omitempty"`
	SelfDeliveryTimeTo     *string                   `json:"self_delivery_time_to,omitempty"`
	Comment                *string                   `json:"comment,omitempty"`
	CargoPlaces            []CreateCargoPlaceRequest `json:"cargo_places" validate:"required,min=1,dive"`
	Pickup                 *CreatePickupRequest      `json:"pickup,omitempty"`
}

type CreateCargoPlaceRequest struct {
	CargoPlaceTypeID int64    `json:"cargo_place_type_id" validate:"required"`
	Quantity         int      `json:"quantity" validate:"required"`
	WeightPerPlaceKG *float64 `json:"weight_per_place_kg,omitempty"`
	LengthCM         *float64 `json:"length_cm,omitempty"`
	WidthCM          *float64 `json:"width_cm,omitempty"`
	HeightCM         *float64 `json:"height_cm,omitempty"`
	Comment          *string  `json:"comment,omitempty"`
}

type CreatePickupRequest struct {
	PickupAddress  string  `json:"pickup_address" validate:"required"`
	PickupDate     string  `json:"pickup_date" validate:"required"`
	PickupTimeFrom *string `json:"pickup_time_from,omitempty"`
	PickupTimeTo   *string `json:"pickup_time_to,omitempty"`
	ContactName    *string `json:"contact_name,omitempty"`
	ContactPhone   *string `json:"contact_phone,omitempty"`
	Comment        *string `json:"comment,omitempty"`
}

type CancelOrderRequest struct {
	Comment *string `json:"comment,omitempty"`
}

type UpdateOrderStatusRequest struct {
	Status  string  `json:"status" validate:"required"`
	Comment *string `json:"comment,omitempty"`
}

type OrdersResponse struct {
	Orders []orders_service.OrderDTO `json:"orders"`
}

type OrderResponse struct {
	Order orders_service.OrderDTO `json:"order"`
}

type OrderHistoryResponse struct {
	History []orders_service.OrderStatusHistoryDTO `json:"history"`
}
