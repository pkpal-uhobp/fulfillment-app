package warehouses_service

import (
	"time"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

type WarehouseDTO struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	WarehouseType string    `json:"warehouse_type"`
	Marketplace   string    `json:"marketplace,omitempty"`
	City          string    `json:"city"`
	Address       string    `json:"address"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
}

type StorageZoneDTO struct {
	ID          int64  `json:"id"`
	WarehouseID int64  `json:"warehouse_id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IsActive    bool   `json:"is_active"`
}

type GateDTO struct {
	ID          int64  `json:"id"`
	WarehouseID int64  `json:"warehouse_id"`
	Name        string `json:"name"`
	IsActive    bool   `json:"is_active"`
}

type ProductTypeDTO struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IsActive    bool   `json:"is_active"`
}

type CargoPlaceTypeDTO struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IsActive    bool   `json:"is_active"`
}

type WarehouseFilter = core_domain.WarehouseFilter

type CreateWarehouseInput struct {
	Name          string
	WarehouseType string
	Marketplace   string
	City          string
	Address       string
}

type PatchWarehouseInput struct {
	Name          *string
	WarehouseType *string
	Marketplace   *string
	City          *string
	Address       *string
	IsActive      *bool
}

type CreateStorageZoneInput struct {
	WarehouseID int64
	Name        string
	Description string
}

type PatchStorageZoneInput struct {
	Name        *string
	Description *string
	IsActive    *bool
}

type CreateGateInput struct {
	WarehouseID int64
	Name        string
}

type PatchGateInput struct {
	Name     *string
	IsActive *bool
}
