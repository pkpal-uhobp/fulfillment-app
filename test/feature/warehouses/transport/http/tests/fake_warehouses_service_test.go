package warehouses_transport_http_tests

import (
	"context"
	"fmt"
	"testing"

	core_logger "github.com/pkpal-uhobp/fulfillment-app/internal/core/logger"
	warehouses_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/warehouses/service"
	warehouses_http "github.com/pkpal-uhobp/fulfillment-app/internal/features/warehouses/transport/http"
)

type fakeWarehousesService struct {
	listWarehousesFn        func(context.Context, warehouses_service.WarehouseFilter) ([]warehouses_service.WarehouseDTO, error)
	getWarehouseFn          func(context.Context, int64) (warehouses_service.WarehouseDTO, error)
	createWarehouseFn       func(context.Context, warehouses_service.CreateWarehouseInput) (warehouses_service.WarehouseDTO, error)
	patchWarehouseFn        func(context.Context, int64, warehouses_service.PatchWarehouseInput) (warehouses_service.WarehouseDTO, error)
	activateWarehouseFn     func(context.Context, int64) error
	deactivateWarehouseFn   func(context.Context, int64) error
	listStorageZonesFn      func(context.Context, int64) ([]warehouses_service.StorageZoneDTO, error)
	createStorageZoneFn     func(context.Context, warehouses_service.CreateStorageZoneInput) (warehouses_service.StorageZoneDTO, error)
	patchStorageZoneFn      func(context.Context, int64, warehouses_service.PatchStorageZoneInput) (warehouses_service.StorageZoneDTO, error)
	activateStorageZoneFn   func(context.Context, int64) error
	deactivateStorageZoneFn func(context.Context, int64) error
	listGatesFn             func(context.Context, int64) ([]warehouses_service.GateDTO, error)
	createGateFn            func(context.Context, warehouses_service.CreateGateInput) (warehouses_service.GateDTO, error)
	patchGateFn             func(context.Context, int64, warehouses_service.PatchGateInput) (warehouses_service.GateDTO, error)
	activateGateFn          func(context.Context, int64) error
	deactivateGateFn        func(context.Context, int64) error
	listProductTypesFn      func(context.Context) ([]warehouses_service.ProductTypeDTO, error)
	listCargoPlaceTypesFn   func(context.Context) ([]warehouses_service.CargoPlaceTypeDTO, error)
}

func (f *fakeWarehousesService) ListWarehouses(ctx context.Context, filter warehouses_service.WarehouseFilter) ([]warehouses_service.WarehouseDTO, error) {
	if f.listWarehousesFn == nil {
		return nil, fmt.Errorf("unexpected ListWarehouses call")
	}
	return f.listWarehousesFn(ctx, filter)
}

func (f *fakeWarehousesService) GetWarehouse(ctx context.Context, warehouseID int64) (warehouses_service.WarehouseDTO, error) {
	if f.getWarehouseFn == nil {
		return warehouses_service.WarehouseDTO{}, fmt.Errorf("unexpected GetWarehouse call")
	}
	return f.getWarehouseFn(ctx, warehouseID)
}

func (f *fakeWarehousesService) CreateWarehouse(ctx context.Context, input warehouses_service.CreateWarehouseInput) (warehouses_service.WarehouseDTO, error) {
	if f.createWarehouseFn == nil {
		return warehouses_service.WarehouseDTO{}, fmt.Errorf("unexpected CreateWarehouse call")
	}
	return f.createWarehouseFn(ctx, input)
}

func (f *fakeWarehousesService) PatchWarehouse(ctx context.Context, warehouseID int64, input warehouses_service.PatchWarehouseInput) (warehouses_service.WarehouseDTO, error) {
	if f.patchWarehouseFn == nil {
		return warehouses_service.WarehouseDTO{}, fmt.Errorf("unexpected PatchWarehouse call")
	}
	return f.patchWarehouseFn(ctx, warehouseID, input)
}

func (f *fakeWarehousesService) ActivateWarehouse(ctx context.Context, warehouseID int64) error {
	if f.activateWarehouseFn == nil {
		return fmt.Errorf("unexpected ActivateWarehouse call")
	}
	return f.activateWarehouseFn(ctx, warehouseID)
}

func (f *fakeWarehousesService) DeactivateWarehouse(ctx context.Context, warehouseID int64) error {
	if f.deactivateWarehouseFn == nil {
		return fmt.Errorf("unexpected DeactivateWarehouse call")
	}
	return f.deactivateWarehouseFn(ctx, warehouseID)
}

