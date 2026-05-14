package warehouses_transport_http_tests

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	warehouses_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/warehouses/service"
)

func TestStorageZoneHandlers(t *testing.T) {
	called := map[string]int{}
	service := &fakeWarehousesService{
		listStorageZonesFn: func(_ context.Context, warehouseID int64) ([]warehouses_service.StorageZoneDTO, error) {
			if warehouseID != 5 {
				t.Fatalf("warehouseID = %d, want 5", warehouseID)
			}
			called["list"]++
			return []warehouses_service.StorageZoneDTO{{ID: 1, WarehouseID: 5, Name: "A-01", IsActive: true}}, nil
		},
		createStorageZoneFn: func(_ context.Context, input warehouses_service.CreateStorageZoneInput) (warehouses_service.StorageZoneDTO, error) {
			if input.WarehouseID != 5 || input.Name != "A-02" {
				t.Fatalf("input = %+v, want warehouse_id=5 name=A-02", input)
			}
			called["create"]++
			return warehouses_service.StorageZoneDTO{ID: 2, WarehouseID: input.WarehouseID, Name: input.Name, Description: input.Description, IsActive: true}, nil
		},
		patchStorageZoneFn: func(_ context.Context, zoneID int64, input warehouses_service.PatchStorageZoneInput) (warehouses_service.StorageZoneDTO, error) {
			if zoneID != 2 {
				t.Fatalf("zoneID = %d, want 2", zoneID)
			}
			if input.Name == nil || *input.Name != "A-03" {
				t.Fatalf("Name = %v, want A-03", input.Name)
			}
			called["patch"]++
			return warehouses_service.StorageZoneDTO{ID: zoneID, WarehouseID: 5, Name: *input.Name, IsActive: true}, nil
		},
		activateStorageZoneFn: func(_ context.Context, zoneID int64) error {
			if zoneID != 2 {
				t.Fatalf("zoneID = %d, want 2", zoneID)
			}
			called["activate"]++
			return nil
		},
		deactivateStorageZoneFn: func(_ context.Context, zoneID int64) error {
			if zoneID != 2 {
				t.Fatalf("zoneID = %d, want 2", zoneID)
			}
			called["deactivate"]++
			return nil
		},
	}
	handler := newTestHandler(t, service)

	listReq := httptest.NewRequest(http.MethodGet, "/storage-zones?warehouse_id=5", nil)
	listRec := httptest.NewRecorder()
	handler.ListStorageZones(listRec, listReq)
	if listRec.Code != http.StatusOK {
		t.Fatalf("list status = %d, want %d; body=%s", listRec.Code, http.StatusOK, listRec.Body.String())
	}

	createReq := httptest.NewRequest(http.MethodPost, "/storage-zones", bytes.NewReader([]byte(`{"warehouse_id":5,"name":"A-02","description":"Зона"}`)))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	handler.CreateStorageZone(createRec, createReq)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d; body=%s", createRec.Code, http.StatusCreated, createRec.Body.String())
	}

	patchReq := httptest.NewRequest(http.MethodPatch, "/storage-zones/2", bytes.NewReader([]byte(`{"name":"A-03"}`)))
	patchReq.Header.Set("Content-Type", "application/json")
	patchReq.SetPathValue("id", "2")
	patchRec := httptest.NewRecorder()
	handler.PatchStorageZone(patchRec, patchReq)
	if patchRec.Code != http.StatusOK {
		t.Fatalf("patch status = %d, want %d; body=%s", patchRec.Code, http.StatusOK, patchRec.Body.String())
	}

	activateReq := httptest.NewRequest(http.MethodPatch, "/storage-zones/2/activate", nil)
	activateReq.SetPathValue("id", "2")
	activateRec := httptest.NewRecorder()
	handler.ActivateStorageZone(activateRec, activateReq)
	if activateRec.Code != http.StatusNoContent {
		t.Fatalf("activate status = %d, want %d", activateRec.Code, http.StatusNoContent)
	}

	deactivateReq := httptest.NewRequest(http.MethodDelete, "/storage-zones/2", nil)
	deactivateReq.SetPathValue("id", "2")
	deactivateRec := httptest.NewRecorder()
	handler.DeactivateStorageZone(deactivateRec, deactivateReq)
	if deactivateRec.Code != http.StatusNoContent {
		t.Fatalf("deactivate status = %d, want %d", deactivateRec.Code, http.StatusNoContent)
	}

	if len(called) != 5 {
		t.Fatalf("called len = %d, want 5; called=%+v", len(called), called)
	}
}

