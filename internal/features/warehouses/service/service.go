package warehouses_service

import (
	"context"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

type WarehousesRepository interface {
	ListWarehouses(
		ctx context.Context,
		filter WarehouseFilter,
	) ([]core_domain.Warehouse, error)

	GetWarehouseByID(
		ctx context.Context,
		warehouseID int64,
	) (core_domain.Warehouse, error)

	CreateWarehouse(
		ctx context.Context,
		warehouse core_domain.Warehouse,
	) (core_domain.Warehouse, error)

	PatchWarehouse(
		ctx context.Context,
		warehouseID int64,
		patch core_domain.WarehousePatch,
	) (core_domain.Warehouse, error)

	DeactivateWarehouse(
		ctx context.Context,
		warehouseID int64,
	) error

	ListStorageZones(
		ctx context.Context,
		warehouseID int64,
	) ([]core_domain.StorageZone, error)

	CreateStorageZone(
		ctx context.Context,
		zone core_domain.StorageZone,
	) (core_domain.StorageZone, error)

	PatchStorageZone(
		ctx context.Context,
		zoneID int64,
		patch core_domain.StorageZonePatch,
	) (core_domain.StorageZone, error)

	ListGates(
		ctx context.Context,
		warehouseID int64,
	) ([]core_domain.Gate, error)

	CreateGate(
		ctx context.Context,
		gate core_domain.Gate,
	) (core_domain.Gate, error)

	PatchGate(
		ctx context.Context,
		gateID int64,
		patch core_domain.GatePatch,
	) (core_domain.Gate, error)

	ListProductTypes(ctx context.Context) ([]core_domain.ProductType, error)

	ListCargoPlaceTypes(ctx context.Context) ([]core_domain.CargoPlaceType, error)
}

type WarehousesService struct {
	repo WarehousesRepository
}

func NewWarehousesService(repo WarehousesRepository) *WarehousesService {
	return &WarehousesService{
		repo: repo,
	}
}
