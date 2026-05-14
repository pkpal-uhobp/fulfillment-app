//go:build integration

package auth_postgres_tests

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
	core_postgres_pool "github.com/pkpal-uhobp/fulfillment-app/internal/core/repository/pool"
	core_postgres_tx "github.com/pkpal-uhobp/fulfillment-app/internal/core/repository/tx"
	auth_repository_postgres "github.com/pkpal-uhobp/fulfillment-app/internal/features/auth/repository/postgres"
)

func newIntegrationRepository(t *testing.T) (*auth_repository_postgres.AuthRepository, *core_postgres_pool.ConnectionPool) {
	t.Helper()

	requiredEnv := []string{
		"POSTGRES_HOST",
		"POSTGRES_USER",
		"POSTGRES_PASSWORD",
		"POSTGRES_DB",
	}
	for _, envName := range requiredEnv {
		if os.Getenv(envName) == "" {
			t.Skipf("skip integration test: %s is not set", envName)
		}
	}

	config := core_postgres_pool.Config{
		Host:              os.Getenv("POSTGRES_HOST"),
		Port:              envOrDefault("POSTGRES_PORT", "5432"),
		User:              os.Getenv("POSTGRES_USER"),
		Password:          os.Getenv("POSTGRES_PASSWORD"),
		Database:          os.Getenv("POSTGRES_DB"),
		SSLMode:           envOrDefault("POSTGRES_SSL_MODE", "disable"),
		MaxConns:          4,
		MinConns:          1,
		MaxConnLifetime:   time.Hour,
		MaxConnIdleTime:   30 * time.Minute,
		HealthCheckPeriod: time.Minute,
		ConnectTimeout:    5 * time.Second,
		QueryTimeout:      5 * time.Second,
	}

	pool, err := core_postgres_pool.NewConnectionPool(context.Background(), config)
	if err != nil {
		t.Fatalf("connect postgres: %v", err)
	}

	t.Cleanup(pool.Close)

	_, err = pool.Exec(context.Background(), "TRUNCATE TABLE issued_tokens, users RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatalf("truncate auth tables: %v", err)
	}

	tx := core_postgres_tx.NewTx(pool)
	repo := auth_repository_postgres.NewAuthRepository(tx)

	return repo, pool
}

func envOrDefault(name string, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

func newDomainUser(t *testing.T, email string) core_domain.User {
	t.Helper()

	user, err := core_domain.NewUser(
		email,
		"hashed-password",
		"Test User",
		nil,
		core_domain.RoleClient,
	)
	if err != nil {
		t.Fatalf("create domain user: %v", err)
	}

	return user
}

func TestCreateAndGetUser(t *testing.T) {
	repo, _ := newIntegrationRepository(t)
	ctx := context.Background()

	email := fmt.Sprintf("client-%s@test.com", uuid.NewString())
	createdUser, err := repo.CreateUser(ctx, newDomainUser(t, email))
	if err != nil {
		t.Fatalf("CreateUser returned error: %v", err)
	}

	if createdUser.ID == 0 {
		t.Fatalf("created user ID must not be zero")
	}

	userByEmail, err := repo.GetUserByEmail(ctx, email)
	if err != nil {
		t.Fatalf("GetUserByEmail returned error: %v", err)
	}

	if userByEmail.ID != createdUser.ID {
		t.Fatalf("userByEmail.ID = %d, want %d", userByEmail.ID, createdUser.ID)
	}

	userByID, err := repo.GetUserByID(ctx, createdUser.ID)
	if err != nil {
		t.Fatalf("GetUserByID returned error: %v", err)
	}

	if userByID.Email != email {
		t.Fatalf("userByID.Email = %q, want %q", userByID.Email, email)
	}
}

func TestCreateUserRejectsDuplicateEmail(t *testing.T) {
	repo, _ := newIntegrationRepository(t)
	ctx := context.Background()

	email := fmt.Sprintf("client-%s@test.com", uuid.NewString())

	_, err := repo.CreateUser(ctx, newDomainUser(t, email))
	if err != nil {
		t.Fatalf("first CreateUser returned error: %v", err)
	}

	_, err = repo.CreateUser(ctx, newDomainUser(t, email))
	if !errors.Is(err, core_errors.ErrConflict) {
		t.Fatalf("err = %v, want ErrConflict", err)
	}
}

func TestCreateGetAndRevokeIssuedToken(t *testing.T) {
	repo, _ := newIntegrationRepository(t)
	ctx := context.Background()

	email := fmt.Sprintf("client-%s@test.com", uuid.NewString())
	createdUser, err := repo.CreateUser(ctx, newDomainUser(t, email))
	if err != nil {
		t.Fatalf("CreateUser returned error: %v", err)
	}

	jti := uuid.New()
	deviceID := uuid.New()

	err = repo.CreateIssuedToken(ctx, core_domain.IssuedToken{
		UserID:    createdUser.ID,
		JTI:       jti,
		TokenType: core_domain.TokenTypeAccess,
		DeviceID:  deviceID,
		ExpiresAt: time.Now().UTC().Add(time.Hour),
	})
	if err != nil {
		t.Fatalf("CreateIssuedToken returned error: %v", err)
	}

	issuedToken, err := repo.GetIssuedTokenByJTI(ctx, jti)
	if err != nil {
		t.Fatalf("GetIssuedTokenByJTI returned error: %v", err)
	}

	if issuedToken.UserID != createdUser.ID {
		t.Fatalf("issuedToken.UserID = %d, want %d", issuedToken.UserID, createdUser.ID)
	}

	err = repo.RevokeTokenByJTI(ctx, jti, "test revoke")
	if err != nil {
		t.Fatalf("RevokeTokenByJTI returned error: %v", err)
	}

	revokedToken, err := repo.GetIssuedTokenByJTI(ctx, jti)
	if err != nil {
		t.Fatalf("GetIssuedTokenByJTI after revoke returned error: %v", err)
	}

	if !revokedToken.Revoked {
		t.Fatalf("token must be revoked")
	}
}
