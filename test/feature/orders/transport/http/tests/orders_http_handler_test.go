package orders_http_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
	core_logger "github.com/pkpal-uhobp/fulfillment-app/internal/core/logger"
	core_http_middleware "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/middleware"
	orders_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/orders/service"
	orders_transport_http "github.com/pkpal-uhobp/fulfillment-app/internal/features/orders/transport/http"
)

func newTestLogger(t *testing.T) *core_logger.Logger {
	t.Helper()

	log, err := core_logger.NewLogger(core_logger.LoggerConfig{
		Level:  "debug",
		Folder: t.TempDir(),
	})
	if err != nil {
		t.Fatalf("create logger: %v", err)
	}
	t.Cleanup(log.Close)
	return log
}

func newTestHandler(t *testing.T, service *fakeOrdersService) *orders_transport_http.OrdersHTTPHandler {
	t.Helper()
	return orders_transport_http.NewOrdersHTTPHandler(newTestLogger(t), service)
}

func decodeJSONBody(t *testing.T, recorder *httptest.ResponseRecorder) map[string]any {
	t.Helper()

	var body map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response body: %v; body=%s", err, recorder.Body.String())
	}
	return body
}

func withCurrentUser(r *http.Request, userID int64, role string) *http.Request {
	ctx := core_http_middleware.WithUser(
		r.Context(),
		core_http_middleware.CurrentUser{
			ID:   userID,
			Role: role,
			JTI:  "test-jti",
		},
	)
	return r.WithContext(ctx)
}

func testOrderDTO(id int64, clientID int64, status string) orders_service.OrderDTO {
	now := time.Now().UTC()
	return orders_service.OrderDTO{
		ID:                     id,
		ClientID:               clientID,
		ReceivingWarehouseID:   10,
		DestinationWarehouseID: 20,
		ProductTypeID:          30,
		HandoverType:           core_domain.HandoverTypeSelfDelivery.String(),
		Status:                 status,
		CargoPlaces: []orders_service.CargoPlaceDTO{
			{ID: 1, OrderID: id, CargoPlaceTypeID: 7, Quantity: 2, CreatedAt: now},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func TestCreateOrderReturnsCreatedAndPassesInputToService(t *testing.T) {
	service := &fakeOrdersService{
		createOrderFn: func(ctx context.Context, actorID int64, actorRole string, input orders_service.CreateOrderInput) (orders_service.OrderDTO, error) {
			if actorID != 42 {
				t.Fatalf("actorID = %d, want 42", actorID)
			}
			if actorRole != core_domain.RoleClient.String() {
				t.Fatalf("actorRole = %q, want client", actorRole)
			}
			if input.ReceivingWarehouseID != 10 || input.DestinationWarehouseID != 20 || input.ProductTypeID != 30 {
				t.Fatalf("unexpected warehouses/product type: %+v", input)
			}
			if input.HandoverType != core_domain.HandoverTypeSelfDelivery.String() {
				t.Fatalf("HandoverType = %q, want self_delivery", input.HandoverType)
			}
			if input.SelfDeliveryDate == nil || *input.SelfDeliveryDate != "2026-05-20" {
				t.Fatalf("SelfDeliveryDate = %v, want 2026-05-20", input.SelfDeliveryDate)
			}
			if len(input.CargoPlaces) != 1 || input.CargoPlaces[0].CargoPlaceTypeID != 7 || input.CargoPlaces[0].Quantity != 2 {
				t.Fatalf("unexpected cargo places: %+v", input.CargoPlaces)
			}
			return testOrderDTO(100, actorID, core_domain.OrderStatusCreated.String()), nil
		},
	}

	handler := newTestHandler(t, service)
	request := httptest.NewRequest(
		http.MethodPost,
		"/orders",
		strings.NewReader(`{
			"receiving_warehouse_id":10,
			"destination_warehouse_id":20,
			"product_type_id":30,
			"handover_type":"self_delivery",
			"self_delivery_date":"2026-05-20",
			"cargo_places":[{"cargo_place_type_id":7,"quantity":2}]
		}`),
	)
	request.Header.Set("Content-Type", "application/json")
	request = withCurrentUser(request, 42, core_domain.RoleClient.String())

	recorder := httptest.NewRecorder()
	handler.CreateOrder(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusCreated, recorder.Body.String())
	}
	body := decodeJSONBody(t, recorder)
	if body["order"] == nil {
		t.Fatalf("response must contain order")
	}
}

func TestCreateOrderReturnsBadRequestForInvalidBody(t *testing.T) {
	handler := newTestHandler(t, &fakeOrdersService{})
	request := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(`{"handover_type":"self_delivery"}`))
	request.Header.Set("Content-Type", "application/json")
	request = withCurrentUser(request, 42, core_domain.RoleClient.String())

	recorder := httptest.NewRecorder()
	handler.CreateOrder(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusBadRequest, recorder.Body.String())
	}
}

