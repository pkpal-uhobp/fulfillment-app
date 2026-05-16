package shipments_service

import "time"

type CreateShipmentInput struct {
	DestinationWarehouseID int64
	GateID                 int64
	PlannedDepartureAt     string
}

type ShipmentFilter struct {
	Status                 string
	DestinationWarehouseID *int64
	GateID                 *int64
	Date                   string
	Page                   int
	Limit                  int
}

type AddShipmentItemInput struct {
	CargoItemID int64
	Comment     *string
}

type UpdateShipmentStatusInput struct {
	Status  string
	Comment *string
}

type ShipmentDTO struct {
	ID                     int64             `json:"id"`
	DestinationWarehouseID int64             `json:"destination_warehouse_id"`
	GateID                 int64             `json:"gate_id"`
	PlannedDepartureAt     time.Time         `json:"planned_departure_at"`
	ActualDepartureAt      *time.Time        `json:"actual_departure_at,omitempty"`
	Status                 string            `json:"status"`
	CreatedBy              int64             `json:"created_by"`
	Items                  []ShipmentItemDTO `json:"items,omitempty"`
	CreatedAt              time.Time         `json:"created_at"`
}

type ShipmentItemDTO struct {
	ID            int64      `json:"id"`
	ShipmentID    int64      `json:"shipment_id"`
	CargoItemID   int64      `json:"cargo_item_id"`
	OrderID       int64      `json:"order_id"`
	QRCode        string     `json:"qr_code"`
	Status        string     `json:"status"`
	StorageZoneID *int64     `json:"storage_zone_id,omitempty"`
	GateID        *int64     `json:"gate_id,omitempty"`
	ReceivedAt    *time.Time `json:"received_at,omitempty"`
	ShippedAt     *time.Time `json:"shipped_at,omitempty"`
}
