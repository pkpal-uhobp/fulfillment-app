package cargoitems_transport_http

import cargoitems_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/cargoitems/service"

type CreateCargoItemRequest struct {
	OrderCargoPlaceID int64   `json:"order_cargo_place_id" validate:"required,gt=0"`
	QRCode            *string `json:"qr_code,omitempty"`
	Comment           *string `json:"comment,omitempty"`
}

type UpdateCargoItemStatusRequest struct {
	Status  string  `json:"status" validate:"required"`
	Comment *string `json:"comment,omitempty"`
}

type AssignStorageZoneRequest struct {
	StorageZoneID int64   `json:"storage_zone_id" validate:"required,gt=0"`
	Comment       *string `json:"comment,omitempty"`
}

type AssignGateRequest struct {
	GateID  int64   `json:"gate_id" validate:"required,gt=0"`
	Comment *string `json:"comment,omitempty"`
}

type CargoItemResponse struct {
	CargoItem cargoitems_service.CargoItemDTO `json:"cargo_item"`
}

type CargoItemsResponse struct {
	CargoItems []cargoitems_service.CargoItemDTO `json:"cargo_items"`
}

type CargoItemHistoryResponse struct {
	History []cargoitems_service.CargoStatusHistoryDTO `json:"history"`
}

type CargoItemLabelResponse struct {
	Label cargoitems_service.CargoItemLabelDTO `json:"label"`
}
