package orders_service

import (
	"fmt"
	"time"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

const (
	dateLayout = "2006-01-02"
	timeLayout = "15:04"
)

func validateHandoverType(value string) (core_domain.HandoverType, error) {
	switch value {
	case core_domain.HandoverTypePickup.String():
		return core_domain.HandoverTypePickup, nil
	case core_domain.HandoverTypeSelfDelivery.String():
		return core_domain.HandoverTypeSelfDelivery, nil
	default:
		return "", fmt.Errorf("%w: invalid handover type", core_errors.ErrInvalidArgument)
	}
}

func validateOrderStatus(value string) (core_domain.OrderStatus, error) {
	switch value {
	case core_domain.OrderStatusCreated.String(),
		core_domain.OrderStatusWaitingPickup.String(),
		core_domain.OrderStatusWaitingDelivery.String(),
		core_domain.OrderStatusReceived.String(),
		core_domain.OrderStatusStored.String(),
		core_domain.OrderStatusAssignedToShipping.String(),
		core_domain.OrderStatusShipped.String(),
		core_domain.OrderStatusDelivered.String(),
		core_domain.OrderStatusCancelled.String():
		return core_domain.OrderStatus(value), nil
	default:
		return "", fmt.Errorf("%w: invalid order status", core_errors.ErrInvalidArgument)
	}
}

func validateDate(value string, field string) error {
	if _, err := time.Parse(dateLayout, value); err != nil {
		return fmt.Errorf("%w: invalid %s format, expected YYYY-MM-DD", core_errors.ErrInvalidArgument, field)
	}
	return nil
}

func validateTime(value string, field string) error {
	if _, err := time.Parse(timeLayout, value); err != nil {
		return fmt.Errorf("%w: invalid %s format, expected HH:MM", core_errors.ErrInvalidArgument, field)
	}
	return nil
}

func validateTimeRange(from *string, to *string, field string) error {
	if from == nil && to == nil {
		return nil
	}
	if from == nil || to == nil {
		return fmt.Errorf("%w: %s time range must contain both from and to", core_errors.ErrInvalidArgument, field)
	}

	parsedFrom, err := time.Parse(timeLayout, *from)
	if err != nil {
		return fmt.Errorf("%w: invalid %s time from format, expected HH:MM", core_errors.ErrInvalidArgument, field)
	}

	parsedTo, err := time.Parse(timeLayout, *to)
	if err != nil {
		return fmt.Errorf("%w: invalid %s time to format, expected HH:MM", core_errors.ErrInvalidArgument, field)
	}

	if !parsedTo.After(parsedFrom) {
		return fmt.Errorf("%w: %s time to must be after time from", core_errors.ErrInvalidArgument, field)
	}

	return nil
}

func validateStatusMatchesHandoverType(
	handoverType core_domain.HandoverType,
	status core_domain.OrderStatus,
) error {
	if handoverType == core_domain.HandoverTypePickup && status == core_domain.OrderStatusWaitingDelivery {
		return fmt.Errorf("%w: pickup order cannot have waiting_delivery status", core_errors.ErrInvalidArgument)
	}
	if handoverType == core_domain.HandoverTypeSelfDelivery && status == core_domain.OrderStatusWaitingPickup {
		return fmt.Errorf("%w: self_delivery order cannot have waiting_pickup status", core_errors.ErrInvalidArgument)
	}
	return nil
}

func isTerminalStatus(status core_domain.OrderStatus) bool {
	return status == core_domain.OrderStatusCancelled || status == core_domain.OrderStatusDelivered
}

func validateOrderStatusTransition(
	handoverType core_domain.HandoverType,
	currentStatus core_domain.OrderStatus,
	newStatus core_domain.OrderStatus,
) error {
	if err := validateStatusMatchesHandoverType(handoverType, newStatus); err != nil {
		return err
	}

	if currentStatus == newStatus {
		return fmt.Errorf("%w: order already has status %s", core_errors.ErrInvalidArgument, newStatus)
	}

	allowedTransitions := map[core_domain.OrderStatus][]core_domain.OrderStatus{
		core_domain.OrderStatusCreated: []core_domain.OrderStatus{
			core_domain.OrderStatusWaitingPickup,
			core_domain.OrderStatusWaitingDelivery,
			core_domain.OrderStatusCancelled,
		},
		core_domain.OrderStatusWaitingPickup: []core_domain.OrderStatus{
			core_domain.OrderStatusReceived,
			core_domain.OrderStatusCancelled,
		},
		core_domain.OrderStatusWaitingDelivery: []core_domain.OrderStatus{
			core_domain.OrderStatusReceived,
			core_domain.OrderStatusCancelled,
		},
		core_domain.OrderStatusReceived: []core_domain.OrderStatus{
			core_domain.OrderStatusStored,
			core_domain.OrderStatusCancelled,
		},
		core_domain.OrderStatusStored: []core_domain.OrderStatus{
			core_domain.OrderStatusAssignedToShipping,
			core_domain.OrderStatusCancelled,
		},
		core_domain.OrderStatusAssignedToShipping: []core_domain.OrderStatus{
			core_domain.OrderStatusShipped,
			core_domain.OrderStatusCancelled,
		},
		core_domain.OrderStatusShipped: []core_domain.OrderStatus{
			core_domain.OrderStatusDelivered,
		},
		core_domain.OrderStatusDelivered: []core_domain.OrderStatus{},
		core_domain.OrderStatusCancelled: []core_domain.OrderStatus{},
	}

	for _, allowedStatus := range allowedTransitions[currentStatus] {
		if allowedStatus == newStatus {
			return nil
		}
	}

	return fmt.Errorf(
		"%w: invalid order status transition from %s to %s",
		core_errors.ErrInvalidArgument,
		currentStatus,
		newStatus,
	)
}
