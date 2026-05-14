package domain

import "time"

type CargoItemStatus string

const (
	CargoItemStatusAccepted    CargoItemStatus = "accepted"
	CargoItemStatusStored      CargoItemStatus = "stored"
	CargoItemStatusReadyToShip CargoItemStatus = "ready_to_ship"
	CargoItemStatusShipped     CargoItemStatus = "shipped"
	CargoItemStatusLost        CargoItemStatus = "lost"
	CargoItemStatusDamaged     CargoItemStatus = "damaged"
	CargoItemStatusCancelled   CargoItemStatus = "cancelled"
)

func (s CargoItemStatus) String() string { return string(s) }

func (s CargoItemStatus) IsValid() bool {
	switch s {
	case CargoItemStatusAccepted,
		CargoItemStatusStored,
		CargoItemStatusReadyToShip,
		CargoItemStatusShipped,
		CargoItemStatusLost,
		CargoItemStatusDamaged,
		CargoItemStatusCancelled:
		return true
	default:
		return false
	}
}

func (s CargoItemStatus) IsTerminal() bool {
	switch s {
	case CargoItemStatusShipped,
		CargoItemStatusLost,
		CargoItemStatusDamaged,
		CargoItemStatusCancelled:
		return true
	default:
		return false
	}
}

type CargoItem struct {
	ID                int64
	OrderID           int64
	OrderCargoPlaceID int64
	CargoPlaceTypeID  int64
	QRCode            string
	Status            CargoItemStatus
	StorageZoneID     *int64
	GateID            *int64
	ReceivedBy        *int64
	ShippedBy         *int64
	ReceivedAt        *time.Time
	ShippedAt         *time.Time
	Comment           *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type CargoStatusHistory struct {
	ID          int64
	CargoItemID int64
	OldStatus   *CargoItemStatus
	NewStatus   CargoItemStatus
	ChangedBy   int64
	Comment     *string
	ChangedAt   time.Time
}

type CargoItemFilter struct {
	OrderID       *int64
	Status        string
	StorageZoneID *int64
	GateID        *int64
	QRCode        string
	ClientID      *int64
	Page          int
	Limit         int
}
