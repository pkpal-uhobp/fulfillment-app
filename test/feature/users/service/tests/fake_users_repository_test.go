package users_service_tests

import (
	"context"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

type fakeUsersRepository struct {
	listUsersFn      func(ctx context.Context, filter core_domain.UserFilter) ([]core_domain.User, error)
	createUserFn     func(ctx context.Context, user core_domain.User) (core_domain.User, error)
	patchUserFn      func(ctx context.Context, userID int64, patch core_domain.UserPatch) (core_domain.User, error)
	blockUserFn      func(ctx context.Context, userID int64, reason *string) error
	deactivateUserFn func(ctx context.Context, userID int64) error
}

func (f *fakeUsersRepository) ListUsers(ctx context.Context, filter core_domain.UserFilter) ([]core_domain.User, error) {
	if f.listUsersFn != nil {
		return f.listUsersFn(ctx, filter)
	}
	return nil, nil
}

func (f *fakeUsersRepository) CreateUser(ctx context.Context, user core_domain.User) (core_domain.User, error) {
	if f.createUserFn != nil {
		return f.createUserFn(ctx, user)
	}
	return user, nil
}

func (f *fakeUsersRepository) PatchUser(ctx context.Context, userID int64, patch core_domain.UserPatch) (core_domain.User, error) {
	if f.patchUserFn != nil {
		return f.patchUserFn(ctx, userID, patch)
	}
	return core_domain.User{}, nil
}

func (f *fakeUsersRepository) BlockUser(ctx context.Context, userID int64, reason *string) error {
	if f.blockUserFn != nil {
		return f.blockUserFn(ctx, userID, reason)
	}
	return nil
}

func (f *fakeUsersRepository) DeactivateUser(ctx context.Context, userID int64) error {
	if f.deactivateUserFn != nil {
		return f.deactivateUserFn(ctx, userID)
	}
	return nil
}
