package cargoitems_transport_http_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_logger "github.com/pkpal-uhobp/fulfillment-app/internal/core/logger"
	core_http_middleware "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/middleware"
	cargoitems_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/cargoitems/service"
	cargoitems_transport_http "github.com/pkpal-uhobp/fulfillment-app/internal/features/cargoitems/transport/http"
)

func TestRoutesExposeCargoItemsEndpoints(t *testing.T) {
	handler := cargoitems_transport_http.NewCargoItemsHTTPHandler(newTestLogger(t), &fakeCargoItemsService{})
	routes := handler.Routes()

	expected := map[string][]string{
		"POST /orders/{id}/cargo-items":       {core_domain.RoleWorker.String(), core_domain.RoleAdmin.String()},
		"GET /cargo-items":                    {core_domain.RoleClient.String(), core_domain.RoleWorker.String(), core_domain.RoleLogist.String(), core_domain.RoleAdmin.String()},
		"GET /cargo-items/scan":               {core_domain.RoleWorker.String(), core_domain.RoleLogist.String(), core_domain.RoleAdmin.String()},
		"GET /cargo-items/{id}":               {core_domain.RoleClient.String(), core_domain.RoleWorker.String(), core_domain.RoleLogist.String(), core_domain.RoleAdmin.String()},
		"GET /cargo-items/{id}/label":         {core_domain.RoleWorker.String(), core_domain.RoleLogist.String(), core_domain.RoleAdmin.String()},
		"GET /cargo-items/{id}/history":       {core_domain.RoleClient.String(), core_domain.RoleWorker.String(), core_domain.RoleLogist.String(), core_domain.RoleAdmin.String()},
		"PATCH /cargo-items/{id}/status":      {core_domain.RoleWorker.String(), core_domain.RoleLogist.String(), core_domain.RoleAdmin.String()},
		"PATCH /cargo-items/{id}/assign-zone": {core_domain.RoleLogist.String(), core_domain.RoleAdmin.String()},
		"PATCH /cargo-items/{id}/assign-gate": {core_domain.RoleLogist.String(), core_domain.RoleAdmin.String()},
	}

	if len(routes) != len(expected) {
		t.Fatalf("len(routes) = %d, want %d", len(routes), len(expected))
	}
	for _, route := range routes {
		key := route.Method + " " + route.Path
		roles, ok := expected[key]
		if !ok {
			t.Fatalf("unexpected route: %s", key)
		}
		if len(route.Roles) != len(roles) {
			t.Fatalf("roles len for %s = %d, want %d", key, len(route.Roles), len(roles))
		}
		for i := range roles {
			if route.Roles[i] != roles[i] {
				t.Fatalf("role[%d] for %s = %q, want %q", i, key, route.Roles[i], roles[i])
			}
		}
	}
}

