package cargoitems_service

import "time"

type CreateCargoItemInput struct {
	OrderID           int64
	OrderCargoPlaceID int64
	QRCode            *string
	Comment           *string
}

type CargoItemFilter struct {
	OrderID       *int64
	Status        string
	StorageZoneID *int64
	GateID        *int64
	QRCode        string
	Page          int
	Limit         int
}

type UpdateCargoItemStatusInput struct {
	Status  string
	Comment *string
}

type AssignStorageZoneInput struct {
	StorageZoneID int64
	Comment       *string
}

type AssignGateInput struct {
	GateID  int64
	Comment *string
}

type CargoItemDTO struct {
	ID                int64      `json:"id"`
	OrderID           int64      `json:"order_id"`
	OrderCargoPlaceID int64      `json:"order_cargo_place_id"`
	CargoPlaceTypeID  int64      `json:"cargo_place_type_id"`
	QRCode            string     `json:"qr_code"`
	Status            string     `json:"status"`
	StorageZoneID     *int64     `json:"storage_zone_id,omitempty"`
	GateID            *int64     `json:"gate_id,omitempty"`
	ReceivedBy        *int64     `json:"received_by,omitempty"`
	ShippedBy         *int64     `json:"shipped_by,omitempty"`
	ReceivedAt        *time.Time `json:"received_at,omitempty"`
	ShippedAt         *time.Time `json:"shipped_at,omitempty"`
	Comment           *string    `json:"comment,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type CargoItemLabelDTO struct {
	CargoItemID       int64      `json:"cargo_item_id"`
	OrderID           int64      `json:"order_id"`
	OrderCargoPlaceID int64      `json:"order_cargo_place_id"`
	CargoPlaceTypeID  int64      `json:"cargo_place_type_id"`
	QRCode            string     `json:"qr_code"`
	QRCodeValue       string     `json:"qr_code_value"`
	Status            string     `json:"status"`
	StorageZoneID     *int64     `json:"storage_zone_id,omitempty"`
	GateID            *int64     `json:"gate_id,omitempty"`
	ReceivedAt        *time.Time `json:"received_at,omitempty"`
	ShippedAt         *time.Time `json:"shipped_at,omitempty"`
	LabelText         string     `json:"label_text"`
}

type CargoStatusHistoryDTO struct {
	ID          int64     `json:"id"`
	CargoItemID int64     `json:"cargo_item_id"`
	OldStatus   *string   `json:"old_status,omitempty"`
	NewStatus   string    `json:"new_status"`
	ChangedBy   int64     `json:"changed_by"`
	Comment     *string   `json:"comment,omitempty"`
	ChangedAt   time.Time `json:"changed_at"`
}