func TestGateHandlers(t *testing.T) {
	called := map[string]int{}
	service := &fakeWarehousesService{
		listGatesFn: func(_ context.Context, warehouseID int64) ([]warehouses_service.GateDTO, error) {
			if warehouseID != 5 {
				t.Fatalf("warehouseID = %d, want 5", warehouseID)
			}
			called["list"]++
			return []warehouses_service.GateDTO{{ID: 1, WarehouseID: 5, Name: "Gate-1", IsActive: true}}, nil
		},
		createGateFn: func(_ context.Context, input warehouses_service.CreateGateInput) (warehouses_service.GateDTO, error) {
			if input.WarehouseID != 5 || input.Name != "Gate-2" {
				t.Fatalf("input = %+v, want warehouse_id=5 name=Gate-2", input)
			}
			called["create"]++
			return warehouses_service.GateDTO{ID: 2, WarehouseID: input.WarehouseID, Name: input.Name, IsActive: true}, nil
		},
		patchGateFn: func(_ context.Context, gateID int64, input warehouses_service.PatchGateInput) (warehouses_service.GateDTO, error) {
			if gateID != 2 {
				t.Fatalf("gateID = %d, want 2", gateID)
			}
			if input.Name == nil || *input.Name != "Gate-3" {
				t.Fatalf("Name = %v, want Gate-3", input.Name)
			}
			called["patch"]++
			return warehouses_service.GateDTO{ID: gateID, WarehouseID: 5, Name: *input.Name, IsActive: true}, nil
		},
		activateGateFn: func(_ context.Context, gateID int64) error {
			if gateID != 2 {
				t.Fatalf("gateID = %d, want 2", gateID)
			}
			called["activate"]++
			return nil
		},
		deactivateGateFn: func(_ context.Context, gateID int64) error {
			if gateID != 2 {
				t.Fatalf("gateID = %d, want 2", gateID)
			}
			called["deactivate"]++
			return nil
		},
	}
	handler := newTestHandler(t, service)

	listReq := httptest.NewRequest(http.MethodGet, "/gates?warehouse_id=5", nil)
	listRec := httptest.NewRecorder()
	handler.ListGates(listRec, listReq)
	if listRec.Code != http.StatusOK {
		t.Fatalf("list status = %d, want %d; body=%s", listRec.Code, http.StatusOK, listRec.Body.String())
	}

	createReq := httptest.NewRequest(http.MethodPost, "/gates", bytes.NewReader([]byte(`{"warehouse_id":5,"name":"Gate-2"}`)))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	handler.CreateGate(createRec, createReq)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d; body=%s", createRec.Code, http.StatusCreated, createRec.Body.String())
	}

	patchReq := httptest.NewRequest(http.MethodPatch, "/gates/2", bytes.NewReader([]byte(`{"name":"Gate-3"}`)))
	patchReq.Header.Set("Content-Type", "application/json")
	patchReq.SetPathValue("id", "2")
	patchRec := httptest.NewRecorder()
	handler.PatchGate(patchRec, patchReq)
	if patchRec.Code != http.StatusOK {
		t.Fatalf("patch status = %d, want %d; body=%s", patchRec.Code, http.StatusOK, patchRec.Body.String())
	}

	activateReq := httptest.NewRequest(http.MethodPatch, "/gates/2/activate", nil)
	activateReq.SetPathValue("id", "2")
	activateRec := httptest.NewRecorder()
	handler.ActivateGate(activateRec, activateReq)
	if activateRec.Code != http.StatusNoContent {
		t.Fatalf("activate status = %d, want %d", activateRec.Code, http.StatusNoContent)
	}

	deactivateReq := httptest.NewRequest(http.MethodDelete, "/gates/2", nil)
	deactivateReq.SetPathValue("id", "2")
	deactivateRec := httptest.NewRecorder()
	handler.DeactivateGate(deactivateRec, deactivateReq)
	if deactivateRec.Code != http.StatusNoContent {
		t.Fatalf("deactivate status = %d, want %d", deactivateRec.Code, http.StatusNoContent)
	}

	if len(called) != 5 {
		t.Fatalf("called len = %d, want 5; called=%+v", len(called), called)
	}
}

func TestCatalogHandlersReturnOK(t *testing.T) {
	service := &fakeWarehousesService{
		listProductTypesFn: func(_ context.Context) ([]warehouses_service.ProductTypeDTO, error) {
			return []warehouses_service.ProductTypeDTO{{ID: 1, Name: "Одежда", IsActive: true}}, nil
		},
		listCargoPlaceTypesFn: func(_ context.Context) ([]warehouses_service.CargoPlaceTypeDTO, error) {
			return []warehouses_service.CargoPlaceTypeDTO{{ID: 2, Name: "Коробка", IsActive: true}}, nil
		},
	}
	handler := newTestHandler(t, service)

	productReq := httptest.NewRequest(http.MethodGet, "/product-types", nil)
	productRec := httptest.NewRecorder()
	handler.ListProductTypes(productRec, productReq)
	if productRec.Code != http.StatusOK {
		t.Fatalf("product types status = %d, want %d", productRec.Code, http.StatusOK)
	}

	cargoReq := httptest.NewRequest(http.MethodGet, "/cargo-place-types", nil)
	cargoRec := httptest.NewRecorder()
	handler.ListCargoPlaceTypes(cargoRec, cargoReq)
	if cargoRec.Code != http.StatusOK {
		t.Fatalf("cargo place types status = %d, want %d", cargoRec.Code, http.StatusOK)
	}
}
