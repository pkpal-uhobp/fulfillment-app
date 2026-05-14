package warehouses_service_tests

import (
	"context"
	"errors"
	"testing"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
	warehouses_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/warehouses/service"
)

func TestListWarehousesSuccessTrimsFilterAndMapsDTO(t *testing.T) {
	marketplace := "WB"
	repo := &fakeWarehousesRepository{
		listWarehousesFn: func(ctx context.Context, filter core_domain.WarehouseFilter) ([]core_domain.Warehouse, error) {
			if filter.WarehouseType != core_domain.WarehouseTypeReceiving.String() {
				t.Fatalf("WarehouseType = %q, want receiving", filter.WarehouseType)
			}
			if filter.Marketplace != "WB" {
				t.Fatalf("Marketplace = %q, want WB", filter.Marketplace)
			}
			if filter.City != "Moscow" {
				t.Fatalf("City = %q, want Moscow", filter.City)
			}

			return []core_domain.Warehouse{
				{
					ID:            10,
					Name:          "Приёмный склад",
					WarehouseType: core_domain.WarehouseTypeReceiving,
					Marketplace:   &marketplace,
					City:          "Moscow",
					Address:       "Test address",
					IsActive:      true,
				},
			}, nil
		},
	}

	service := warehouses_service.NewWarehousesService(repo)

	warehouses, err := service.ListWarehouses(context.Background(), warehouses_service.WarehouseFilter{
		WarehouseType: " receiving ",
		Marketplace:   " WB ",
		City:          " Moscow ",
	})

	if err != nil {
		t.Fatalf("ListWarehouses returned error: %v", err)
	}
	if len(warehouses) != 1 {
		t.Fatalf("len = %d, want 1", len(warehouses))
	}
	if warehouses[0].ID != 10 {
		t.Fatalf("ID = %d, want 10", warehouses[0].ID)
	}
	if warehouses[0].Marketplace != "WB" {
		t.Fatalf("Marketplace = %q, want WB", warehouses[0].Marketplace)
	}
}

func TestListWarehousesRejectsInvalidWarehouseType(t *testing.T) {
	service := warehouses_service.NewWarehousesService(&fakeWarehousesRepository{})

	_, err := service.ListWarehouses(context.Background(), warehouses_service.WarehouseFilter{
		WarehouseType: "wrong_type",
	})

	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("err = %v, want ErrInvalidArgument", err)
	}
}

func TestGetWarehouseRejectsInvalidID(t *testing.T) {
	service := warehouses_service.NewWarehousesService(&fakeWarehousesRepository{})

	_, err := service.GetWarehouse(context.Background(), 0)

	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("err = %v, want ErrInvalidArgument", err)
	}
}

func TestCreateWarehouseSuccess(t *testing.T) {
	repo := &fakeWarehousesRepository{
		createWarehouseFn: func(ctx context.Context, warehouse core_domain.Warehouse) (core_domain.Warehouse, error) {
			if warehouse.Name != "Склад Москва" {
				t.Fatalf("Name = %q, want Склад Москва", warehouse.Name)
			}
			if warehouse.WarehouseType != core_domain.WarehouseTypeBoth {
				t.Fatalf("WarehouseType = %q, want both", warehouse.WarehouseType)
			}
			if warehouse.Marketplace == nil || *warehouse.Marketplace != "Ozon" {
				t.Fatalf("Marketplace = %v, want Ozon", warehouse.Marketplace)
			}
			if !warehouse.IsActive {
				t.Fatalf("new warehouse must be active")
			}

			warehouse.ID = 11
			return warehouse, nil
		},
	}

	service := warehouses_service.NewWarehousesService(repo)

	dto, err := service.CreateWarehouse(context.Background(), warehouses_service.CreateWarehouseInput{
		Name:          "  Склад Москва  ",
		WarehouseType: " both ",
		Marketplace:   " Ozon ",
		City:          " Москва ",
		Address:       " ул. Тестовая, 1 ",
	})

	if err != nil {
		t.Fatalf("CreateWarehouse returned error: %v", err)
	}
	if dto.ID != 11 {
		t.Fatalf("ID = %d, want 11", dto.ID)
	}
	if dto.WarehouseType != core_domain.WarehouseTypeBoth.String() {
		t.Fatalf("WarehouseType = %q, want both", dto.WarehouseType)
	}
	if !dto.IsActive {
		t.Fatalf("dto.IsActive must be true")
	}
}

func TestCreateWarehouseRejectsInvalidInput(t *testing.T) {
	service := warehouses_service.NewWarehousesService(&fakeWarehousesRepository{})

	_, err := service.CreateWarehouse(context.Background(), warehouses_service.CreateWarehouseInput{
		Name:          "",
		WarehouseType: "wrong_type",
		City:          "Москва",
		Address:       "Адрес",
	})

	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("err = %v, want ErrInvalidArgument", err)
	}
}

