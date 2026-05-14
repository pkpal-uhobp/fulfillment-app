package auth_service_tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
	auth_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/auth/service"
	"golang.org/x/crypto/bcrypt"
)

func newTestAuthService(repo *fakeAuthRepository) *auth_service.AuthService {
	return auth_service.NewAuthService(
		fakeTxManager{},
		repo,
		auth_service.Config{
			JWTSecret:  "test-secret",
			AccessTTL:  15 * time.Minute,
			RefreshTTL: 24 * time.Hour,
		},
	)
}

func mustCreateUser(
	t *testing.T,
	repo *fakeAuthRepository,
	email string,
	password string,
	role core_domain.Role,
) core_domain.User {
	t.Helper()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	user, err := core_domain.NewUser(
		email,
		string(hash),
		"Test User",
		nil,
		role,
	)
	if err != nil {
		t.Fatalf("create domain user: %v", err)
	}

	createdUser, err := repo.CreateUser(context.Background(), user)
	if err != nil {
		t.Fatalf("save user: %v", err)
	}

	return createdUser
}

func TestRegisterSuccess(t *testing.T) {
	repo := newFakeAuthRepository()
	service := newTestAuthService(repo)

	user, tokens, err := service.Register(
		context.Background(),
		auth_service.RegisterInput{
			Email:    "CLIENT@TEST.COM",
			Password: "123456",
			FullName: "Test Client",
			Phone:    "+79990000001",
		},
	)

	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	if user.ID == 0 {
		t.Fatalf("user.ID must not be zero")
	}

	if user.Email != "client@test.com" {
		t.Fatalf("user.Email = %q, want client@test.com", user.Email)
	}

	if user.Role != core_domain.RoleClient.String() {
		t.Fatalf("user.Role = %q, want client", user.Role)
	}

	if tokens.AccessToken == "" {
		t.Fatalf("access token must not be empty")
	}

	if tokens.RefreshToken == "" {
		t.Fatalf("refresh token must not be empty")
	}

	if tokens.DeviceID == "" {
		t.Fatalf("device id must not be empty")
	}

	if len(repo.tokensByJTI) != 2 {
		t.Fatalf("tokens count = %d, want 2", len(repo.tokensByJTI))
	}
}

func TestRegisterRejectsShortPassword(t *testing.T) {
	repo := newFakeAuthRepository()
	service := newTestAuthService(repo)

	_, _, err := service.Register(
		context.Background(),
		auth_service.RegisterInput{
			Email:    "client@test.com",
			Password: "123",
			FullName: "Test Client",
		},
	)

	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("err = %v, want ErrInvalidArgument", err)
	}
}

func TestRegisterRejectsInvalidEmail(t *testing.T) {
	repo := newFakeAuthRepository()
	service := newTestAuthService(repo)

	_, _, err := service.Register(
		context.Background(),
		auth_service.RegisterInput{
			Email:    "bad-email",
			Password: "123456",
			FullName: "Test Client",
		},
	)

	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("err = %v, want ErrInvalidArgument", err)
	}
}

func TestRegisterRejectsDuplicateEmail(t *testing.T) {
	repo := newFakeAuthRepository()
	service := newTestAuthService(repo)

	input := auth_service.RegisterInput{
		Email:    "client@test.com",
		Password: "123456",
		FullName: "Test Client",
	}

	_, _, err := service.Register(context.Background(), input)
	if err != nil {
		t.Fatalf("first Register returned error: %v", err)
	}

	_, _, err = service.Register(context.Background(), input)
	if !errors.Is(err, core_errors.ErrConflict) {
		t.Fatalf("err = %v, want ErrConflict", err)
	}
}

func TestLoginSuccess(t *testing.T) {
	repo := newFakeAuthRepository()
	service := newTestAuthService(repo)

	createdUser := mustCreateUser(
		t,
		repo,
		"client@test.com",
		"123456",
		core_domain.RoleClient,
	)

	user, tokens, err := service.Login(
		context.Background(),
		auth_service.LoginInput{
			Email:    " CLIENT@TEST.COM ",
			Password: "123456",
		},
	)

	if err != nil {
		t.Fatalf("Login returned error: %v", err)
	}

	if user.ID != createdUser.ID {
		t.Fatalf("user.ID = %d, want %d", user.ID, createdUser.ID)
	}

	if tokens.AccessToken == "" {
		t.Fatalf("access token must not be empty")
	}

	if tokens.RefreshToken == "" {
		t.Fatalf("refresh token must not be empty")
	}

	if _, err := uuid.Parse(tokens.DeviceID); err != nil {
		t.Fatalf("invalid device id: %v", err)
	}

	if len(repo.tokensByJTI) != 2 {
		t.Fatalf("tokens count = %d, want 2", len(repo.tokensByJTI))
	}
}

