package auth_service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (s *AuthService) Register(
	ctx context.Context,
	input RegisterInput,
) (UserDTO, TokenPair, error) {
	email := normalizeEmail(input.Email)

	if email == "" ||
		strings.TrimSpace(input.Password) == "" ||
		strings.TrimSpace(input.FullName) == "" {
		return UserDTO{}, TokenPair{}, fmt.Errorf(
			"%w: required fields are empty",
			core_errors.ErrInvalidArgument,
		)
	}

	if len([]rune(input.Password)) < 6 {
		return UserDTO{}, TokenPair{}, fmt.Errorf(
			"%w: password is too short",
			core_errors.ErrInvalidArgument,
		)
	}

	passwordHash, err := hashPassword(input.Password)
	if err != nil {
		return UserDTO{}, TokenPair{}, fmt.Errorf("hash password: %w", err)
	}

	var (
		user   core_domain.User
		tokens TokenPair
	)

	err = s.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		var phone *string

		if strings.TrimSpace(input.Phone) != "" {
			p := strings.TrimSpace(input.Phone)
			phone = &p
		}

		createdUser, err := s.repo.CreateUser(ctx, core_domain.User{
			Email:        email,
			PasswordHash: passwordHash,
			FullName:     strings.TrimSpace(input.FullName),
			Phone:        phone,
			Role:         core_domain.RoleClient,
		})
		if err != nil {
			return err
		}

		user = createdUser

		tokenPair, err := s.issueTokenPair(ctx, user, uuid.New())
		if err != nil {
			return err
		}

		tokens = tokenPair

		return nil
	})
	if err != nil {
		return UserDTO{}, TokenPair{}, err
	}

	return toUserDTO(user), tokens, nil
}
