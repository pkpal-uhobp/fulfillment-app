package users_service

import (
	"context"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

type UsersRepository interface {
	ListUsers(ctx context.Context, filter core_domain.UserFilter) ([]core_domain.User, error)
	CreateUser(ctx context.Context, user core_domain.User) (core_domain.User, error)
	PatchUser(ctx context.Context, userID int64, patch core_domain.UserPatch) (core_domain.User, error)
	BlockUser(ctx context.Context, userID int64, reason *string) error
	DeactivateUser(ctx context.Context, userID int64) error
}

type UsersService struct {
	repo UsersRepository
}

func NewUsersService(repo UsersRepository) *UsersService {
	return &UsersService{repo: repo}
}
