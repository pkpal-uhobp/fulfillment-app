package users_transport_http_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_logger "github.com/pkpal-uhobp/fulfillment-app/internal/core/logger"
	core_http_middleware "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/middleware"
	users_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/users/service"
	users_transport_http "github.com/pkpal-uhobp/fulfillment-app/internal/features/users/transport/http"
)

func newTestLogger(t *testing.T) *core_logger.Logger {
	t.Helper()
	log, err := core_logger.NewLogger(core_logger.LoggerConfig{Level: "debug", Folder: t.TempDir()})
	if err != nil {
		t.Fatalf("create logger: %v", err)
	}
	t.Cleanup(log.Close)
	return log
}

func newTestHandler(t *testing.T, service *fakeUsersService) *users_transport_http.UsersHTTPHandler {
	t.Helper()
	return users_transport_http.NewUsersHTTPHandler(newTestLogger(t), service)
}

func withAdmin(req *http.Request) *http.Request {
	ctx := core_http_middleware.WithUser(req.Context(), core_http_middleware.CurrentUser{ID: 1, Role: core_domain.RoleAdmin.String(), JTI: "test"})
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

func TestListUsersReturnsOKAndPassesFilters(t *testing.T) {
	service := &fakeUsersService{
		listUsersFn: func(_ context.Context, actorRole string, filter users_service.UserFilter) ([]users_service.UserDTO, error) {
			if actorRole != core_domain.RoleAdmin.String() {
				t.Fatalf("actorRole = %q, want admin", actorRole)
			}
			if filter.Role != "worker" || filter.Search != "ivan" {
				t.Fatalf("unexpected filter: %+v", filter)
			}
			if filter.IsActive == nil || *filter.IsActive != true {
				t.Fatalf("IsActive = %v, want true", filter.IsActive)
			}
			return []users_service.UserDTO{{ID: 2, Email: "worker@test.com", FullName: "Worker", Role: "worker", IsActive: true}}, nil
		},
	}
	handler := newTestHandler(t, service)
	req := withAdmin(httptest.NewRequest(http.MethodGet, "/users?role=worker&search=ivan&is_active=true", nil))
	rec := httptest.NewRecorder()
	handler.ListUsers(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	body := decodeJSON(t, rec)
	if body["users"] == nil {
		t.Fatalf("response must contain users")
	}
}

func TestCreateUserReturnsCreated(t *testing.T) {
	service := &fakeUsersService{
		createUserFn: func(_ context.Context, actorRole string, input users_service.CreateUserInput) (users_service.UserDTO, error) {
			if input.Email != "worker@test.com" || input.Role != "worker" {
				t.Fatalf("unexpected input: %+v", input)
			}
			return users_service.UserDTO{ID: 2, Email: input.Email, FullName: input.FullName, Role: input.Role, IsActive: true}, nil
		},
	}
	handler := newTestHandler(t, service)
	body := bytes.NewReader([]byte(`{"email":"worker@test.com","password":"123456","full_name":"Worker Test","role":"worker"}`))
	req := withAdmin(httptest.NewRequest(http.MethodPost, "/users", body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.CreateUser(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusCreated, rec.Body.String())
	}
}

func TestPatchUserReturnsOK(t *testing.T) {
	service := &fakeUsersService{
		patchUserFn: func(_ context.Context, actorID int64, actorRole string, userID int64, input users_service.PatchUserInput) (users_service.UserDTO, error) {
			if actorID != 1 || userID != 2 {
				t.Fatalf("actorID=%d userID=%d", actorID, userID)
			}
			if input.FullName == nil || *input.FullName != "Updated User" {
				t.Fatalf("unexpected input: %+v", input)
			}
			return users_service.UserDTO{ID: userID, Email: "u@test.com", FullName: *input.FullName, Role: "worker", IsActive: true}, nil
		},
	}
	handler := newTestHandler(t, service)
	req := withAdmin(httptest.NewRequest(http.MethodPatch, "/users/2", bytes.NewReader([]byte(`{"full_name":"Updated User"}`))))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "2")
	rec := httptest.NewRecorder()
	handler.PatchUser(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
}

func TestBlockAndDeactivateReturnNoContent(t *testing.T) {
	calls := map[string]int{}
	service := &fakeUsersService{
		blockUserFn: func(_ context.Context, actorID int64, actorRole string, userID int64, input users_service.BlockUserInput) error {
			if userID != 2 {
				t.Fatalf("userID = %d, want 2", userID)
			}
			calls["block"]++
			return nil
		},
		deactivateUserFn: func(_ context.Context, actorID int64, actorRole string, userID int64) error {
			if userID != 2 {
				t.Fatalf("userID = %d, want 2", userID)
			}
			calls["delete"]++
			return nil
		},
	}
	handler := newTestHandler(t, service)

	blockReq := withAdmin(httptest.NewRequest(http.MethodPatch, "/users/2/block", bytes.NewReader([]byte(`{"reason":"test"}`))))
	blockReq.Header.Set("Content-Type", "application/json")
	blockReq.SetPathValue("id", "2")
	blockRec := httptest.NewRecorder()
	handler.BlockUser(blockRec, blockReq)
	if blockRec.Code != http.StatusNoContent {
		t.Fatalf("block status = %d, want %d", blockRec.Code, http.StatusNoContent)
	}

	deleteReq := withAdmin(httptest.NewRequest(http.MethodDelete, "/users/2", nil))
	deleteReq.SetPathValue("id", "2")
	deleteRec := httptest.NewRecorder()
	handler.DeactivateUser(deleteRec, deleteReq)
	if deleteRec.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want %d", deleteRec.Code, http.StatusNoContent)
	}

	if calls["block"] != 1 || calls["delete"] != 1 {
		t.Fatalf("unexpected calls: %+v", calls)
	}
}
