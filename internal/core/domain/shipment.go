package domain

import "time"

type ShipmentStatus string

const (
	ShipmentStatusPlanned   ShipmentStatus = "planned"
	ShipmentStatusLoading   ShipmentStatus = "loading"
	ShipmentStatusShipped   ShipmentStatus = "shipped"
	ShipmentStatusCompleted ShipmentStatus = "completed"
	ShipmentStatusCancelled ShipmentStatus = "cancelled"
)

func (s ShipmentStatus) String() string {
	return string(s)
}

func (s ShipmentStatus) IsValid() bool {
	switch s {
	case ShipmentStatusPlanned,
		ShipmentStatusLoading,
		ShipmentStatusShipped,
		ShipmentStatusCompleted,
		ShipmentStatusCancelled:
		return true
	default:
		return false
	}
}

func (s ShipmentStatus) IsTerminal() bool {
	switch s {
	case ShipmentStatusCompleted, ShipmentStatusCancelled:
		return true
	default:
		return false
	}
}

type Shipment struct {
	ID                     int64
	DestinationWarehouseID int64
	GateID                 int64
	PlannedDepartureAt     time.Time
	ActualDepartureAt      *time.Time
	Status                 ShipmentStatus
	CreatedBy              int64
	CreatedAt              time.Time
}

type ShipmentItem struct {
	ID          int64
	ShipmentID  int64
	CargoItemID int64
}

type ShipmentItemDetails struct {
	Item      ShipmentItem
	CargoItem CargoItem
}

type ShipmentDetails struct {
	Shipment Shipment
	Items    []ShipmentItemDetails
}

type ShipmentFilter struct {
	Status                 string
	DestinationWarehouseID *int64
	GateID                 *int64
	Date                   string
	Page                   int
	Limit                  int
}