func TestLoginRejectsWrongPassword(t *testing.T) {
	repo := newFakeAuthRepository()
	service := newTestAuthService(repo)

	mustCreateUser(
		t,
		repo,
		"client@test.com",
		"123456",
		core_domain.RoleClient,
	)

	_, _, err := service.Login(
		context.Background(),
		auth_service.LoginInput{
			Email:    "client@test.com",
			Password: "wrong-password",
		},
	)

	if !errors.Is(err, core_errors.ErrUnauthorized) {
		t.Fatalf("err = %v, want ErrUnauthorized", err)
	}
}

func TestLoginRejectsUnknownUser(t *testing.T) {
	repo := newFakeAuthRepository()
	service := newTestAuthService(repo)

	_, _, err := service.Login(
		context.Background(),
		auth_service.LoginInput{
			Email:    "unknown@test.com",
			Password: "123456",
		},
	)

	if !errors.Is(err, core_errors.ErrUnauthorized) {
		t.Fatalf("err = %v, want ErrUnauthorized", err)
	}
}

func TestLoginRejectsInvalidDeviceID(t *testing.T) {
	repo := newFakeAuthRepository()
	service := newTestAuthService(repo)

	mustCreateUser(
		t,
		repo,
		"client@test.com",
		"123456",
		core_domain.RoleClient,
	)

	_, _, err := service.Login(
		context.Background(),
		auth_service.LoginInput{
			Email:    "client@test.com",
			Password: "123456",
			DeviceID: "bad-device-id",
		},
	)

	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("err = %v, want ErrInvalidArgument", err)
	}
}

func TestGetMeSuccess(t *testing.T) {
	repo := newFakeAuthRepository()
	service := newTestAuthService(repo)

	createdUser := mustCreateUser(
		t,
		repo,
		"client@test.com",
		"123456",
		core_domain.RoleClient,
	)

	user, err := service.GetMe(context.Background(), createdUser.ID)
	if err != nil {
		t.Fatalf("GetMe returned error: %v", err)
	}

	if user.ID != createdUser.ID {
		t.Fatalf("user.ID = %d, want %d", user.ID, createdUser.ID)
	}

	if user.Email != "client@test.com" {
		t.Fatalf("user.Email = %q, want client@test.com", user.Email)
	}
}

func TestGetMeRejectsInvalidUserID(t *testing.T) {
	repo := newFakeAuthRepository()
	service := newTestAuthService(repo)

	_, err := service.GetMe(context.Background(), 0)

	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("err = %v, want ErrInvalidArgument", err)
	}
}

func TestGetMeRejectsBlockedUser(t *testing.T) {
	repo := newFakeAuthRepository()
	service := newTestAuthService(repo)

	createdUser := mustCreateUser(
		t,
		repo,
		"client@test.com",
		"123456",
		core_domain.RoleClient,
	)

	createdUser.IsBlocked = true
	repo.usersByID[createdUser.ID] = createdUser

	_, err := service.GetMe(context.Background(), createdUser.ID)

	if !errors.Is(err, core_errors.ErrForbidden) {
		t.Fatalf("err = %v, want ErrForbidden", err)
	}
}

func TestVerifyAccessTokenSuccess(t *testing.T) {
	repo := newFakeAuthRepository()
	service := newTestAuthService(repo)

	user, tokens, err := service.Register(
		context.Background(),
		auth_service.RegisterInput{
			Email:    "client@test.com",
			Password: "123456",
			FullName: "Test Client",
		},
	)
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	claims, err := service.VerifyAccessToken(
		context.Background(),
		tokens.AccessToken,
	)

	if err != nil {
		t.Fatalf("VerifyAccessToken returned error: %v", err)
	}

	if claims.UserID != user.ID {
		t.Fatalf("claims.UserID = %d, want %d", claims.UserID, user.ID)
	}

	if claims.Role != core_domain.RoleClient.String() {
		t.Fatalf("claims.Role = %q, want client", claims.Role)
	}

	if claims.JTI == "" {
		t.Fatalf("claims.JTI must not be empty")
	}

	if claims.DeviceID == "" {
		t.Fatalf("claims.DeviceID must not be empty")
	}
}

