package warehouses_service_tests

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

type fakeWarehousesRepository struct {
	listWarehousesFn        func(context.Context, core_domain.WarehouseFilter) ([]core_domain.Warehouse, error)
	getWarehouseByIDFn      func(context.Context, int64) (core_domain.Warehouse, error)
	createWarehouseFn       func(context.Context, core_domain.Warehouse) (core_domain.Warehouse, error)
	patchWarehouseFn        func(context.Context, int64, core_domain.WarehousePatch) (core_domain.Warehouse, error)
	activateWarehouseFn     func(context.Context, int64) error
	deactivateWarehouseFn   func(context.Context, int64) error
	listStorageZonesFn      func(context.Context, int64) ([]core_domain.StorageZone, error)
	createStorageZoneFn     func(context.Context, core_domain.StorageZone) (core_domain.StorageZone, error)
	patchStorageZoneFn      func(context.Context, int64, core_domain.StorageZonePatch) (core_domain.StorageZone, error)
	activateStorageZoneFn   func(context.Context, int64) error
	deactivateStorageZoneFn func(context.Context, int64) error
	listGatesFn             func(context.Context, int64) ([]core_domain.Gate, error)
	createGateFn            func(context.Context, core_domain.Gate) (core_domain.Gate, error)
	patchGateFn             func(context.Context, int64, core_domain.GatePatch) (core_domain.Gate, error)
	activateGateFn          func(context.Context, int64) error
	deactivateGateFn        func(context.Context, int64) error
	listProductTypesFn      func(context.Context) ([]core_domain.ProductType, error)
	listCargoPlaceTypesFn   func(context.Context) ([]core_domain.CargoPlaceType, error)
}

func (f *fakeWarehousesRepository) ListWarehouses(ctx context.Context, filter core_domain.WarehouseFilter) ([]core_domain.Warehouse, error) {
	if f.listWarehousesFn == nil {
		return nil, fmt.Errorf("unexpected ListWarehouses call")
	}
	return f.listWarehousesFn(ctx, filter)
}

func (f *fakeWarehousesRepository) GetWarehouseByID(ctx context.Context, warehouseID int64) (core_domain.Warehouse, error) {
	if f.getWarehouseByIDFn == nil {
		return core_domain.Warehouse{}, fmt.Errorf("unexpected GetWarehouseByID call")
	}
	return f.getWarehouseByIDFn(ctx, warehouseID)
}

func (f *fakeWarehousesRepository) CreateWarehouse(ctx context.Context, warehouse core_domain.Warehouse) (core_domain.Warehouse, error) {
	if f.createWarehouseFn == nil {
		return core_domain.Warehouse{}, fmt.Errorf("unexpected CreateWarehouse call")
	}
	return f.createWarehouseFn(ctx, warehouse)
}

func (f *fakeWarehousesRepository) PatchWarehouse(ctx context.Context, warehouseID int64, patch core_domain.WarehousePatch) (core_domain.Warehouse, error) {
	if f.patchWarehouseFn == nil {
		return core_domain.Warehouse{}, fmt.Errorf("unexpected PatchWarehouse call")
	}
	return f.patchWarehouseFn(ctx, warehouseID, patch)
}

func (f *fakeWarehousesRepository) ActivateWarehouse(ctx context.Context, warehouseID int64) error {
	if f.activateWarehouseFn == nil {
		return fmt.Errorf("unexpected ActivateWarehouse call")
	}
	return f.activateWarehouseFn(ctx, warehouseID)
}

func (f *fakeWarehousesRepository) DeactivateWarehouse(ctx context.Context, warehouseID int64) error {
	if f.deactivateWarehouseFn == nil {
		return fmt.Errorf("unexpected DeactivateWarehouse call")
	}
	return f.deactivateWarehouseFn(ctx, warehouseID)
}