func TestPatchWarehouseSuccess(t *testing.T) {
	isActive := false
	repo := &fakeWarehousesRepository{
		patchWarehouseFn: func(ctx context.Context, warehouseID int64, patch core_domain.WarehousePatch) (core_domain.Warehouse, error) {
			if warehouseID != 15 {
				t.Fatalf("warehouseID = %d, want 15", warehouseID)
			}
			if patch.Name == nil || *patch.Name != "Новый склад" {
				t.Fatalf("patch.Name = %v, want Новый склад", patch.Name)
			}
			if patch.City == nil || *patch.City != "Казань" {
				t.Fatalf("patch.City = %v, want Казань", patch.City)
			}
			if patch.IsActive == nil || *patch.IsActive != false {
				t.Fatalf("patch.IsActive = %v, want false", patch.IsActive)
			}

			return core_domain.Warehouse{
				ID:            warehouseID,
				Name:          *patch.Name,
				WarehouseType: core_domain.WarehouseTypeDestination,
				City:          *patch.City,
				Address:       "Адрес",
				IsActive:      *patch.IsActive,
			}, nil
		},
	}

	service := warehouses_service.NewWarehousesService(repo)

	dto, err := service.PatchWarehouse(context.Background(), 15, warehouses_service.PatchWarehouseInput{
		Name:     ptr(" Новый склад "),
		City:     ptr(" Казань "),
		IsActive: &isActive,
	})

	if err != nil {
		t.Fatalf("PatchWarehouse returned error: %v", err)
	}
	if dto.Name != "Новый склад" {
		t.Fatalf("Name = %q, want Новый склад", dto.Name)
	}
	if dto.IsActive {
		t.Fatalf("IsActive = true, want false")
	}
}

func TestPatchWarehouseRejectsEmptyName(t *testing.T) {
	service := warehouses_service.NewWarehousesService(&fakeWarehousesRepository{})

	_, err := service.PatchWarehouse(context.Background(), 1, warehouses_service.PatchWarehouseInput{
		Name: ptr("   "),
	})

	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("err = %v, want ErrInvalidArgument", err)
	}
}

func TestCreateStorageZoneSuccess(t *testing.T) {
	repo := &fakeWarehousesRepository{
		createStorageZoneFn: func(ctx context.Context, zone core_domain.StorageZone) (core_domain.StorageZone, error) {
			if zone.WarehouseID != 7 {
				t.Fatalf("WarehouseID = %d, want 7", zone.WarehouseID)
			}
			if zone.Name != "A-01" {
				t.Fatalf("Name = %q, want A-01", zone.Name)
			}
			if zone.Description == nil || *zone.Description != "Хранение коробок" {
				t.Fatalf("Description = %v, want Хранение коробок", zone.Description)
			}
			zone.ID = 100
			return zone, nil
		},
	}

	service := warehouses_service.NewWarehousesService(repo)

	dto, err := service.CreateStorageZone(context.Background(), warehouses_service.CreateStorageZoneInput{
		WarehouseID: 7,
		Name:        " A-01 ",
		Description: " Хранение коробок ",
	})

	if err != nil {
		t.Fatalf("CreateStorageZone returned error: %v", err)
	}
	if dto.ID != 100 {
		t.Fatalf("ID = %d, want 100", dto.ID)
	}
	if !dto.IsActive {
		t.Fatalf("storage zone must be active")
	}
}

func TestCreateStorageZoneRejectsInvalidWarehouseID(t *testing.T) {
	service := warehouses_service.NewWarehousesService(&fakeWarehousesRepository{})

	_, err := service.CreateStorageZone(context.Background(), warehouses_service.CreateStorageZoneInput{
		WarehouseID: 0,
		Name:        "A-01",
	})

	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("err = %v, want ErrInvalidArgument", err)
	}
}

func TestPatchStorageZoneRejectsEmptyName(t *testing.T) {
	service := warehouses_service.NewWarehousesService(&fakeWarehousesRepository{})

	_, err := service.PatchStorageZone(context.Background(), 1, warehouses_service.PatchStorageZoneInput{
		Name: ptr(" "),
	})

	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("err = %v, want ErrInvalidArgument", err)
	}
}

func TestCreateGateSuccess(t *testing.T) {
	repo := &fakeWarehousesRepository{
		createGateFn: func(ctx context.Context, gate core_domain.Gate) (core_domain.Gate, error) {
			if gate.WarehouseID != 8 {
				t.Fatalf("WarehouseID = %d, want 8", gate.WarehouseID)
			}
			if gate.Name != "Gate-1" {
				t.Fatalf("Name = %q, want Gate-1", gate.Name)
			}
			gate.ID = 200
			return gate, nil
		},
	}

	service := warehouses_service.NewWarehousesService(repo)

	dto, err := service.CreateGate(context.Background(), warehouses_service.CreateGateInput{
		WarehouseID: 8,
		Name:        " Gate-1 ",
	})

	if err != nil {
		t.Fatalf("CreateGate returned error: %v", err)
	}
	if dto.ID != 200 {
		t.Fatalf("ID = %d, want 200", dto.ID)
	}
}

func TestCreateGateRejectsInvalidInput(t *testing.T) {
	service := warehouses_service.NewWarehousesService(&fakeWarehousesRepository{})

	_, err := service.CreateGate(context.Background(), warehouses_service.CreateGateInput{
		WarehouseID: 0,
		Name:        "Gate-1",
	})

	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("err = %v, want ErrInvalidArgument", err)
	}
}

