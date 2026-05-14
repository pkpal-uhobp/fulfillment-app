package orders_service

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (s *OrdersService) GetOrder(
	ctx context.Context,
	orderID int64,
	actorID int64,
	actorRole string,
) (OrderDTO, error) {
	if orderID <= 0 {
		return OrderDTO{}, fmt.Errorf("%w: invalid order id", core_errors.ErrInvalidArgument)
	}

	order, err := s.repo.GetOrder(ctx, orderID)
	if err != nil {
		return OrderDTO{}, fmt.Errorf("get order: %w", err)
	}

	if actorRole == core_domain.RoleClient.String() && order.Order.ClientID != actorID {
		return OrderDTO{}, fmt.Errorf("%w: access denied", core_errors.ErrForbidden)
	}

	return mapOrderDetailsToDTO(order), nil
}
