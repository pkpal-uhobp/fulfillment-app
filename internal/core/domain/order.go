package domain

import "time"

type HandoverType string

const (
	HandoverTypePickup       HandoverType = "pickup"
	HandoverTypeSelfDelivery HandoverType = "self_delivery"
)

func (t HandoverType) String() string {
	return string(t)
}

type OrderStatus string

const (
	OrderStatusCreated            OrderStatus = "created"
	OrderStatusWaitingPickup      OrderStatus = "waiting_pickup"
	OrderStatusWaitingDelivery    OrderStatus = "waiting_delivery"
	OrderStatusReceived           OrderStatus = "received"
	OrderStatusStored             OrderStatus = "stored"
	OrderStatusAssignedToShipping OrderStatus = "assigned_to_shipping"
	OrderStatusShipped            OrderStatus = "shipped"
	OrderStatusDelivered          OrderStatus = "delivered"
	OrderStatusCancelled          OrderStatus = "cancelled"
)

func (s OrderStatus) String() string {
	return string(s)
}

type Order struct {
	ID                     int64
	ClientID               int64
	ReceivingWarehouseID   int64
	DestinationWarehouseID int64
	ProductTypeID          int64
	HandoverType           HandoverType
	SelfDeliveryDate       *string
	SelfDeliveryTimeFrom   *string
	SelfDeliveryTimeTo     *string
	Status                 OrderStatus
	Comment                *string
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

type OrderCargoPlace struct {
	ID               int64
	OrderID          int64
	CargoPlaceTypeID int64
	Quantity         int
	WeightPerPlaceKG *float64
	LengthCM         *float64
	WidthCM          *float64
	HeightCM         *float64
	Comment          *string
	CreatedAt        time.Time
}

type PickupRequest struct {
	ID               int64
	OrderID          int64
	PickupAddress    string
	PickupDate       string
	PickupTimeFrom   *string
	PickupTimeTo     *string
	ContactName      *string
	ContactPhone     *string
	Status           string
	AssignedLogistID *int64
	Comment          *string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type OrderDetails struct {
	Order       Order
	CargoPlaces []OrderCargoPlace
	Pickup      *PickupRequest
}

type OrderFilter struct {
	ClientID     *int64
	Status       string
	HandoverType string
}