func TestCreateCargoItemHTTP(t *testing.T) {
	now := time.Now().UTC()
	service := &fakeCargoItemsService{
		createCargoItemFn: func(ctx context.Context, actorID int64, actorRole string, input cargoitems_service.CreateCargoItemInput) (cargoitems_service.CargoItemDTO, error) {
			if actorID != 5 || actorRole != core_domain.RoleWorker.String() {
				t.Fatalf("actor = %d/%q, want 5/worker", actorID, actorRole)
			}
			if input.OrderID != 10 || input.OrderCargoPlaceID != 20 {
				t.Fatalf("input ids = %d/%d, want 10/20", input.OrderID, input.OrderCargoPlaceID)
			}
			if input.QRCode == nil || *input.QRCode != "QR-001" {
				t.Fatalf("QRCode = %v, want QR-001", input.QRCode)
			}
			return cargoitems_service.CargoItemDTO{ID: 100, OrderID: input.OrderID, QRCode: *input.QRCode, Status: "accepted", CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	handler := cargoitems_transport_http.NewCargoItemsHTTPHandler(newTestLogger(t), service)

	requestBody := []byte(`{"order_cargo_place_id":20,"qr_code":"QR-001"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/10/cargo-items", bytes.NewReader(requestBody))
	req.SetPathValue("id", "10")
	req = req.WithContext(core_http_middleware.WithUser(req.Context(), core_http_middleware.CurrentUser{ID: 5, Role: core_domain.RoleWorker.String()}))
	rr := httptest.NewRecorder()

	handler.CreateCargoItem(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body = %s", rr.Code, http.StatusCreated, rr.Body.String())
	}
	var response cargoitems_transport_http.CargoItemResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.CargoItem.ID != 100 || response.CargoItem.QRCode != "QR-001" {
		t.Fatalf("unexpected response: %+v", response)
	}
}

func TestListCargoItemsHTTPParsesQuery(t *testing.T) {
	orderID := int64(10)
	zoneID := int64(9)
	gateID := int64(4)
	service := &fakeCargoItemsService{
		listCargoItemsFn: func(ctx context.Context, actorID int64, actorRole string, filter cargoitems_service.CargoItemFilter) ([]cargoitems_service.CargoItemDTO, error) {
			if actorID != 42 || actorRole != core_domain.RoleClient.String() {
				t.Fatalf("actor = %d/%q, want 42/client", actorID, actorRole)
			}
			if filter.OrderID == nil || *filter.OrderID != orderID {
				t.Fatalf("OrderID = %v, want %d", filter.OrderID, orderID)
			}
			if filter.StorageZoneID == nil || *filter.StorageZoneID != zoneID {
				t.Fatalf("StorageZoneID = %v, want %d", filter.StorageZoneID, zoneID)
			}
			if filter.GateID == nil || *filter.GateID != gateID {
				t.Fatalf("GateID = %v, want %d", filter.GateID, gateID)
			}
			if filter.Status != "stored" || filter.QRCode != "QR-777" || filter.Page != 2 || filter.Limit != 50 {
				t.Fatalf("unexpected filter: %+v", filter)
			}
			return []cargoitems_service.CargoItemDTO{{ID: 7, OrderID: orderID, QRCode: "QR-777", Status: "stored"}}, nil
		},
	}
	handler := cargoitems_transport_http.NewCargoItemsHTTPHandler(newTestLogger(t), service)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cargo-items?order_id=10&storage_zone_id=9&gate_id=4&status=stored&qr_code=QR-777&page=2&limit=50", nil)
	req = req.WithContext(core_http_middleware.WithUser(req.Context(), core_http_middleware.CurrentUser{ID: 42, Role: core_domain.RoleClient.String()}))
	rr := httptest.NewRecorder()

	handler.ListCargoItems(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", rr.Code, http.StatusOK, rr.Body.String())
	}
	var response cargoitems_transport_http.CargoItemsResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(response.CargoItems) != 1 || response.CargoItems[0].ID != 7 {
		t.Fatalf("unexpected response: %+v", response)
	}
}

func TestGetCargoItemHTTP(t *testing.T) {
	service := &fakeCargoItemsService{
		getCargoItemFn: func(ctx context.Context, cargoItemID int64, actorID int64, actorRole string) (cargoitems_service.CargoItemDTO, error) {
			if cargoItemID != 7 || actorID != 42 || actorRole != core_domain.RoleClient.String() {
				t.Fatalf("args = %d/%d/%q, want 7/42/client", cargoItemID, actorID, actorRole)
			}
			return cargoitems_service.CargoItemDTO{ID: cargoItemID, QRCode: "QR-777", Status: "stored"}, nil
		},
	}
	handler := cargoitems_transport_http.NewCargoItemsHTTPHandler(newTestLogger(t), service)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cargo-items/7", nil)
	req.SetPathValue("id", "7")
	req = req.WithContext(core_http_middleware.WithUser(req.Context(), core_http_middleware.CurrentUser{ID: 42, Role: core_domain.RoleClient.String()}))
	rr := httptest.NewRecorder()

	handler.GetCargoItem(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", rr.Code, http.StatusOK, rr.Body.String())
	}
}

func TestGetCargoItemHistoryHTTP(t *testing.T) {
	service := &fakeCargoItemsService{
		getCargoItemHistoryFn: func(ctx context.Context, cargoItemID int64, actorID int64, actorRole string) ([]cargoitems_service.CargoStatusHistoryDTO, error) {
			if cargoItemID != 7 || actorID != 5 || actorRole != core_domain.RoleWorker.String() {
				t.Fatalf("args = %d/%d/%q, want 7/5/worker", cargoItemID, actorID, actorRole)
			}
			return []cargoitems_service.CargoStatusHistoryDTO{{ID: 1, CargoItemID: cargoItemID, NewStatus: "accepted", ChangedBy: actorID, ChangedAt: time.Now().UTC()}}, nil
		},
	}
	handler := cargoitems_transport_http.NewCargoItemsHTTPHandler(newTestLogger(t), service)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cargo-items/7/history", nil)
	req.SetPathValue("id", "7")
	req = req.WithContext(core_http_middleware.WithUser(req.Context(), core_http_middleware.CurrentUser{ID: 5, Role: core_domain.RoleWorker.String()}))
	rr := httptest.NewRecorder()

	handler.GetCargoItemHistory(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", rr.Code, http.StatusOK, rr.Body.String())
	}
	var response cargoitems_transport_http.CargoItemHistoryResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(response.History) != 1 || response.History[0].CargoItemID != 7 {
		t.Fatalf("unexpected response: %+v", response)
	}
}

func TestUpdateCargoItemStatusHTTP(t *testing.T) {
	service := &fakeCargoItemsService{
		updateCargoItemStatusFn: func(ctx context.Context, cargoItemID int64, actorID int64, actorRole string, input cargoitems_service.UpdateCargoItemStatusInput) (cargoitems_service.CargoItemDTO, error) {
			if cargoItemID != 7 || actorID != 5 || actorRole != core_domain.RoleWorker.String() || input.Status != "stored" {
				t.Fatalf("unexpected args: %d/%d/%q/%+v", cargoItemID, actorID, actorRole, input)
			}
			return cargoitems_service.CargoItemDTO{ID: cargoItemID, Status: input.Status}, nil
		},
	}
	handler := cargoitems_transport_http.NewCargoItemsHTTPHandler(newTestLogger(t), service)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/cargo-items/7/status", bytes.NewReader([]byte(`{"status":"stored"}`)))
	req.SetPathValue("id", "7")
	req = req.WithContext(core_http_middleware.WithUser(req.Context(), core_http_middleware.CurrentUser{ID: 5, Role: core_domain.RoleWorker.String()}))
	rr := httptest.NewRecorder()

	handler.UpdateCargoItemStatus(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", rr.Code, http.StatusOK, rr.Body.String())
	}
}

func TestUpdateCargoItemStatusHTTPRejectsInvalidBody(t *testing.T) {
	handler := cargoitems_transport_http.NewCargoItemsHTTPHandler(newTestLogger(t), &fakeCargoItemsService{})

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/cargo-items/7/status", bytes.NewReader([]byte(`{}`)))
	req.SetPathValue("id", "7")
	req = req.WithContext(core_http_middleware.WithUser(req.Context(), core_http_middleware.CurrentUser{ID: 5, Role: core_domain.RoleWorker.String()}))
	rr := httptest.NewRecorder()

	handler.UpdateCargoItemStatus(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body = %s", rr.Code, http.StatusBadRequest, rr.Body.String())
	}
}

func TestAssignStorageZoneHTTP(t *testing.T) {
	service := &fakeCargoItemsService{
		assignStorageZoneFn: func(ctx context.Context, cargoItemID int64, actorID int64, actorRole string, input cargoitems_service.AssignStorageZoneInput) (cargoitems_service.CargoItemDTO, error) {
			if cargoItemID != 7 || actorID != 6 || actorRole != core_domain.RoleLogist.String() || input.StorageZoneID != 9 {
				t.Fatalf("unexpected args: %d/%d/%q/%+v", cargoItemID, actorID, actorRole, input)
			}
			return cargoitems_service.CargoItemDTO{ID: cargoItemID, Status: "accepted", StorageZoneID: &input.StorageZoneID}, nil
		},
	}
	handler := cargoitems_transport_http.NewCargoItemsHTTPHandler(newTestLogger(t), service)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/cargo-items/7/assign-zone", bytes.NewReader([]byte(`{"storage_zone_id":9}`)))
	req.SetPathValue("id", "7")
	req = req.WithContext(core_http_middleware.WithUser(req.Context(), core_http_middleware.CurrentUser{ID: 6, Role: core_domain.RoleLogist.String()}))
	rr := httptest.NewRecorder()

	handler.AssignStorageZone(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", rr.Code, http.StatusOK, rr.Body.String())
	}
}

func TestAssignGateHTTP(t *testing.T) {
	service := &fakeCargoItemsService{
		assignGateFn: func(ctx context.Context, cargoItemID int64, actorID int64, actorRole string, input cargoitems_service.AssignGateInput) (cargoitems_service.CargoItemDTO, error) {
			if cargoItemID != 7 || actorID != 6 || actorRole != core_domain.RoleLogist.String() || input.GateID != 4 {
				t.Fatalf("unexpected args: %d/%d/%q/%+v", cargoItemID, actorID, actorRole, input)
			}
			return cargoitems_service.CargoItemDTO{ID: cargoItemID, Status: "stored", GateID: &input.GateID}, nil
		},
	}
	handler := cargoitems_transport_http.NewCargoItemsHTTPHandler(newTestLogger(t), service)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/cargo-items/7/assign-gate", bytes.NewReader([]byte(`{"gate_id":4}`)))
	req.SetPathValue("id", "7")
	req = req.WithContext(core_http_middleware.WithUser(req.Context(), core_http_middleware.CurrentUser{ID: 6, Role: core_domain.RoleLogist.String()}))
	rr := httptest.NewRecorder()

	handler.AssignGate(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", rr.Code, http.StatusOK, rr.Body.String())
	}
}

func newTestLogger(t *testing.T) *core_logger.Logger {
	t.Helper()
	log, err := core_logger.NewLogger(core_logger.LoggerConfig{
		Level:  "debug",
		Folder: t.TempDir(),
	})
	if err != nil {
		t.Fatalf("new logger: %v", err)
	}
	t.Cleanup(log.Close)
	return log
}

func TestScanCargoItemHTTP(t *testing.T) {
	service := &fakeCargoItemsService{
		scanCargoItemFn: func(ctx context.Context, actorID int64, actorRole string, qrCode string) (cargoitems_service.CargoItemDTO, error) {
			if actorID != 5 || actorRole != core_domain.RoleWorker.String() || qrCode != "QR-SCAN-001" {
				t.Fatalf("unexpected args: %d/%q/%q", actorID, actorRole, qrCode)
			}
			return cargoitems_service.CargoItemDTO{ID: 7, QRCode: qrCode, Status: "stored"}, nil
		},
	}
	handler := cargoitems_transport_http.NewCargoItemsHTTPHandler(newTestLogger(t), service)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cargo-items/scan?qr_code=QR-SCAN-001", nil)
	req = req.WithContext(core_http_middleware.WithUser(req.Context(), core_http_middleware.CurrentUser{ID: 5, Role: core_domain.RoleWorker.String()}))
	rr := httptest.NewRecorder()

	handler.ScanCargoItem(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", rr.Code, http.StatusOK, rr.Body.String())
	}
	var response cargoitems_transport_http.CargoItemResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.CargoItem.ID != 7 || response.CargoItem.QRCode != "QR-SCAN-001" {
		t.Fatalf("unexpected response: %+v", response)
	}
}

func TestGetCargoItemLabelHTTP(t *testing.T) {
	service := &fakeCargoItemsService{
		getCargoItemLabelFn: func(ctx context.Context, cargoItemID int64, actorID int64, actorRole string) (cargoitems_service.CargoItemLabelDTO, error) {
			if cargoItemID != 7 || actorID != 42 || actorRole != core_domain.RoleClient.String() {
				t.Fatalf("unexpected args: %d/%d/%q", cargoItemID, actorID, actorRole)
			}
			return cargoitems_service.CargoItemLabelDTO{CargoItemID: cargoItemID, OrderID: 10, QRCode: "QR-LABEL-001", QRCodeValue: "QR-LABEL-001", Status: "stored", LabelText: "ORDER-10 | CARGO-7 | QR-LABEL-001"}, nil
		},
	}
	handler := cargoitems_transport_http.NewCargoItemsHTTPHandler(newTestLogger(t), service)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cargo-items/7/label", nil)
	req.SetPathValue("id", "7")
	req = req.WithContext(core_http_middleware.WithUser(req.Context(), core_http_middleware.CurrentUser{ID: 42, Role: core_domain.RoleClient.String()}))
	rr := httptest.NewRecorder()

	handler.GetCargoItemLabel(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", rr.Code, http.StatusOK, rr.Body.String())
	}
	var response cargoitems_transport_http.CargoItemLabelResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Label.CargoItemID != 7 || response.Label.QRCodeValue != "QR-LABEL-001" {
		t.Fatalf("unexpected response: %+v", response)
	}
}
