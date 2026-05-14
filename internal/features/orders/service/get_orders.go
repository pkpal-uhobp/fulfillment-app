package orders_service

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

const (
	defaultOrdersPage  = 1
	defaultOrdersLimit = 20
	maxOrdersLimit     = 100
)

func (s *OrdersService) ListOrders(
	ctx context.Context,
	actorID int64,
	actorRole string,
	filter OrderFilter,
) ([]OrderDTO, error) {
	clientID := filter.ClientID
	if actorRole == core_domain.RoleClient.String() {
		clientID = &actorID
	}

	if filter.Status != "" {
		if _, err := validateOrderStatus(filter.Status); err != nil {
			return nil, err
		}
	}

	if filter.HandoverType != "" {
		if _, err := validateHandoverType(filter.HandoverType); err != nil {
			return nil, err
		}
	}

	page := filter.Page
	if page <= 0 {
		page = defaultOrdersPage
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = defaultOrdersLimit
	}
	if limit > maxOrdersLimit {
		limit = maxOrdersLimit
	}

	orders, err := s.repo.ListOrders(
		ctx,
		core_domain.OrderFilter{
			ClientID:               clientID,
			Status:                 filter.Status,
			HandoverType:           filter.HandoverType,
			WarehouseID:            filter.WarehouseID,
			ReceivingWarehouseID:   filter.ReceivingWarehouseID,
			DestinationWarehouseID: filter.DestinationWarehouseID,
			Page:                   page,
			Limit:                  limit,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("list orders: %w", err)
	}

	return mapOrderDetailsListToDTO(orders), nil
}
