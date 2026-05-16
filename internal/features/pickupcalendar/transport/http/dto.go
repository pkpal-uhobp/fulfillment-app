package pickupcalendar_transport_http

import pickupcalendar_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/pickupcalendar/service"

type BlockDateRequest struct {
	WarehouseID int64   `json:"warehouse_id" validate:"required,gt=0"`
	BlockedDate string  `json:"blocked_date" validate:"required"`
	Reason      *string `json:"reason,omitempty"`
}

type SetCapacityRequest struct {
	WarehouseID   int64  `json:"warehouse_id" validate:"required,gt=0"`
	PickupDate    string `json:"pickup_date" validate:"required"`
	MaxOrders     int    `json:"max_orders" validate:"gte=0"`
	CurrentOrders int    `json:"current_orders" validate:"gte=0"`
	IsClosed      bool   `json:"is_closed"`
}

type PickupCalendarResponse struct {
	Days []pickupcalendar_service.PickupCalendarDayDTO `json:"days"`
}

type PickupCalendarBlockResponse struct {
	Block pickupcalendar_service.PickupCalendarBlockDTO `json:"block"`
}

type PickupCapacityResponse struct {
	Capacity pickupcalendar_service.PickupCapacityDTO `json:"capacity"`
}
