package orders_service

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
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

	orders, err := s.repo.ListOrders(
		ctx,
		core_domain.OrderFilter{
			ClientID:     clientID,
			Status:      filter.Status,
			HandoverType: filter.HandoverType,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("list orders: %w", err)
	}

	return mapOrderDetailsListToDTO(orders), nil
}
