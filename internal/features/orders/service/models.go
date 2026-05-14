package orders_service

import "time"

type CargoPlaceInput struct {
	CargoPlaceTypeID int64
	Quantity         int
	WeightPerPlaceKG *float64
	LengthCM         *float64
	WidthCM          *float64
	HeightCM         *float64
	Comment          *string
}

type PickupInput struct {
	PickupAddress  string
	PickupDate     string
	PickupTimeFrom *string
	PickupTimeTo   *string
	ContactName    *string
	ContactPhone   *string
	Comment        *string
}

type CreateOrderInput struct {
	ReceivingWarehouseID   int64
	DestinationWarehouseID int64
	ProductTypeID          int64
	HandoverType           string
	SelfDeliveryDate       *string
	SelfDeliveryTimeFrom   *string
	SelfDeliveryTimeTo     *string
	Comment                *string
	CargoPlaces            []CargoPlaceInput
	Pickup                 *PickupInput
}

type OrderFilter struct {
	ClientID               *int64
	Status                 string
	HandoverType           string
	WarehouseID            *int64
	ReceivingWarehouseID   *int64
	DestinationWarehouseID *int64
	Page                   int
	Limit                  int
}

type UpdateOrderStatusInput struct {
	Status  string
	Comment *string
}

type CancelOrderInput struct {
	Comment *string
}

type OrderDTO struct {
	ID                     int64           `json:"id"`
	ClientID               int64           `json:"client_id"`
	ReceivingWarehouseID   int64           `json:"receiving_warehouse_id"`
	DestinationWarehouseID int64           `json:"destination_warehouse_id"`
	ProductTypeID          int64           `json:"product_type_id"`
	HandoverType           string          `json:"handover_type"`
	SelfDeliveryDate       *string         `json:"self_delivery_date,omitempty"`
	SelfDeliveryTimeFrom   *string         `json:"self_delivery_time_from,omitempty"`
	SelfDeliveryTimeTo     *string         `json:"self_delivery_time_to,omitempty"`
	Status                 string          `json:"status"`
	Comment                *string         `json:"comment,omitempty"`
	CargoPlaces            []CargoPlaceDTO `json:"cargo_places"`
	Pickup                 *PickupDTO      `json:"pickup,omitempty"`
	CreatedAt              time.Time       `json:"created_at"`
	UpdatedAt              time.Time       `json:"updated_at"`
}

type CargoPlaceDTO struct {
	ID               int64     `json:"id"`
	OrderID          int64     `json:"order_id"`
	CargoPlaceTypeID int64     `json:"cargo_place_type_id"`
	Quantity         int       `json:"quantity"`
	WeightPerPlaceKG *float64  `json:"weight_per_place_kg,omitempty"`
	LengthCM         *float64  `json:"length_cm,omitempty"`
	WidthCM          *float64  `json:"width_cm,omitempty"`
	HeightCM         *float64  `json:"height_cm,omitempty"`
	Comment          *string   `json:"comment,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

type PickupDTO struct {
	ID               int64     `json:"id"`
	OrderID          int64     `json:"order_id"`
	PickupAddress    string    `json:"pickup_address"`
	PickupDate       string    `json:"pickup_date"`
	PickupTimeFrom   *string   `json:"pickup_time_from,omitempty"`
	PickupTimeTo     *string   `json:"pickup_time_to,omitempty"`
	ContactName      *string   `json:"contact_name,omitempty"`
	ContactPhone     *string   `json:"contact_phone,omitempty"`
	Status           string    `json:"status"`
	AssignedLogistID *int64    `json:"assigned_logist_id,omitempty"`
	Comment          *string   `json:"comment,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type OrderStatusHistoryDTO struct {
	ID        int64     `json:"id"`
	OrderID   int64     `json:"order_id"`
	OldStatus *string   `json:"old_status,omitempty"`
	NewStatus string    `json:"new_status"`
	ChangedBy int64     `json:"changed_by"`
	Comment   *string   `json:"comment,omitempty"`
	ChangedAt time.Time `json:"changed_at"`
}
