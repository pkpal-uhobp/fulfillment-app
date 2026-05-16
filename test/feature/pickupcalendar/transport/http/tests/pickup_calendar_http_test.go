package pickupcalendar_transport_http_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_http_middleware "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/middleware"
	pickupcalendar_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/pickupcalendar/service"
)

func decodeJSON(t *testing.T, recorder *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var body map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v; body=%s", err, recorder.Body.String())
	}
	return body
}

func withUser(req *http.Request, id int64, role string) *http.Request {
	ctx := core_http_middleware.WithUser(req.Context(), core_http_middleware.CurrentUser{
		ID:   id,
		Role: role,
		JTI:  "test-jti",
	})
	return req.WithContext(ctx)
}

func TestGetCalendarHTTP(t *testing.T) {
	service := &fakePickupCalendarService{
		getCalendarFn: func(_ context.Context, actorID int64, actorRole string, filter pickupcalendar_service.CalendarFilter) ([]pickupcalendar_service.PickupCalendarDayDTO, error) {
			if actorID != 5 || actorRole != core_domain.RoleClient.String() {
				t.Fatalf("actor = %d/%s, want 5/client", actorID, actorRole)
			}
			if filter.WarehouseID != 10 || filter.DateFrom != "2026-06-01" || filter.DateTo != "2026-06-02" {
				t.Fatalf("unexpected filter: %+v", filter)
			}
			return []pickupcalendar_service.PickupCalendarDayDTO{
				{WarehouseID: 10, Date: "2026-06-01", MaxOrders: 5, CurrentOrders: 0, IsClosed: false},
			}, nil
		},
	}
	handler := newTestHandler(t, service)
	req := httptest.NewRequest(http.MethodGet, "/pickup-calendar?warehouse_id=10&date_from=2026-06-01&date_to=2026-06-02", nil)
	req = withUser(req, 5, core_domain.RoleClient.String())
	rec := httptest.NewRecorder()

	handler.GetCalendar(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	body := decodeJSON(t, rec)
	days := body["days"].([]any)
	if len(days) != 1 {
		t.Fatalf("days len = %d, want 1", len(days))
	}
}

func TestBlockDateHTTP(t *testing.T) {
	service := &fakePickupCalendarService{
		blockDateFn: func(_ context.Context, actorID int64, actorRole string, input pickupcalendar_service.BlockDateInput) (pickupcalendar_service.PickupCalendarBlockDTO, error) {
			if actorID != 8 || actorRole != core_domain.RoleLogist.String() {
				t.Fatalf("actor = %d/%s, want 8/logist", actorID, actorRole)
			}
			if input.WarehouseID != 10 || input.BlockedDate != "2026-06-10" {
				t.Fatalf("unexpected input: %+v", input)
			}
			return pickupcalendar_service.PickupCalendarBlockDTO{ID: 1, WarehouseID: input.WarehouseID, BlockedDate: input.BlockedDate, CreatedBy: actorID}, nil
		},
	}
	handler := newTestHandler(t, service)
	body := []byte(`{"warehouse_id":10,"blocked_date":"2026-06-10","reason":"Перегрузка"}`)
	req := httptest.NewRequest(http.MethodPost, "/pickup-calendar/blocks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withUser(req, 8, core_domain.RoleLogist.String())
	rec := httptest.NewRecorder()

	handler.BlockDate(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusCreated, rec.Body.String())
	}
}

func TestBlockDateHTTPRejectsInvalidBody(t *testing.T) {
	handler := newTestHandler(t, &fakePickupCalendarService{})
	req := httptest.NewRequest(http.MethodPost, "/pickup-calendar/blocks", bytes.NewReader([]byte(`{"warehouse_id":0}`)))
	req.Header.Set("Content-Type", "application/json")
	req = withUser(req, 8, core_domain.RoleLogist.String())
	rec := httptest.NewRecorder()

	handler.BlockDate(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}

func TestUnblockDateHTTP(t *testing.T) {
	called := false
	service := &fakePickupCalendarService{
		unblockDateFn: func(_ context.Context, actorID int64, actorRole string, blockID int64) error {
			if actorID != 8 || actorRole != core_domain.RoleAdmin.String() || blockID != 33 {
				t.Fatalf("unexpected args: actor=%d/%s blockID=%d", actorID, actorRole, blockID)
			}
			called = true
			return nil
		},
	}
	handler := newTestHandler(t, service)
	req := httptest.NewRequest(http.MethodDelete, "/pickup-calendar/blocks/33", nil)
	req.SetPathValue("id", "33")
	req = withUser(req, 8, core_domain.RoleAdmin.String())
	rec := httptest.NewRecorder()

	handler.UnblockDate(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusNoContent, rec.Body.String())
	}
	if !called {
		t.Fatal("UnblockDate service was not called")
	}
}

func TestSetCapacityHTTP(t *testing.T) {
	service := &fakePickupCalendarService{
		setCapacityFn: func(_ context.Context, actorID int64, actorRole string, input pickupcalendar_service.SetCapacityInput) (pickupcalendar_service.PickupCapacityDTO, error) {
			if actorID != 8 || actorRole != core_domain.RoleLogist.String() {
				t.Fatalf("actor = %d/%s, want 8/logist", actorID, actorRole)
			}
			if input.WarehouseID != 10 || input.PickupDate != "2026-06-11" || input.MaxOrders != 7 || input.CurrentOrders != 2 || input.IsClosed {
				t.Fatalf("unexpected input: %+v", input)
			}
			return pickupcalendar_service.PickupCapacityDTO{ID: 4, WarehouseID: input.WarehouseID, PickupDate: input.PickupDate, MaxOrders: input.MaxOrders, CurrentOrders: input.CurrentOrders}, nil
		},
	}
	handler := newTestHandler(t, service)
	body := []byte(`{"warehouse_id":10,"pickup_date":"2026-06-11","max_orders":7,"current_orders":2,"is_closed":false}`)
	req := httptest.NewRequest(http.MethodPatch, "/pickup-calendar/capacity", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withUser(req, 8, core_domain.RoleLogist.String())
	rec := httptest.NewRecorder()

	handler.SetCapacity(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
}

func TestGetCalendarHTTPRejectsInvalidWarehouseID(t *testing.T) {
	handler := newTestHandler(t, &fakePickupCalendarService{})
	req := httptest.NewRequest(http.MethodGet, "/pickup-calendar?warehouse_id=bad", nil)
	req = withUser(req, 5, core_domain.RoleClient.String())
	rec := httptest.NewRecorder()

	handler.GetCalendar(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}