func (f *fakeWarehousesRepository) ListStorageZones(ctx context.Context, warehouseID int64) ([]core_domain.StorageZone, error) {
	if f.listStorageZonesFn == nil {
		return nil, fmt.Errorf("unexpected ListStorageZones call")
	}
	return f.listStorageZonesFn(ctx, warehouseID)
}

func (f *fakeWarehousesRepository) CreateStorageZone(ctx context.Context, zone core_domain.StorageZone) (core_domain.StorageZone, error) {
	if f.createStorageZoneFn == nil {
		return core_domain.StorageZone{}, fmt.Errorf("unexpected CreateStorageZone call")
	}
	return f.createStorageZoneFn(ctx, zone)
}

func (f *fakeWarehousesRepository) PatchStorageZone(ctx context.Context, zoneID int64, patch core_domain.StorageZonePatch) (core_domain.StorageZone, error) {
	if f.patchStorageZoneFn == nil {
		return core_domain.StorageZone{}, fmt.Errorf("unexpected PatchStorageZone call")
	}
	return f.patchStorageZoneFn(ctx, zoneID, patch)
}

func (f *fakeWarehousesRepository) ActivateStorageZone(ctx context.Context, zoneID int64) error {
	if f.activateStorageZoneFn == nil {
		return fmt.Errorf("unexpected ActivateStorageZone call")
	}
	return f.activateStorageZoneFn(ctx, zoneID)
}

func (f *fakeWarehousesRepository) DeactivateStorageZone(ctx context.Context, zoneID int64) error {
	if f.deactivateStorageZoneFn == nil {
		return fmt.Errorf("unexpected DeactivateStorageZone call")
	}
	return f.deactivateStorageZoneFn(ctx, zoneID)
}

func (f *fakeWarehousesRepository) ListGates(ctx context.Context, warehouseID int64) ([]core_domain.Gate, error) {
	if f.listGatesFn == nil {
		return nil, fmt.Errorf("unexpected ListGates call")
	}
	return f.listGatesFn(ctx, warehouseID)
}

func (f *fakeWarehousesRepository) CreateGate(ctx context.Context, gate core_domain.Gate) (core_domain.Gate, error) {
	if f.createGateFn == nil {
		return core_domain.Gate{}, fmt.Errorf("unexpected CreateGate call")
	}
	return f.createGateFn(ctx, gate)
}

func (f *fakeWarehousesRepository) PatchGate(ctx context.Context, gateID int64, patch core_domain.GatePatch) (core_domain.Gate, error) {
	if f.patchGateFn == nil {
		return core_domain.Gate{}, fmt.Errorf("unexpected PatchGate call")
	}
	return f.patchGateFn(ctx, gateID, patch)
}

func (f *fakeWarehousesRepository) ActivateGate(ctx context.Context, gateID int64) error {
	if f.activateGateFn == nil {
		return fmt.Errorf("unexpected ActivateGate call")
	}
	return f.activateGateFn(ctx, gateID)
}

func (f *fakeWarehousesRepository) DeactivateGate(ctx context.Context, gateID int64) error {
	if f.deactivateGateFn == nil {
		return fmt.Errorf("unexpected DeactivateGate call")
	}
	return f.deactivateGateFn(ctx, gateID)
}

func (f *fakeWarehousesRepository) ListProductTypes(ctx context.Context) ([]core_domain.ProductType, error) {
	if f.listProductTypesFn == nil {
		return nil, fmt.Errorf("unexpected ListProductTypes call")
	}
	return f.listProductTypesFn(ctx)
}

func (f *fakeWarehousesRepository) ListCargoPlaceTypes(ctx context.Context) ([]core_domain.CargoPlaceType, error) {
	if f.listCargoPlaceTypesFn == nil {
		return nil, fmt.Errorf("unexpected ListCargoPlaceTypes call")
	}
	return f.listCargoPlaceTypesFn(ctx)
}

func ptr[T any](value T) *T {
	return &value
}
