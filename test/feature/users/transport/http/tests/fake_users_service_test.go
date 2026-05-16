package users_transport_http_tests

import (
	"context"

	users_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/users/service"
)

type fakeUsersService struct {
	listUsersFn      func(ctx context.Context, actorRole string, filter users_service.UserFilter) ([]users_service.UserDTO, error)
	createUserFn     func(ctx context.Context, actorRole string, input users_service.CreateUserInput) (users_service.UserDTO, error)
	patchUserFn      func(ctx context.Context, actorID int64, actorRole string, userID int64, input users_service.PatchUserInput) (users_service.UserDTO, error)
	blockUserFn      func(ctx context.Context, actorID int64, actorRole string, userID int64, input users_service.BlockUserInput) error
	deactivateUserFn func(ctx context.Context, actorID int64, actorRole string, userID int64) error
}

func (f *fakeUsersService) ListUsers(ctx context.Context, actorRole string, filter users_service.UserFilter) ([]users_service.UserDTO, error) {
	if f.listUsersFn != nil {
		return f.listUsersFn(ctx, actorRole, filter)
	}
	return nil, nil
}

func (f *fakeUsersService) CreateUser(ctx context.Context, actorRole string, input users_service.CreateUserInput) (users_service.UserDTO, error) {
	if f.createUserFn != nil {
		return f.createUserFn(ctx, actorRole, input)
	}
	return users_service.UserDTO{}, nil
}

func (f *fakeUsersService) PatchUser(ctx context.Context, actorID int64, actorRole string, userID int64, input users_service.PatchUserInput) (users_service.UserDTO, error) {
	if f.patchUserFn != nil {
		return f.patchUserFn(ctx, actorID, actorRole, userID, input)
	}
	return users_service.UserDTO{}, nil
}

func (f *fakeUsersService) BlockUser(ctx context.Context, actorID int64, actorRole string, userID int64, input users_service.BlockUserInput) error {
	if f.blockUserFn != nil {
		return f.blockUserFn(ctx, actorID, actorRole, userID, input)
	}
	return nil
}

func (f *fakeUsersService) DeactivateUser(ctx context.Context, actorID int64, actorRole string, userID int64) error {
	if f.deactivateUserFn != nil {
		return f.deactivateUserFn(ctx, actorID, actorRole, userID)
	}
	return nil
}
