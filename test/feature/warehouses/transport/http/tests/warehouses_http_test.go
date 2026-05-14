package warehouses_transport_http_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
	warehouses_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/warehouses/service"
)

func decodeJSON(t *testing.T, recorder *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var body map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v; body=%s", err, recorder.Body.String())
	}
	return body
}

func TestListWarehousesReturnsOKAndPassesFilters(t *testing.T) {
	service := &fakeWarehousesService{
		listWarehousesFn: func(_ context.Context, filter warehouses_service.WarehouseFilter) ([]warehouses_service.WarehouseDTO, error) {
			if filter.WarehouseType != "receiving" {
				t.Fatalf("WarehouseType = %q, want receiving", filter.WarehouseType)
			}
			if filter.Marketplace != "WB" {
				t.Fatalf("Marketplace = %q, want WB", filter.Marketplace)
			}
			if filter.City != "Moscow" {
				t.Fatalf("City = %q, want Moscow", filter.City)
			}
			return []warehouses_service.WarehouseDTO{
				{ID: 1, Name: "Склад 1", WarehouseType: "receiving", Marketplace: "WB", City: "Moscow", Address: "Address", IsActive: true},
			}, nil
		},
	}
	handler := newTestHandler(t, service)

	req := httptest.NewRequest(http.MethodGet, "/warehouses?warehouse_type=receiving&marketplace=WB&city=Moscow", nil)
	rec := httptest.NewRecorder()

	handler.ListWarehouses(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	body := decodeJSON(t, rec)
	warehouses := body["warehouses"].([]any)
	if len(warehouses) != 1 {
		t.Fatalf("warehouses len = %d, want 1", len(warehouses))
	}
}

func TestGetWarehouseReturnsOK(t *testing.T) {
	service := &fakeWarehousesService{
		getWarehouseFn: func(_ context.Context, warehouseID int64) (warehouses_service.WarehouseDTO, error) {
			if warehouseID != 7 {
				t.Fatalf("warehouseID = %d, want 7", warehouseID)
			}
			return warehouses_service.WarehouseDTO{ID: warehouseID, Name: "Склад", WarehouseType: "both", City: "Москва", Address: "Адрес", IsActive: true}, nil
		},
	}
	handler := newTestHandler(t, service)

	req := httptest.NewRequest(http.MethodGet, "/warehouses/7", nil)
	req.SetPathValue("id", "7")
	rec := httptest.NewRecorder()

	handler.GetWarehouse(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
}

func TestGetWarehouseReturnsBadRequestForInvalidID(t *testing.T) {
	handler := newTestHandler(t, &fakeWarehousesService{})

	req := httptest.NewRequest(http.MethodGet, "/warehouses/bad", nil)
	req.SetPathValue("id", "bad")
	rec := httptest.NewRecorder()

	handler.GetWarehouse(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}

func TestCreateWarehouseReturnsCreatedAndPassesInput(t *testing.T) {
	service := &fakeWarehousesService{
		createWarehouseFn: func(_ context.Context, input warehouses_service.CreateWarehouseInput) (warehouses_service.WarehouseDTO, error) {
			if input.Name != "Склад Москва" {
				t.Fatalf("Name = %q, want Склад Москва", input.Name)
			}
			if input.WarehouseType != "both" {
				t.Fatalf("WarehouseType = %q, want both", input.WarehouseType)
			}
			return warehouses_service.WarehouseDTO{ID: 10, Name: input.Name, WarehouseType: input.WarehouseType, City: input.City, Address: input.Address, IsActive: true}, nil
		},
	}
	handler := newTestHandler(t, service)

	body := []byte(`{"name":"Склад Москва","warehouse_type":"both","marketplace":"Ozon","city":"Москва","address":"ул. Тестовая, 1"}`)
	req := httptest.NewRequest(http.MethodPost, "/warehouses", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.CreateWarehouse(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusCreated, rec.Body.String())
	}
}

func TestCreateWarehouseReturnsBadRequestForInvalidJSON(t *testing.T) {
	handler := newTestHandler(t, &fakeWarehousesService{})

	req := httptest.NewRequest(http.MethodPost, "/warehouses", bytes.NewReader([]byte(`{"name":""}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.CreateWarehouse(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}

func TestPatchWarehouseReturnsOK(t *testing.T) {
	service := &fakeWarehousesService{
		patchWarehouseFn: func(_ context.Context, warehouseID int64, input warehouses_service.PatchWarehouseInput) (warehouses_service.WarehouseDTO, error) {
			if warehouseID != 10 {
				t.Fatalf("warehouseID = %d, want 10", warehouseID)
			}
			if input.Name == nil || *input.Name != "Обновленный склад" {
				t.Fatalf("Name = %v, want Обновленный склад", input.Name)
			}
			if input.IsActive == nil || *input.IsActive != false {
				t.Fatalf("IsActive = %v, want false", input.IsActive)
			}
			return warehouses_service.WarehouseDTO{ID: warehouseID, Name: *input.Name, WarehouseType: "both", City: "Москва", Address: "Адрес", IsActive: *input.IsActive}, nil
		},
	}
	handler := newTestHandler(t, service)

	body := []byte(`{"name":"Обновленный склад","is_active":false}`)
	req := httptest.NewRequest(http.MethodPatch, "/warehouses/10", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "10")
	rec := httptest.NewRecorder()

	handler.PatchWarehouse(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
}

func TestWarehouseStatusHandlersReturnNoContent(t *testing.T) {
	called := map[string]int{}
	service := &fakeWarehousesService{
		activateWarehouseFn: func(_ context.Context, warehouseID int64) error {
			if warehouseID != 10 {
				t.Fatalf("warehouseID = %d, want 10", warehouseID)
			}
			called["activate"]++
			return nil
		},
		deactivateWarehouseFn: func(_ context.Context, warehouseID int64) error {
			if warehouseID != 10 {
				t.Fatalf("warehouseID = %d, want 10", warehouseID)
			}
			called["deactivate"]++
			return nil
		},
	}
	handler := newTestHandler(t, service)

	activateReq := httptest.NewRequest(http.MethodPatch, "/warehouses/10/activate", nil)
	activateReq.SetPathValue("id", "10")
	activateRec := httptest.NewRecorder()
	handler.ActivateWarehouse(activateRec, activateReq)
	if activateRec.Code != http.StatusNoContent {
		t.Fatalf("activate status = %d, want %d", activateRec.Code, http.StatusNoContent)
	}

	deactivateReq := httptest.NewRequest(http.MethodDelete, "/warehouses/10", nil)
	deactivateReq.SetPathValue("id", "10")
	deactivateRec := httptest.NewRecorder()
	handler.DeactivateWarehouse(deactivateRec, deactivateReq)
	if deactivateRec.Code != http.StatusNoContent {
		t.Fatalf("deactivate status = %d, want %d", deactivateRec.Code, http.StatusNoContent)
	}

	if called["activate"] != 1 || called["deactivate"] != 1 {
		t.Fatalf("unexpected calls: %+v", called)
	}
}

func TestCreateWarehouseReturnsConflictFromService(t *testing.T) {
	service := &fakeWarehousesService{
		createWarehouseFn: func(_ context.Context, input warehouses_service.CreateWarehouseInput) (warehouses_service.WarehouseDTO, error) {
			return warehouses_service.WarehouseDTO{}, core_errors.ErrConflict
		},
	}
	handler := newTestHandler(t, service)

	body := []byte(`{"name":"Склад Москва","warehouse_type":"both","city":"Москва","address":"Адрес"}`)
	req := httptest.NewRequest(http.MethodPost, "/warehouses", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.CreateWarehouse(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusConflict, rec.Body.String())
	}
}
