package auth_service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"

	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (s *AuthService) Login(
	ctx context.Context,
	input LoginInput,
) (UserDTO, TokenPair, error) {
	email := normalizeEmail(input.Email)

	if email == "" || strings.TrimSpace(input.Password) == "" {
		return UserDTO{}, TokenPair{}, fmt.Errorf(
			"%w: empty email or password",
			core_errors.ErrInvalidArgument,
		)
	}

	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return UserDTO{}, TokenPair{}, fmt.Errorf(
				"%w: invalid email or password",
				core_errors.ErrUnauthorized,
			)
		}

		return UserDTO{}, TokenPair{}, err
	}

	if err := checkUserAvailable(user); err != nil {
		return UserDTO{}, TokenPair{}, err
	}

	if !checkPassword(input.Password, user.PasswordHash) {
		return UserDTO{}, TokenPair{}, fmt.Errorf(
			"%w: invalid email or password",
			core_errors.ErrUnauthorized,
		)
	}

	deviceID := uuid.New()

	if strings.TrimSpace(input.DeviceID) != "" {
		parsedDeviceID, err := uuid.Parse(strings.TrimSpace(input.DeviceID))
		if err != nil {
			return UserDTO{}, TokenPair{}, fmt.Errorf(
				"%w: invalid device id",
				core_errors.ErrInvalidArgument,
			)
		}

		deviceID = parsedDeviceID
	}

	var tokens TokenPair

	err = s.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		tokenPair, err := s.issueTokenPair(ctx, user, deviceID)
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
