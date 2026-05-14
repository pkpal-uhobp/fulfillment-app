package orders_service

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (s *OrdersService) CancelOrder(
	ctx context.Context,
	orderID int64,
	actorID int64,
	actorRole string,
	input CancelOrderInput,
) error {
	if orderID <= 0 {
		return fmt.Errorf("%w: invalid order id", core_errors.ErrInvalidArgument)
	}

	order, err := s.repo.GetOrder(ctx, orderID)
	if err != nil {
		return fmt.Errorf("get order before cancel: %w", err)
	}

	if actorRole == core_domain.RoleClient.String() && order.Order.ClientID != actorID {
		return fmt.Errorf("%w: access denied", core_errors.ErrForbidden)
	}

	if isTerminalStatus(order.Order.Status) {
		return fmt.Errorf("%w: order cannot be cancelled from current status", core_errors.ErrInvalidArgument)
	}

	if err := validateOrderStatusTransition(order.Order.HandoverType, order.Order.Status, core_domain.OrderStatusCancelled); err != nil {
		return err
	}

	if err := s.repo.CancelOrder(ctx, orderID, actorID, input.Comment); err != nil {
		return fmt.Errorf("cancel order: %w", err)
	}

	return nil
}
