package pickupcalendar_service

import "time"

type CalendarFilter struct {
	WarehouseID int64
	DateFrom    string
	DateTo      string
}

type BlockDateInput struct {
	WarehouseID int64
	BlockedDate string
	Reason      *string
}

type SetCapacityInput struct {
	WarehouseID   int64
	PickupDate    string
	MaxOrders     int
	CurrentOrders int
	IsClosed      bool
}

type PickupCalendarDayDTO struct {
	WarehouseID   int64                   `json:"warehouse_id"`
	Date          string                  `json:"date"`
	MaxOrders     int                     `json:"max_orders"`
	CurrentOrders int                     `json:"current_orders"`
	IsClosed      bool                    `json:"is_closed"`
	Block         *PickupCalendarBlockDTO `json:"block,omitempty"`
}

type PickupCalendarBlockDTO struct {
	ID          int64     `json:"id"`
	WarehouseID int64     `json:"warehouse_id"`
	BlockedDate string    `json:"blocked_date"`
	Reason      *string   `json:"reason,omitempty"`
	CreatedBy   int64     `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

type PickupCapacityDTO struct {
	ID            int64  `json:"id"`
	WarehouseID   int64  `json:"warehouse_id"`
	PickupDate    string `json:"pickup_date"`
	MaxOrders     int    `json:"max_orders"`
	CurrentOrders int    `json:"current_orders"`
	IsClosed      bool   `json:"is_closed"`
}