func TestVerifyAccessTokenRejectsRefreshToken(t *testing.T) {
	repo := newFakeAuthRepository()
	service := newTestAuthService(repo)

	_, tokens, err := service.Register(
		context.Background(),
		auth_service.RegisterInput{
			Email:    "client@test.com",
			Password: "123456",
			FullName: "Test Client",
		},
	)
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	_, err = service.VerifyAccessToken(
		context.Background(),
		tokens.RefreshToken,
	)

	if !errors.Is(err, core_errors.ErrUnauthorized) {
		t.Fatalf("err = %v, want ErrUnauthorized", err)
	}
}

func TestVerifyAccessTokenRejectsRevokedToken(t *testing.T) {
	repo := newFakeAuthRepository()
	service := newTestAuthService(repo)

	_, tokens, err := service.Register(
		context.Background(),
		auth_service.RegisterInput{
			Email:    "client@test.com",
			Password: "123456",
			FullName: "Test Client",
		},
	)
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	claims, err := service.VerifyAccessToken(
		context.Background(),
		tokens.AccessToken,
	)
	if err != nil {
		t.Fatalf("first VerifyAccessToken returned error: %v", err)
	}

	jti := uuid.MustParse(claims.JTI)

	err = repo.RevokeTokenByJTI(context.Background(), jti, "test revoke")
	if err != nil {
		t.Fatalf("revoke token: %v", err)
	}

	_, err = service.VerifyAccessToken(
		context.Background(),
		tokens.AccessToken,
	)

	if !errors.Is(err, core_errors.ErrUnauthorized) {
		t.Fatalf("err = %v, want ErrUnauthorized", err)
	}
}

func TestRefreshSuccessRotatesRefreshToken(t *testing.T) {
	repo := newFakeAuthRepository()
	service := newTestAuthService(repo)

	_, tokens, err := service.Register(
		context.Background(),
		auth_service.RegisterInput{
			Email:    "client@test.com",
			Password: "123456",
			FullName: "Test Client",
		},
	)
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	newTokens, err := service.Refresh(
		context.Background(),
		auth_service.RefreshInput{
			RefreshToken: tokens.RefreshToken,
		},
	)

	if err != nil {
		t.Fatalf("Refresh returned error: %v", err)
	}

	if newTokens.AccessToken == "" {
		t.Fatalf("new access token must not be empty")
	}

	if newTokens.RefreshToken == "" {
		t.Fatalf("new refresh token must not be empty")
	}

	if newTokens.RefreshToken == tokens.RefreshToken {
		t.Fatalf("new refresh token must be different from old refresh token")
	}

	if len(repo.tokensByJTI) != 4 {
		t.Fatalf("tokens count = %d, want 4", len(repo.tokensByJTI))
	}

	_, err = service.Refresh(
		context.Background(),
		auth_service.RefreshInput{
			RefreshToken: tokens.RefreshToken,
		},
	)

	if !errors.Is(err, core_errors.ErrUnauthorized) {
		t.Fatalf("old refresh token err = %v, want ErrUnauthorized", err)
	}
}

func TestRefreshRejectsAccessToken(t *testing.T) {
	repo := newFakeAuthRepository()
	service := newTestAuthService(repo)

	_, tokens, err := service.Register(
		context.Background(),
		auth_service.RegisterInput{
			Email:    "client@test.com",
			Password: "123456",
			FullName: "Test Client",
		},
	)
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	_, err = service.Refresh(
		context.Background(),
		auth_service.RefreshInput{
			RefreshToken: tokens.AccessToken,
		},
	)

	if !errors.Is(err, core_errors.ErrUnauthorized) {
		t.Fatalf("err = %v, want ErrUnauthorized", err)
	}
}

func TestLogoutSuccessRevokesDeviceTokens(t *testing.T) {
	repo := newFakeAuthRepository()
	service := newTestAuthService(repo)

	_, tokens, err := service.Register(
		context.Background(),
		auth_service.RegisterInput{
			Email:    "client@test.com",
			Password: "123456",
			FullName: "Test Client",
		},
	)
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	err = service.Logout(
		context.Background(),
		auth_service.LogoutInput{
			RefreshToken: tokens.RefreshToken,
		},
	)

	if err != nil {
		t.Fatalf("Logout returned error: %v", err)
	}

	_, err = service.VerifyAccessToken(
		context.Background(),
		tokens.AccessToken,
	)

	if !errors.Is(err, core_errors.ErrUnauthorized) {
		t.Fatalf("err = %v, want ErrUnauthorized after logout", err)
	}
}

func TestLogoutRejectsEmptyRefreshToken(t *testing.T) {
	repo := newFakeAuthRepository()
	service := newTestAuthService(repo)

	err := service.Logout(
		context.Background(),
		auth_service.LogoutInput{
			RefreshToken: "",
		},
	)

	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("err = %v, want ErrInvalidArgument", err)
	}
}
