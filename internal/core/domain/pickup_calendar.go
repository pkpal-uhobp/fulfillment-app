package domain

import "time"

type PickupCalendarBlock struct {
	ID          int64
	WarehouseID int64
	BlockedDate string
	Reason      *string
	CreatedBy   int64
	CreatedAt   time.Time
}

type PickupCapacity struct {
	ID            int64
	WarehouseID   int64
	PickupDate    string
	MaxOrders     int
	CurrentOrders int
	IsClosed      bool
}

type PickupCalendarDay struct {
	WarehouseID   int64
	Date          string
	MaxOrders     int
	CurrentOrders int
	IsClosed      bool
	Block         *PickupCalendarBlock
}

type PickupCalendarFilter struct {
	WarehouseID int64
	DateFrom    string
	DateTo      string
}
