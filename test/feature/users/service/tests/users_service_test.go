package users_service_tests

import (
	"context"
	"errors"
	"testing"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
	users_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/users/service"
)

func TestCreateUserSuccess(t *testing.T) {
	repo := &fakeUsersRepository{
		createUserFn: func(_ context.Context, user core_domain.User) (core_domain.User, error) {
			if user.Email != "worker@test.com" {
				t.Fatalf("Email = %q, want worker@test.com", user.Email)
			}
			if user.PasswordHash == "123456" || user.PasswordHash == "" {
				t.Fatalf("password must be hashed")
			}
			if user.Role != core_domain.RoleWorker {
				t.Fatalf("Role = %q, want worker", user.Role)
			}
			user.ID = 10
			return user, nil
		},
	}
	service := users_service.NewUsersService(repo)

	created, err := service.CreateUser(context.Background(), core_domain.RoleAdmin.String(), users_service.CreateUserInput{
		Email:    "worker@test.com",
		Password: "123456",
		FullName: "Worker Test",
		Phone:    "+79990000001",
		Role:     core_domain.RoleWorker.String(),
	})
	if err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}
	if created.ID != 10 || created.Role != core_domain.RoleWorker.String() {
		t.Fatalf("unexpected created user: %+v", created)
	}
}

func TestCreateUserRejectsNonAdmin(t *testing.T) {
	service := users_service.NewUsersService(&fakeUsersRepository{})
	_, err := service.CreateUser(context.Background(), core_domain.RoleLogist.String(), users_service.CreateUserInput{})
	if !errors.Is(err, core_errors.ErrForbidden) {
		t.Fatalf("error = %v, want forbidden", err)
	}
}

func TestCreateUserRejectsShortPassword(t *testing.T) {
	service := users_service.NewUsersService(&fakeUsersRepository{})
	_, err := service.CreateUser(context.Background(), core_domain.RoleAdmin.String(), users_service.CreateUserInput{
		Email:    "worker@test.com",
		Password: "123",
		FullName: "Worker Test",
		Role:     core_domain.RoleWorker.String(),
	})
	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("error = %v, want invalid argument", err)
	}
}

func TestListUsersNormalizesPaginationAndFilter(t *testing.T) {
	repo := &fakeUsersRepository{
		listUsersFn: func(_ context.Context, filter core_domain.UserFilter) ([]core_domain.User, error) {
			if filter.Role != core_domain.RoleWorker.String() {
				t.Fatalf("Role = %q, want worker", filter.Role)
			}
			if filter.Page != 1 {
				t.Fatalf("Page = %d, want 1", filter.Page)
			}
			if filter.Limit != 20 {
				t.Fatalf("Limit = %d, want 20", filter.Limit)
			}
			return []core_domain.User{{ID: 1, Email: "w@test.com", FullName: "Worker", Role: core_domain.RoleWorker, IsActive: true}}, nil
		},
	}
	service := users_service.NewUsersService(repo)

	users, err := service.ListUsers(context.Background(), core_domain.RoleAdmin.String(), users_service.UserFilter{Role: " worker "})
	if err != nil {
		t.Fatalf("ListUsers() error = %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("len(users) = %d, want 1", len(users))
	}
}

func TestPatchUserRejectsOwnRoleChange(t *testing.T) {
	role := core_domain.RoleWorker.String()
	service := users_service.NewUsersService(&fakeUsersRepository{})
	_, err := service.PatchUser(context.Background(), 1, core_domain.RoleAdmin.String(), 1, users_service.PatchUserInput{Role: &role})
	if !errors.Is(err, core_errors.ErrForbidden) {
		t.Fatalf("error = %v, want forbidden", err)
	}
}

func TestBlockUserRejectsSelfBlock(t *testing.T) {
	service := users_service.NewUsersService(&fakeUsersRepository{})
	err := service.BlockUser(context.Background(), 1, core_domain.RoleAdmin.String(), 1, users_service.BlockUserInput{})
	if !errors.Is(err, core_errors.ErrForbidden) {
		t.Fatalf("error = %v, want forbidden", err)
	}
}
