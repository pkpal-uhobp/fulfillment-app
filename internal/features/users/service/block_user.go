package users_service

import (
	"context"
	"fmt"

	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (s *UsersService) BlockUser(
	ctx context.Context,
	actorID int64,
	actorRole string,
	userID int64,
	input BlockUserInput,
) error {
	if err := requireAdmin(actorRole); err != nil {
		return err
	}
	if userID <= 0 {
		return fmt.Errorf("%w: invalid user id", core_errors.ErrInvalidArgument)
	}
	if userID == actorID {
		return fmt.Errorf("%w: cannot block yourself", core_errors.ErrForbidden)
	}
	return s.repo.BlockUser(ctx, userID, input.Reason)
}

func (s *UsersService) DeactivateUser(
	ctx context.Context,
	actorID int64,
	actorRole string,
	userID int64,
) error {
	if err := requireAdmin(actorRole); err != nil {
		return err
	}
	if userID <= 0 {
		return fmt.Errorf("%w: invalid user id", core_errors.ErrInvalidArgument)
	}
	if userID == actorID {
		return fmt.Errorf("%w: cannot deactivate yourself", core_errors.ErrForbidden)
	}
	return s.repo.DeactivateUser(ctx, userID)
}
