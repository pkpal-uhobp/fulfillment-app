package users_service

import (
	"context"
	"fmt"
	"net/mail"
	"strings"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (s *UsersService) PatchUser(
	ctx context.Context,
	actorID int64,
	actorRole string,
	userID int64,
	input PatchUserInput,
) (UserDTO, error) {
	if err := requireAdmin(actorRole); err != nil {
		return UserDTO{}, err
	}
	if userID <= 0 {
		return UserDTO{}, fmt.Errorf("%w: invalid user id", core_errors.ErrInvalidArgument)
	}

	var patch core_domain.UserPatch

	if input.Email != nil {
		value := strings.ToLower(strings.TrimSpace(*input.Email))
		if value == "" {
			return UserDTO{}, fmt.Errorf("%w: email is required", core_errors.ErrInvalidArgument)
		}
		if strings.Contains(value, " ") {
			return UserDTO{}, fmt.Errorf("%w: invalid email", core_errors.ErrInvalidArgument)
		}
		if _, err := mail.ParseAddress(value); err != nil {
			return UserDTO{}, fmt.Errorf("%w: invalid email", core_errors.ErrInvalidArgument)
		}
		patch.Email = &value
	}

	if input.Password != nil {
		value := strings.TrimSpace(*input.Password)
		if value != "" {
			if len([]rune(value)) < 6 {
				return UserDTO{}, fmt.Errorf("%w: password is too short", core_errors.ErrInvalidArgument)
			}
			passwordHash, err := hashPassword(value)
			if err != nil {
				return UserDTO{}, fmt.Errorf("hash password: %w", err)
			}
			patch.PasswordHash = &passwordHash
		}
	}

	if input.FullName != nil {
		value := strings.TrimSpace(*input.FullName)
		if value == "" {
			return UserDTO{}, fmt.Errorf("%w: full name is required", core_errors.ErrInvalidArgument)
		}
		patch.FullName = &value
	}

	if input.Phone != nil {
		patch.PhoneProvided = true
		value := strings.TrimSpace(*input.Phone)
		if value != "" {
			patch.Phone = &value
		}
	}

	if input.Role != nil {
		role, err := validateRole(*input.Role)
		if err != nil {
			return UserDTO{}, err
		}
		if userID == actorID {
			return UserDTO{}, fmt.Errorf("%w: cannot change own role", core_errors.ErrForbidden)
		}
		patch.Role = &role
	}

	if input.IsActive != nil {
		if userID == actorID && !*input.IsActive {
			return UserDTO{}, fmt.Errorf("%w: cannot deactivate yourself", core_errors.ErrForbidden)
		}
		patch.IsActive = input.IsActive
	}

	if input.IsBlocked != nil {
		if userID == actorID && *input.IsBlocked {
			return UserDTO{}, fmt.Errorf("%w: cannot block yourself", core_errors.ErrForbidden)
		}
		patch.IsBlocked = input.IsBlocked
	}

	updated, err := s.repo.PatchUser(ctx, userID, patch)
	if err != nil {
		return UserDTO{}, err
	}

	return toUserDTO(updated), nil
}
