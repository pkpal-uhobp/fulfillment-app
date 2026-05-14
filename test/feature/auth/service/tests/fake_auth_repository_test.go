package auth_service_tests

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

type fakeTxManager struct{}

func (fakeTxManager) WithinTransaction(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	return fn(ctx)
}

type fakeAuthRepository struct {
	usersByID      map[int64]core_domain.User
	userIDByEmail map[string]int64
	tokensByJTI   map[uuid.UUID]core_domain.IssuedToken

	nextUserID  int64
	nextTokenID int64
}

func newFakeAuthRepository() *fakeAuthRepository {
	return &fakeAuthRepository{
		usersByID:      make(map[int64]core_domain.User),
		userIDByEmail: make(map[string]int64),
		tokensByJTI:   make(map[uuid.UUID]core_domain.IssuedToken),
		nextUserID:    1,
		nextTokenID:   1,
	}
}

func (r *fakeAuthRepository) CreateUser(
	ctx context.Context,
	user core_domain.User,
) (core_domain.User, error) {
	email := strings.ToLower(strings.TrimSpace(user.Email))

	if _, exists := r.userIDByEmail[email]; exists {
		return core_domain.User{}, fmt.Errorf("%w: user already exists", core_errors.ErrConflict)
	}

	user.ID = r.nextUserID
	r.nextUserID++

	r.usersByID[user.ID] = user
	r.userIDByEmail[email] = user.ID

	return user, nil
}

func (r *fakeAuthRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (core_domain.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	userID, exists := r.userIDByEmail[email]
	if !exists {
		return core_domain.User{}, core_errors.ErrNotFound
	}

	return r.usersByID[userID], nil
}

func (r *fakeAuthRepository) GetUserByID(
	ctx context.Context,
	userID int64,
) (core_domain.User, error) {
	user, exists := r.usersByID[userID]
	if !exists {
		return core_domain.User{}, core_errors.ErrNotFound
	}

	return user, nil
}

func (r *fakeAuthRepository) CreateIssuedToken(
	ctx context.Context,
	token core_domain.IssuedToken,
) error {
	token.ID = r.nextTokenID
	r.nextTokenID++

	r.tokensByJTI[token.JTI] = token

	return nil
}

func (r *fakeAuthRepository) GetIssuedTokenByJTI(
	ctx context.Context,
	jti uuid.UUID,
) (core_domain.IssuedToken, error) {
	token, exists := r.tokensByJTI[jti]
	if !exists {
		return core_domain.IssuedToken{}, core_errors.ErrNotFound
	}

	return token, nil
}

func (r *fakeAuthRepository) RevokeTokenByJTI(
	ctx context.Context,
	jti uuid.UUID,
	reason string,
) error {
	token, exists := r.tokensByJTI[jti]
	if !exists || token.Revoked {
		return core_errors.ErrNotFound
	}

	token.Revoked = true
	r.tokensByJTI[jti] = token

	return nil
}

func (r *fakeAuthRepository) RevokeActiveTokensByDevice(
	ctx context.Context,
	userID int64,
	deviceID uuid.UUID,
	reason string,
) error {
	now := time.Now().UTC()

	for jti, token := range r.tokensByJTI {
		if token.UserID == userID &&
			token.DeviceID == deviceID &&
			!token.Revoked &&
			now.Before(token.ExpiresAt) {
			token.Revoked = true
			r.tokensByJTI[jti] = token
		}
	}

	return nil
}
