package shipments_transport_http

import shipments_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/shipments/service"

type CreateShipmentRequest struct {
	DestinationWarehouseID int64  `json:"destination_warehouse_id" validate:"required"`
	GateID                 int64  `json:"gate_id" validate:"required"`
	PlannedDepartureAt     string `json:"planned_departure_at" validate:"required"`
}

type AddShipmentItemRequest struct {
	CargoItemID int64   `json:"cargo_item_id" validate:"required"`
	Comment     *string `json:"comment,omitempty"`
}

type UpdateShipmentStatusRequest struct {
	Status  string  `json:"status" validate:"required"`
	Comment *string `json:"comment,omitempty"`
}

type ShipmentsResponse struct {
	Shipments []shipments_service.ShipmentDTO `json:"shipments"`
}

type ShipmentResponse struct {
	Shipment shipments_service.ShipmentDTO `json:"shipment"`
}