func (f *fakeWarehousesService) ListStorageZones(ctx context.Context, warehouseID int64) ([]warehouses_service.StorageZoneDTO, error) {
	if f.listStorageZonesFn == nil {
		return nil, fmt.Errorf("unexpected ListStorageZones call")
	}
	return f.listStorageZonesFn(ctx, warehouseID)
}

func (f *fakeWarehousesService) CreateStorageZone(ctx context.Context, input warehouses_service.CreateStorageZoneInput) (warehouses_service.StorageZoneDTO, error) {
	if f.createStorageZoneFn == nil {
		return warehouses_service.StorageZoneDTO{}, fmt.Errorf("unexpected CreateStorageZone call")
	}
	return f.createStorageZoneFn(ctx, input)
}

func (f *fakeWarehousesService) PatchStorageZone(ctx context.Context, zoneID int64, input warehouses_service.PatchStorageZoneInput) (warehouses_service.StorageZoneDTO, error) {
	if f.patchStorageZoneFn == nil {
		return warehouses_service.StorageZoneDTO{}, fmt.Errorf("unexpected PatchStorageZone call")
	}
	return f.patchStorageZoneFn(ctx, zoneID, input)
}

func (f *fakeWarehousesService) ActivateStorageZone(ctx context.Context, zoneID int64) error {
	if f.activateStorageZoneFn == nil {
		return fmt.Errorf("unexpected ActivateStorageZone call")
	}
	return f.activateStorageZoneFn(ctx, zoneID)
}

func (f *fakeWarehousesService) DeactivateStorageZone(ctx context.Context, zoneID int64) error {
	if f.deactivateStorageZoneFn == nil {
		return fmt.Errorf("unexpected DeactivateStorageZone call")
	}
	return f.deactivateStorageZoneFn(ctx, zoneID)
}

func (f *fakeWarehousesService) ListGates(ctx context.Context, warehouseID int64) ([]warehouses_service.GateDTO, error) {
	if f.listGatesFn == nil {
		return nil, fmt.Errorf("unexpected ListGates call")
	}
	return f.listGatesFn(ctx, warehouseID)
}

func (f *fakeWarehousesService) CreateGate(ctx context.Context, input warehouses_service.CreateGateInput) (warehouses_service.GateDTO, error) {
	if f.createGateFn == nil {
		return warehouses_service.GateDTO{}, fmt.Errorf("unexpected CreateGate call")
	}
	return f.createGateFn(ctx, input)
}

func (f *fakeWarehousesService) PatchGate(ctx context.Context, gateID int64, input warehouses_service.PatchGateInput) (warehouses_service.GateDTO, error) {
	if f.patchGateFn == nil {
		return warehouses_service.GateDTO{}, fmt.Errorf("unexpected PatchGate call")
	}
	return f.patchGateFn(ctx, gateID, input)
}

func (f *fakeWarehousesService) ActivateGate(ctx context.Context, gateID int64) error {
	if f.activateGateFn == nil {
		return fmt.Errorf("unexpected ActivateGate call")
	}
	return f.activateGateFn(ctx, gateID)
}

func (f *fakeWarehousesService) DeactivateGate(ctx context.Context, gateID int64) error {
	if f.deactivateGateFn == nil {
		return fmt.Errorf("unexpected DeactivateGate call")
	}
	return f.deactivateGateFn(ctx, gateID)
}

func (f *fakeWarehousesService) ListProductTypes(ctx context.Context) ([]warehouses_service.ProductTypeDTO, error) {
	if f.listProductTypesFn == nil {
		return nil, fmt.Errorf("unexpected ListProductTypes call")
	}
	return f.listProductTypesFn(ctx)
}

func (f *fakeWarehousesService) ListCargoPlaceTypes(ctx context.Context) ([]warehouses_service.CargoPlaceTypeDTO, error) {
	if f.listCargoPlaceTypesFn == nil {
		return nil, fmt.Errorf("unexpected ListCargoPlaceTypes call")
	}
	return f.listCargoPlaceTypesFn(ctx)
}

func newTestHandler(t *testing.T, service *fakeWarehousesService) *warehouses_http.WarehousesHTTPHandler {
	t.Helper()

	log, err := core_logger.NewLogger(core_logger.LoggerConfig{
		Level:  "debug",
		Folder: t.TempDir(),
	})
	if err != nil {
		t.Fatalf("create logger: %v", err)
	}
	t.Cleanup(log.Close)

	return warehouses_http.NewWarehousesHTTPHandler(log, service)
}

func ptr[T any](value T) *T {
	return &value
}
