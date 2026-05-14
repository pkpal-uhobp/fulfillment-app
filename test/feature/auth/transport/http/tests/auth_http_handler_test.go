package auth_http_tests

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
	auth_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/auth/service"
	auth_transport_http "github.com/pkpal-uhobp/fulfillment-app/internal/features/auth/transport/http"
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

func newTestHandler(t *testing.T, service *fakeAuthService) *auth_transport_http.AuthHTTPHandler {
	t.Helper()

	return auth_transport_http.NewAuthHTTPHandler(
		newTestLogger(t),
		service,
		nil,
	)
}

func decodeJSONBody(t *testing.T, recorder *httptest.ResponseRecorder) map[string]any {
	t.Helper()

	var body map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response body: %v; body=%s", err, recorder.Body.String())
	}

	return body
}

func testTokenPair() auth_service.TokenPair {
	now := time.Now().UTC()

	return auth_service.TokenPair{
		AccessToken:      "access-token",
		RefreshToken:     "refresh-token",
		AccessExpiresAt:  now.Add(15 * time.Minute),
		RefreshExpiresAt: now.Add(24 * time.Hour),
		DeviceID:         "00000000-0000-0000-0000-000000000001",
	}
}

func TestRegisterReturnsCreatedAndPassesInputToService(t *testing.T) {
	service := &fakeAuthService{
		registerFn: func(ctx context.Context, input auth_service.RegisterInput) (auth_service.UserDTO, auth_service.TokenPair, error) {
			if input.Email != "client@test.com" {
				t.Fatalf("input.Email = %q, want client@test.com", input.Email)
			}
			if input.Password != "123456" {
				t.Fatalf("input.Password = %q, want 123456", input.Password)
			}
			if input.FullName != "Test Client" {
				t.Fatalf("input.FullName = %q, want Test Client", input.FullName)
			}
			if input.Phone != "+79990000001" {
				t.Fatalf("input.Phone = %q, want +79990000001", input.Phone)
			}

			return auth_service.UserDTO{
				ID:       1,
				Email:    "client@test.com",
				FullName: "Test Client",
				Phone:    "+79990000001",
				Role:     core_domain.RoleClient.String(),
			}, testTokenPair(), nil
		},
	}

	handler := newTestHandler(t, service)

	request := httptest.NewRequest(
		http.MethodPost,
		"/auth/register",
		strings.NewReader(`{
			"email":"client@test.com",
			"password":"123456",
			"full_name":"Test Client",
			"phone":"+79990000001"
		}`),
	)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	handler.Register(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusCreated, recorder.Body.String())
	}

	body := decodeJSONBody(t, recorder)
	if body["user"] == nil {
		t.Fatalf("response must contain user")
	}
	if body["tokens"] == nil {
		t.Fatalf("response must contain tokens")
	}
}

func TestRegisterReturnsBadRequestForInvalidJSON(t *testing.T) {
	handler := newTestHandler(t, &fakeAuthService{})

	request := httptest.NewRequest(
		http.MethodPost,
		"/auth/register",
		strings.NewReader(`{"email":"bad-email"}`),
	)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	handler.Register(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusBadRequest, recorder.Body.String())
	}
}

func TestLoginReturnsOK(t *testing.T) {
	service := &fakeAuthService{
		loginFn: func(ctx context.Context, input auth_service.LoginInput) (auth_service.UserDTO, auth_service.TokenPair, error) {
			if input.Email != "client@test.com" {
				t.Fatalf("input.Email = %q, want client@test.com", input.Email)
			}
			if input.Password != "123456" {
				t.Fatalf("input.Password = %q, want 123456", input.Password)
			}

			return auth_service.UserDTO{
				ID:       1,
				Email:    "client@test.com",
				FullName: "Test Client",
				Role:     core_domain.RoleClient.String(),
			}, testTokenPair(), nil
		},
	}

	handler := newTestHandler(t, service)

	request := httptest.NewRequest(
		http.MethodPost,
		"/auth/login",
		strings.NewReader(`{"email":"client@test.com","password":"123456"}`),
	)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	handler.Login(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
}

func TestLoginReturnsUnauthorizedWhenServiceRejectsCredentials(t *testing.T) {
	service := &fakeAuthService{
		loginFn: func(ctx context.Context, input auth_service.LoginInput) (auth_service.UserDTO, auth_service.TokenPair, error) {
			return auth_service.UserDTO{}, auth_service.TokenPair{}, fmt.Errorf("%w: invalid email or password", core_errors.ErrUnauthorized)
		},
	}

	handler := newTestHandler(t, service)

	request := httptest.NewRequest(
		http.MethodPost,
		"/auth/login",
		strings.NewReader(`{"email":"client@test.com","password":"wrong"}`),
	)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	handler.Login(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusUnauthorized, recorder.Body.String())
	}
}

func TestRefreshReturnsOK(t *testing.T) {
	service := &fakeAuthService{
		refreshFn: func(ctx context.Context, input auth_service.RefreshInput) (auth_service.TokenPair, error) {
			if input.RefreshToken != "old-refresh-token" {
				t.Fatalf("input.RefreshToken = %q, want old-refresh-token", input.RefreshToken)
			}
			return testTokenPair(), nil
		},
	}

	handler := newTestHandler(t, service)

	request := httptest.NewRequest(
		http.MethodPost,
		"/auth/refresh",
		strings.NewReader(`{"refresh_token":"old-refresh-token"}`),
	)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	handler.Refresh(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusOK, recorder.Body.String())
	}

	body := decodeJSONBody(t, recorder)
	if body["tokens"] == nil {
		t.Fatalf("response must contain tokens")
	}
}

func TestLogoutReturnsNoContent(t *testing.T) {
	service := &fakeAuthService{
		logoutFn: func(ctx context.Context, input auth_service.LogoutInput) error {
			if input.RefreshToken != "refresh-token" {
				t.Fatalf("input.RefreshToken = %q, want refresh-token", input.RefreshToken)
			}
			return nil
		},
	}

	handler := newTestHandler(t, service)

	request := httptest.NewRequest(
		http.MethodPost,
		"/auth/logout",
		strings.NewReader(`{"refresh_token":"refresh-token"}`),
	)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	handler.Logout(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusNoContent, recorder.Body.String())
	}
}

func TestMeReturnsCurrentUser(t *testing.T) {
	service := &fakeAuthService{
		getMeFn: func(ctx context.Context, userID int64) (auth_service.UserDTO, error) {
			if userID != 10 {
				t.Fatalf("userID = %d, want 10", userID)
			}
			return auth_service.UserDTO{
				ID:       10,
				Email:    "client@test.com",
				FullName: "Test Client",
				Role:     core_domain.RoleClient.String(),
			}, nil
		},
	}

	handler := newTestHandler(t, service)

	request := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	ctx := core_http_middleware.WithUser(
		request.Context(),
		core_http_middleware.CurrentUser{
			ID:   10,
			Role: core_domain.RoleClient.String(),
			JTI:  "test-jti",
		},
	)
	request = request.WithContext(ctx)

	recorder := httptest.NewRecorder()
	handler.Me(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusOK, recorder.Body.String())
	}

	body := decodeJSONBody(t, recorder)
	if body["user"] == nil {
		t.Fatalf("response must contain user")
	}
}

func TestMeReturnsUnauthorizedWithoutCurrentUser(t *testing.T) {
	handler := newTestHandler(t, &fakeAuthService{})

	request := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	recorder := httptest.NewRecorder()

	handler.Me(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusUnauthorized, recorder.Body.String())
	}
}

