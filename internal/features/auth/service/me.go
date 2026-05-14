package auth_service

import (
	"context"
	"fmt"
	"strings"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (s *AuthService) GetMe(
	ctx context.Context,
	userID int64,
) (UserDTO, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return UserDTO{}, err
	}

	if err := checkUserAvailable(user); err != nil {
		return UserDTO{}, err
	}

	return toUserDTO(user), nil
}

func checkUserAvailable(user core_domain.User) error {
	if !user.IsActive {
		return fmt.Errorf("%w: user is inactive", core_errors.ErrForbidden)
	}

	if user.IsBlocked {
		return fmt.Errorf("%w: user is blocked", core_errors.ErrForbidden)
	}

	return nil
}

func toUserDTO(user core_domain.User) UserDTO {
	dto := UserDTO{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.Role.String(),
	}

	if user.Phone != nil {
		dto.Phone = *user.Phone
	}

	return dto
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
