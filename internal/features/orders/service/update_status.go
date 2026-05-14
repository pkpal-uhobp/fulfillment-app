package orders_service

import (
	"context"
	"fmt"

	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (s *OrdersService) UpdateOrderStatus(
	ctx context.Context,
	orderID int64,
	actorID int64,
	input UpdateOrderStatusInput,
) (OrderDTO, error) {
	if orderID <= 0 {
		return OrderDTO{}, fmt.Errorf("%w: invalid order id", core_errors.ErrInvalidArgument)
	}

	newStatus, err := validateOrderStatus(input.Status)
	if err != nil {
		return OrderDTO{}, err
	}

	order, err := s.repo.GetOrder(ctx, orderID)
	if err != nil {
		return OrderDTO{}, fmt.Errorf("get order before status update: %w", err)
	}

	if err := validateOrderStatusTransition(order.Order.HandoverType, order.Order.Status, newStatus); err != nil {
		return OrderDTO{}, err
	}

	if err := s.repo.UpdateOrderStatus(ctx, orderID, newStatus.String(), actorID, input.Comment); err != nil {
		return OrderDTO{}, fmt.Errorf("update order status: %w", err)
	}

	updatedOrder, err := s.repo.GetOrder(ctx, orderID)
	if err != nil {
		return OrderDTO{}, fmt.Errorf("get updated order: %w", err)
	}

	return mapOrderDetailsToDTO(updatedOrder), nil
}
