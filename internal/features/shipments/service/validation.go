package shipments_service

import (
	"fmt"
	"strings"
	"time"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func requireLogistOrAdmin(actorRole string) error {
	if actorRole != core_domain.RoleLogist.String() && actorRole != core_domain.RoleAdmin.String() {
		return fmt.Errorf("%w: only logist or admin can manage shipments", core_errors.ErrForbidden)
	}
	return nil
}

func validateShipmentID(id int64, name string) error {
	if id <= 0 {
		return fmt.Errorf("%w: invalid %s", core_errors.ErrInvalidArgument, name)
	}
	return nil
}

func parsePlannedDeparture(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, fmt.Errorf("%w: planned departure is required", core_errors.ErrInvalidArgument)
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: planned departure must be RFC3339", core_errors.ErrInvalidArgument)
	}
	return parsed, nil
}

func validateShipmentStatus(value string) (core_domain.ShipmentStatus, error) {
	status := core_domain.ShipmentStatus(strings.TrimSpace(value))
	if !status.IsValid() {
		return "", fmt.Errorf("%w: invalid shipment status", core_errors.ErrInvalidArgument)
	}
	return status, nil
}

func validateShipmentStatusTransition(current core_domain.ShipmentStatus, next core_domain.ShipmentStatus) error {
	if current == next {
		return nil
	}

	switch current {
	case core_domain.ShipmentStatusPlanned:
		if next == core_domain.ShipmentStatusLoading || next == core_domain.ShipmentStatusCancelled {
			return nil
		}
	case core_domain.ShipmentStatusLoading:
		if next == core_domain.ShipmentStatusShipped || next == core_domain.ShipmentStatusCancelled {
			return nil
		}
	case core_domain.ShipmentStatusShipped:
		if next == core_domain.ShipmentStatusCompleted {
			return nil
		}
	}

	return fmt.Errorf("%w: invalid shipment status transition from %s to %s", core_errors.ErrConflict, current, next)
}
