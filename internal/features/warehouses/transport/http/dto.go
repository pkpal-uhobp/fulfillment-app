package warehouses_transport_http

import warehouses_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/warehouses/service"

type WarehousesResponse struct {
	Warehouses []warehouses_service.WarehouseDTO `json:"warehouses"`
}

type WarehouseResponse struct {
	Warehouse warehouses_service.WarehouseDTO `json:"warehouse"`
}

type CreateWarehouseRequest struct {
	Name          string `json:"name" validate:"required,min=2,max=255"`
	WarehouseType string `json:"warehouse_type" validate:"required"`
	Marketplace   string `json:"marketplace,omitempty" validate:"omitempty,max=100"`
	City          string `json:"city" validate:"required,min=2,max=100"`
	Address       string `json:"address" validate:"required,min=5,max=500"`
}

type PatchWarehouseRequest struct {
	Name          *string `json:"name,omitempty" validate:"omitempty,min=2,max=255"`
	WarehouseType *string `json:"warehouse_type,omitempty"`
	Marketplace   *string `json:"marketplace,omitempty" validate:"omitempty,max=100"`
	City          *string `json:"city,omitempty" validate:"omitempty,min=2,max=100"`
	Address       *string `json:"address,omitempty" validate:"omitempty,min=5,max=500"`
	IsActive      *bool   `json:"is_active,omitempty"`
}

type StorageZonesResponse struct {
	StorageZones []warehouses_service.StorageZoneDTO `json:"storage_zones"`
}

type StorageZoneResponse struct {
	StorageZone warehouses_service.StorageZoneDTO `json:"storage_zone"`
}

type CreateStorageZoneRequest struct {
	WarehouseID int64  `json:"warehouse_id" validate:"required,gt=0"`
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description,omitempty" validate:"omitempty,max=1000"`
}

type PatchStorageZoneRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

type GatesResponse struct {
	Gates []warehouses_service.GateDTO `json:"gates"`
}

type GateResponse struct {
	Gate warehouses_service.GateDTO `json:"gate"`
}

type CreateGateRequest struct {
	WarehouseID int64  `json:"warehouse_id" validate:"required,gt=0"`
	Name        string `json:"name" validate:"required,min=1,max=100"`
}

type PatchGateRequest struct {
	Name     *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	IsActive *bool   `json:"is_active,omitempty"`
}

type ProductTypesResponse struct {
	ProductTypes []warehouses_service.ProductTypeDTO `json:"product_types"`
}

type CargoPlaceTypesResponse struct {
	CargoPlaceTypes []warehouses_service.CargoPlaceTypeDTO `json:"cargo_place_types"`
}
