package users_service

import (
	"context"
	"fmt"
	"strings"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (s *UsersService) CreateUser(
	ctx context.Context,
	actorRole string,
	input CreateUserInput,
) (UserDTO, error) {
	if err := requireAdmin(actorRole); err != nil {
		return UserDTO{}, err
	}

	password := strings.TrimSpace(input.Password)
	if password == "" {
		return UserDTO{}, fmt.Errorf("%w: password is required", core_errors.ErrInvalidArgument)
	}
	if len([]rune(password)) < 6 {
		return UserDTO{}, fmt.Errorf("%w: password is too short", core_errors.ErrInvalidArgument)
	}

	role, err := validateRole(input.Role)
	if err != nil {
		return UserDTO{}, err
	}

	passwordHash, err := hashPassword(password)
	if err != nil {
		return UserDTO{}, fmt.Errorf("hash password: %w", err)
	}

	var phone *string
	if strings.TrimSpace(input.Phone) != "" {
		value := strings.TrimSpace(input.Phone)
		phone = &value
	}

	user, err := core_domain.NewUser(
		input.Email,
		passwordHash,
		input.FullName,
		phone,
		role,
	)
	if err != nil {
		return UserDTO{}, err
	}

	created, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return UserDTO{}, err
	}
	return toUserDTO(created), nil
}
