package orders_service

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (s *OrdersService) GetOrderHistory(
	ctx context.Context,
	orderID int64,
	actorID int64,
	actorRole string,
) ([]OrderStatusHistoryDTO, error) {
	if orderID <= 0 {
		return nil, fmt.Errorf("%w: invalid order id", core_errors.ErrInvalidArgument)
	}

	order, err := s.repo.GetOrder(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("get order before history: %w", err)
	}

	if actorRole == core_domain.RoleClient.String() && order.Order.ClientID != actorID {
		return nil, fmt.Errorf("%w: access denied", core_errors.ErrForbidden)
	}

	history, err := s.repo.ListOrderStatusHistory(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("get order status history: %w", err)
	}

	return mapOrderStatusHistoryListToDTO(history), nil
}