func TestListOrdersReturnsOKAndPassesFilterToService(t *testing.T) {
	service := &fakeOrdersService{
		listOrdersFn: func(ctx context.Context, actorID int64, actorRole string, filter orders_service.OrderFilter) ([]orders_service.OrderDTO, error) {
			if actorID != 1 || actorRole != core_domain.RoleAdmin.String() {
				t.Fatalf("unexpected actor: %d/%s", actorID, actorRole)
			}
			if filter.ClientID == nil || *filter.ClientID != 42 {
				t.Fatalf("ClientID = %v, want 42", filter.ClientID)
			}
			if filter.Status != core_domain.OrderStatusCreated.String() {
				t.Fatalf("Status = %q, want created", filter.Status)
			}
			if filter.HandoverType != core_domain.HandoverTypeSelfDelivery.String() {
				t.Fatalf("HandoverType = %q, want self_delivery", filter.HandoverType)
			}
			if filter.WarehouseID == nil || *filter.WarehouseID != 10 {
				t.Fatalf("WarehouseID = %v, want 10", filter.WarehouseID)
			}
			if filter.Page != 2 || filter.Limit != 5 {
				t.Fatalf("Page/Limit = %d/%d, want 2/5", filter.Page, filter.Limit)
			}
			return []orders_service.OrderDTO{testOrderDTO(100, 42, core_domain.OrderStatusCreated.String())}, nil
		},
	}

	handler := newTestHandler(t, service)
	request := httptest.NewRequest(http.MethodGet, "/orders?client_id=42&status=created&handover_type=self_delivery&warehouse_id=10&page=2&limit=5", nil)
	request = withCurrentUser(request, 1, core_domain.RoleAdmin.String())

	recorder := httptest.NewRecorder()
	handler.ListOrders(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	body := decodeJSONBody(t, recorder)
	if body["orders"] == nil {
		t.Fatalf("response must contain orders")
	}
}

func TestListOrdersReturnsBadRequestForInvalidQuery(t *testing.T) {
	handler := newTestHandler(t, &fakeOrdersService{})
	request := httptest.NewRequest(http.MethodGet, "/orders?page=bad", nil)
	request = withCurrentUser(request, 1, core_domain.RoleAdmin.String())

	recorder := httptest.NewRecorder()
	handler.ListOrders(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusBadRequest, recorder.Body.String())
	}
}

func TestGetOrderReturnsOK(t *testing.T) {
	service := &fakeOrdersService{
		getOrderFn: func(ctx context.Context, orderID int64, actorID int64, actorRole string) (orders_service.OrderDTO, error) {
			if orderID != 100 || actorID != 42 || actorRole != core_domain.RoleClient.String() {
				t.Fatalf("unexpected args: orderID=%d actor=%d/%s", orderID, actorID, actorRole)
			}
			return testOrderDTO(orderID, actorID, core_domain.OrderStatusCreated.String()), nil
		},
	}

	handler := newTestHandler(t, service)
	request := httptest.NewRequest(http.MethodGet, "/orders/100", nil)
	request.SetPathValue("id", "100")
	request = withCurrentUser(request, 42, core_domain.RoleClient.String())

	recorder := httptest.NewRecorder()
	handler.GetOrder(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
}

func TestGetOrderReturnsBadRequestForInvalidID(t *testing.T) {
	handler := newTestHandler(t, &fakeOrdersService{})
	request := httptest.NewRequest(http.MethodGet, "/orders/bad", nil)
	request.SetPathValue("id", "bad")
	request = withCurrentUser(request, 42, core_domain.RoleClient.String())

	recorder := httptest.NewRecorder()
	handler.GetOrder(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusBadRequest, recorder.Body.String())
	}
}

