package users_service

import (
	"fmt"
	"strings"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

func requireAdmin(role string) error {
	if core_domain.Role(role) != core_domain.RoleAdmin {
		return fmt.Errorf("%w: admin role required", core_errors.ErrForbidden)
	}
	return nil
}

func normalizePageLimit(page, limit int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	return page, limit
}

func validateRole(value string) (core_domain.Role, error) {
	role := core_domain.Role(strings.TrimSpace(value))
	if !role.IsValid() {
		return "", fmt.Errorf("%w: invalid user role", core_errors.ErrInvalidArgument)
	}
	return role, nil
}
