package auth_service

import (
	"context"

	"github.com/google/uuid"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

type TransactionManager interface {
	WithinTransaction(
		ctx context.Context,
		fn func(ctx context.Context) error,
	) error
}

type AuthRepository interface {
	CreateUser(
		ctx context.Context,
		user core_domain.User,
	) (core_domain.User, error)

	GetUserByEmail(
		ctx context.Context,
		email string,
	) (core_domain.User, error)

	GetUserByID(
		ctx context.Context,
		userID int64,
	) (core_domain.User, error)

	CreateIssuedToken(
		ctx context.Context,
		token core_domain.IssuedToken,
	) error

	GetIssuedTokenByJTI(
		ctx context.Context,
		jti uuid.UUID,
	) (core_domain.IssuedToken, error)

	RevokeTokenByJTI(
		ctx context.Context,
		jti uuid.UUID,
		reason string,
	) error

	RevokeActiveTokensByDevice(
		ctx context.Context,
		userID int64,
		deviceID uuid.UUID,
		reason string,
	) error
}

type AuthService struct {
	tx     TransactionManager
	repo   AuthRepository
	config Config
}

func NewAuthService(
	tx TransactionManager,
	repo AuthRepository,
	config Config,
) *AuthService {
	return &AuthService{
		tx:     tx,
		repo:   repo,
		config: config,
	}
}