func TestGetOrderHistoryReturnsOK(t *testing.T) {
	service := &fakeOrdersService{
		getOrderHistoryFn: func(ctx context.Context, orderID int64, actorID int64, actorRole string) ([]orders_service.OrderStatusHistoryDTO, error) {
			if orderID != 100 || actorID != 42 || actorRole != core_domain.RoleClient.String() {
				t.Fatalf("unexpected args: orderID=%d actor=%d/%s", orderID, actorID, actorRole)
			}
			oldStatus := core_domain.OrderStatusCreated.String()
			return []orders_service.OrderStatusHistoryDTO{
				{ID: 1, OrderID: orderID, OldStatus: &oldStatus, NewStatus: core_domain.OrderStatusWaitingDelivery.String(), ChangedBy: 1, ChangedAt: time.Now().UTC()},
			}, nil
		},
	}

	handler := newTestHandler(t, service)
	request := httptest.NewRequest(http.MethodGet, "/orders/100/history", nil)
	request.SetPathValue("id", "100")
	request = withCurrentUser(request, 42, core_domain.RoleClient.String())

	recorder := httptest.NewRecorder()
	handler.GetOrderHistory(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	body := decodeJSONBody(t, recorder)
	if body["history"] == nil {
		t.Fatalf("response must contain history")
	}
}

func TestCancelOrderReturnsNoContentAndPassesComment(t *testing.T) {
	comment := "Клиент отменил заявку"
	service := &fakeOrdersService{
		cancelOrderFn: func(ctx context.Context, orderID int64, actorID int64, actorRole string, input orders_service.CancelOrderInput) error {
			if orderID != 100 || actorID != 42 || actorRole != core_domain.RoleClient.String() {
				t.Fatalf("unexpected args: orderID=%d actor=%d/%s", orderID, actorID, actorRole)
			}
			if input.Comment == nil || *input.Comment != comment {
				t.Fatalf("Comment = %v, want %s", input.Comment, comment)
			}
			return nil
		},
	}

	handler := newTestHandler(t, service)
	request := httptest.NewRequest(http.MethodPatch, "/orders/100/cancel", strings.NewReader(fmt.Sprintf(`{"comment":%q}`, comment)))
	request.Header.Set("Content-Type", "application/json")
	request.SetPathValue("id", "100")
	request = withCurrentUser(request, 42, core_domain.RoleClient.String())

	recorder := httptest.NewRecorder()
	handler.CancelOrder(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusNoContent, recorder.Body.String())
	}
}

func TestUpdateOrderStatusReturnsOKAndPassesInput(t *testing.T) {
	comment := "Ожидает сдачи на склад"
	service := &fakeOrdersService{
		updateOrderStatusFn: func(ctx context.Context, orderID int64, actorID int64, input orders_service.UpdateOrderStatusInput) (orders_service.OrderDTO, error) {
			if orderID != 100 || actorID != 5 {
				t.Fatalf("unexpected args: orderID=%d actorID=%d", orderID, actorID)
			}
			if input.Status != core_domain.OrderStatusWaitingDelivery.String() {
				t.Fatalf("Status = %q, want waiting_delivery", input.Status)
			}
			if input.Comment == nil || *input.Comment != comment {
				t.Fatalf("Comment = %v, want %s", input.Comment, comment)
			}
			return testOrderDTO(orderID, 42, input.Status), nil
		},
	}

	handler := newTestHandler(t, service)
	request := httptest.NewRequest(http.MethodPatch, "/orders/100/status", strings.NewReader(fmt.Sprintf(`{"status":"waiting_delivery","comment":%q}`, comment)))
	request.Header.Set("Content-Type", "application/json")
	request.SetPathValue("id", "100")
	request = withCurrentUser(request, 5, core_domain.RoleLogist.String())

	recorder := httptest.NewRecorder()
	handler.UpdateOrderStatus(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	body := decodeJSONBody(t, recorder)
	if body["order"] == nil {
		t.Fatalf("response must contain order")
	}
}

func TestUpdateOrderStatusReturnsBadRequestFromService(t *testing.T) {
	service := &fakeOrdersService{
		updateOrderStatusFn: func(ctx context.Context, orderID int64, actorID int64, input orders_service.UpdateOrderStatusInput) (orders_service.OrderDTO, error) {
			return orders_service.OrderDTO{}, fmt.Errorf("%w: invalid transition", core_errors.ErrInvalidArgument)
		},
	}

	handler := newTestHandler(t, service)
	request := httptest.NewRequest(http.MethodPatch, "/orders/100/status", strings.NewReader(`{"status":"received"}`))
	request.Header.Set("Content-Type", "application/json")
	request.SetPathValue("id", "100")
	request = withCurrentUser(request, 5, core_domain.RoleLogist.String())

	recorder := httptest.NewRecorder()
	handler.UpdateOrderStatus(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusBadRequest, recorder.Body.String())
	}
}
