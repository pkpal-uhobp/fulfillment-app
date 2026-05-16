package shipments_transport_http_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_http_middleware "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/middleware"
	shipments_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/shipments/service"
)

func withUser(req *http.Request, id int64, role string) *http.Request {
	ctx := core_http_middleware.WithUser(req.Context(), core_http_middleware.CurrentUser{
		ID:   id,
		Role: role,
		JTI:  "test-jti",
	})
	return req.WithContext(ctx)
}

func decodeJSON(t *testing.T, recorder *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var body map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v; body=%s", err, recorder.Body.String())
	}
	return body
}

func TestRoutesExposeShipmentEndpoints(t *testing.T) {
	handler := newTestHandler(t, &fakeShipmentsService{})
	routes := handler.Routes()
	if len(routes) != 6 {
		t.Fatalf("routes len = %d, want 6", len(routes))
	}
}

func TestCreateShipmentHTTP(t *testing.T) {
	service := &fakeShipmentsService{
		createShipmentFn: func(_ context.Context, actorID int64, actorRole string, input shipments_service.CreateShipmentInput) (shipments_service.ShipmentDTO, error) {
			if actorID != 7 || actorRole != core_domain.RoleLogist.String() {
				t.Fatalf("actor = %d/%s, want 7/logist", actorID, actorRole)
			}
			if input.DestinationWarehouseID != 10 || input.GateID != 20 || input.PlannedDepartureAt != "2026-06-20T12:00:00Z" {
				t.Fatalf("unexpected input: %+v", input)
			}
			return shipments_service.ShipmentDTO{ID: 100, DestinationWarehouseID: 10, GateID: 20, Status: core_domain.ShipmentStatusPlanned.String()}, nil
		},
	}
	handler := newTestHandler(t, service)
	body := []byte(`{"destination_warehouse_id":10,"gate_id":20,"planned_departure_at":"2026-06-20T12:00:00Z"}`)
	req := httptest.NewRequest(http.MethodPost, "/shipments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withUser(req, 7, core_domain.RoleLogist.String())
	rec := httptest.NewRecorder()

	handler.CreateShipment(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusCreated, rec.Body.String())
	}
	_ = decodeJSON(t, rec)
}

func TestListShipmentsHTTPParsesQuery(t *testing.T) {
	service := &fakeShipmentsService{
		listShipmentsFn: func(_ context.Context, actorID int64, actorRole string, filter shipments_service.ShipmentFilter) ([]shipments_service.ShipmentDTO, error) {
			if actorID != 8 || actorRole != core_domain.RoleAdmin.String() {
				t.Fatalf("actor = %d/%s, want 8/admin", actorID, actorRole)
			}
			if filter.Status != "planned" || filter.DestinationWarehouseID == nil || *filter.DestinationWarehouseID != 10 || filter.GateID == nil || *filter.GateID != 20 || filter.Date != "2026-06-20" {
				t.Fatalf("unexpected filter: %+v", filter)
			}
			return []shipments_service.ShipmentDTO{{ID: 100, Status: "planned"}}, nil
		},
	}
	handler := newTestHandler(t, service)
	req := httptest.NewRequest(http.MethodGet, "/shipments?status=planned&destination_warehouse_id=10&gate_id=20&date=2026-06-20", nil)
	req = withUser(req, 8, core_domain.RoleAdmin.String())
	rec := httptest.NewRecorder()

	handler.ListShipments(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
}

func TestAddShipmentItemHTTP(t *testing.T) {
	service := &fakeShipmentsService{
		addShipmentItemFn: func(_ context.Context, shipmentID int64, actorID int64, actorRole string, input shipments_service.AddShipmentItemInput) (shipments_service.ShipmentDTO, error) {
			if shipmentID != 100 || actorID != 7 || actorRole != core_domain.RoleLogist.String() || input.CargoItemID != 55 {
				t.Fatalf("unexpected args: shipment=%d actor=%d/%s input=%+v", shipmentID, actorID, actorRole, input)
			}
			return shipments_service.ShipmentDTO{ID: shipmentID, Items: []shipments_service.ShipmentItemDTO{{CargoItemID: input.CargoItemID}}}, nil
		},
	}
	handler := newTestHandler(t, service)
	req := httptest.NewRequest(http.MethodPost, "/shipments/100/items", bytes.NewReader([]byte(`{"cargo_item_id":55}`)))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "100")
	req = withUser(req, 7, core_domain.RoleLogist.String())
	rec := httptest.NewRecorder()

	handler.AddShipmentItem(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
}

func TestUpdateShipmentStatusHTTP(t *testing.T) {
	service := &fakeShipmentsService{
		updateShipmentStatusFn: func(_ context.Context, shipmentID int64, actorID int64, actorRole string, input shipments_service.UpdateShipmentStatusInput) (shipments_service.ShipmentDTO, error) {
			if shipmentID != 100 || actorID != 7 || input.Status != core_domain.ShipmentStatusLoading.String() {
				t.Fatalf("unexpected args: shipment=%d actor=%d input=%+v", shipmentID, actorID, input)
			}
			return shipments_service.ShipmentDTO{ID: shipmentID, Status: input.Status}, nil
		},
	}
	handler := newTestHandler(t, service)
	req := httptest.NewRequest(http.MethodPatch, "/shipments/100/status", bytes.NewReader([]byte(`{"status":"loading"}`)))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "100")
	req = withUser(req, 7, core_domain.RoleLogist.String())
	rec := httptest.NewRecorder()

	handler.UpdateShipmentStatus(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
}

func TestRemoveShipmentItemHTTP(t *testing.T) {
	called := false
	service := &fakeShipmentsService{
		removeShipmentItemFn: func(_ context.Context, shipmentID int64, cargoItemID int64, actorID int64, actorRole string) error {
			if shipmentID != 100 || cargoItemID != 55 || actorID != 7 {
				t.Fatalf("unexpected args: shipment=%d cargo=%d actor=%d", shipmentID, cargoItemID, actorID)
			}
			called = true
			return nil
		},
	}
	handler := newTestHandler(t, service)
	req := httptest.NewRequest(http.MethodDelete, "/shipments/100/items/55", nil)
	req.SetPathValue("id", "100")
	req.SetPathValue("cargo_item_id", "55")
	req = withUser(req, 7, core_domain.RoleAdmin.String())
	rec := httptest.NewRecorder()

	handler.RemoveShipmentItem(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusNoContent, rec.Body.String())
	}
	if !called {
		t.Fatal("RemoveShipmentItem service was not called")
	}
}
