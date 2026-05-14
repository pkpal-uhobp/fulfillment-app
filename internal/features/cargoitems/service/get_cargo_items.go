package cargoitems_service

import (
	"context"
	"strings"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

func (s *CargoItemsService) ListCargoItems(
	ctx context.Context,
	actorID int64,
	actorRole string,
	filter CargoItemFilter,
) ([]CargoItemDTO, error) {
	page, limit := normalizePageLimit(filter.Page, filter.Limit)

	status := strings.TrimSpace(filter.Status)
	if status != "" {
		if _, err := validateCargoItemStatus(status); err != nil {
			return nil, err
		}
	}

	domainFilter := core_domain.CargoItemFilter{
		OrderID:       filter.OrderID,
		Status:        status,
		StorageZoneID: filter.StorageZoneID,
		GateID:        filter.GateID,
		QRCode:        strings.TrimSpace(filter.QRCode),
		Page:          page,
		Limit:         limit,
	}

	if actorRole == core_domain.RoleClient.String() {
		domainFilter.ClientID = &actorID
	}

	items, err := s.repo.ListCargoItems(ctx, domainFilter)
	if err != nil {
		return nil, err
	}

	return mapCargoItemsToDTO(items), nil
}
