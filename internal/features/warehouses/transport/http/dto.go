package warehouses_transport_http

import warehouses_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/warehouses/service"

type CreateWarehouseRequest struct {
	Name          string `json:"name" validate:"required"`
	WarehouseType string `json:"warehouse_type" validate:"required"`
	Marketplace   string `json:"marketplace,omitempty"`
	City          string `json:"city" validate:"required"`
	Address       string `json:"address" validate:"required"`
}

type PatchWarehouseRequest struct {
	Name          *string `json:"name,omitempty"`
	WarehouseType *string `json:"warehouse_type,omitempty"`
	Marketplace   *string `json:"marketplace,omitempty"`
	City          *string `json:"city,omitempty"`
	Address       *string `json:"address,omitempty"`
	IsActive      *bool   `json:"is_active,omitempty"`
}

type CreateStorageZoneRequest struct {
	WarehouseID int64  `json:"warehouse_id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`
}

type PatchStorageZoneRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

type CreateGateRequest struct {
	WarehouseID int64  `json:"warehouse_id" validate:"required"`
	Name        string `json:"name" validate:"required"`
}

type PatchGateRequest struct {
	Name     *string `json:"name,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}

type WarehousesResponse struct {
	Warehouses []warehouses_service.WarehouseDTO `json:"warehouses"`
}

type WarehouseResponse struct {
	Warehouse warehouses_service.WarehouseDTO `json:"warehouse"`
}

type StorageZonesResponse struct {
	StorageZones []warehouses_service.StorageZoneDTO `json:"storage_zones"`
}

type StorageZoneResponse struct {
	StorageZone warehouses_service.StorageZoneDTO `json:"storage_zone"`
}

type GatesResponse struct {
	Gates []warehouses_service.GateDTO `json:"gates"`
}

type GateResponse struct {
	Gate warehouses_service.GateDTO `json:"gate"`
}

type ProductTypesResponse struct {
	ProductTypes []warehouses_service.ProductTypeDTO `json:"product_types"`
}

type CargoPlaceTypesResponse struct {
	CargoPlaceTypes []warehouses_service.CargoPlaceTypeDTO `json:"cargo_place_types"`
}
