package users_service

import (
	"context"
	"strings"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

func (s *UsersService) ListUsers(
	ctx context.Context,
	actorRole string,
	filter UserFilter,
) ([]UserDTO, error) {
	if err := requireAdmin(actorRole); err != nil {
		return nil, err
	}

	filter.Role = strings.TrimSpace(filter.Role)
	filter.Search = strings.TrimSpace(filter.Search)
	if filter.Role != "" {
		if _, err := validateRole(filter.Role); err != nil {
			return nil, err
		}
	}
	page, limit := normalizePageLimit(filter.Page, filter.Limit)

	users, err := s.repo.ListUsers(ctx, core_domain.UserFilter{
		Role:      filter.Role,
		IsActive:  filter.IsActive,
		IsBlocked: filter.IsBlocked,
		Search:    filter.Search,
		Page:      page,
		Limit:     limit,
	})
	if err != nil {
		return nil, err
	}
	return toUserDTOs(users), nil
}