func TestCatalogMethodsReturnDTOs(t *testing.T) {
	repo := &fakeWarehousesRepository{
		listProductTypesFn: func(ctx context.Context) ([]core_domain.ProductType, error) {
			return []core_domain.ProductType{{ID: 1, Name: "Одежда", IsActive: true}}, nil
		},
		listCargoPlaceTypesFn: func(ctx context.Context) ([]core_domain.CargoPlaceType, error) {
			return []core_domain.CargoPlaceType{{ID: 2, Name: "Коробка", IsActive: true}}, nil
		},
	}

	service := warehouses_service.NewWarehousesService(repo)

	productTypes, err := service.ListProductTypes(context.Background())
	if err != nil {
		t.Fatalf("ListProductTypes returned error: %v", err)
	}
	if len(productTypes) != 1 || productTypes[0].Name != "Одежда" {
		t.Fatalf("unexpected productTypes: %+v", productTypes)
	}

	cargoPlaceTypes, err := service.ListCargoPlaceTypes(context.Background())
	if err != nil {
		t.Fatalf("ListCargoPlaceTypes returned error: %v", err)
	}
	if len(cargoPlaceTypes) != 1 || cargoPlaceTypes[0].Name != "Коробка" {
		t.Fatalf("unexpected cargoPlaceTypes: %+v", cargoPlaceTypes)
	}
}

func TestStatusMethodsValidateIDsAndCallRepository(t *testing.T) {
	calls := map[string]int{}
	repo := &fakeWarehousesRepository{
		activateWarehouseFn: func(ctx context.Context, warehouseID int64) error {
			if warehouseID != 1 {
				t.Fatalf("warehouseID = %d, want 1", warehouseID)
			}
			calls["activateWarehouse"]++
			return nil
		},
		deactivateWarehouseFn: func(ctx context.Context, warehouseID int64) error {
			if warehouseID != 1 {
				t.Fatalf("warehouseID = %d, want 1", warehouseID)
			}
			calls["deactivateWarehouse"]++
			return nil
		},
		activateStorageZoneFn: func(ctx context.Context, zoneID int64) error {
			if zoneID != 2 {
				t.Fatalf("zoneID = %d, want 2", zoneID)
			}
			calls["activateStorageZone"]++
			return nil
		},
		deactivateStorageZoneFn: func(ctx context.Context, zoneID int64) error {
			if zoneID != 2 {
				t.Fatalf("zoneID = %d, want 2", zoneID)
			}
			calls["deactivateStorageZone"]++
			return nil
		},
		activateGateFn: func(ctx context.Context, gateID int64) error {
			if gateID != 3 {
				t.Fatalf("gateID = %d, want 3", gateID)
			}
			calls["activateGate"]++
			return nil
		},
		deactivateGateFn: func(ctx context.Context, gateID int64) error {
			if gateID != 3 {
				t.Fatalf("gateID = %d, want 3", gateID)
			}
			calls["deactivateGate"]++
			return nil
		},
	}

	service := warehouses_service.NewWarehousesService(repo)

	if err := service.ActivateWarehouse(context.Background(), 1); err != nil {
		t.Fatalf("ActivateWarehouse returned error: %v", err)
	}
	if err := service.DeactivateWarehouse(context.Background(), 1); err != nil {
		t.Fatalf("DeactivateWarehouse returned error: %v", err)
	}
	if err := service.ActivateStorageZone(context.Background(), 2); err != nil {
		t.Fatalf("ActivateStorageZone returned error: %v", err)
	}
	if err := service.DeactivateStorageZone(context.Background(), 2); err != nil {
		t.Fatalf("DeactivateStorageZone returned error: %v", err)
	}
	if err := service.ActivateGate(context.Background(), 3); err != nil {
		t.Fatalf("ActivateGate returned error: %v", err)
	}
	if err := service.DeactivateGate(context.Background(), 3); err != nil {
		t.Fatalf("DeactivateGate returned error: %v", err)
	}

	if len(calls) != 6 {
		t.Fatalf("calls len = %d, want 6; calls = %+v", len(calls), calls)
	}

	invalidChecks := []struct {
		name string
		fn   func() error
	}{
		{"ActivateWarehouse", func() error { return service.ActivateWarehouse(context.Background(), 0) }},
		{"DeactivateWarehouse", func() error { return service.DeactivateWarehouse(context.Background(), 0) }},
		{"ActivateStorageZone", func() error { return service.ActivateStorageZone(context.Background(), 0) }},
		{"DeactivateStorageZone", func() error { return service.DeactivateStorageZone(context.Background(), 0) }},
		{"ActivateGate", func() error { return service.ActivateGate(context.Background(), 0) }},
		{"DeactivateGate", func() error { return service.DeactivateGate(context.Background(), 0) }},
	}

	for _, tt := range invalidChecks {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn()
			if !errors.Is(err, core_errors.ErrInvalidArgument) {
				t.Fatalf("err = %v, want ErrInvalidArgument", err)
			}
		})
	}
}
