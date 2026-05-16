package pickupcalendar_service

import (
	"fmt"
	"strings"
	"time"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

const dateLayout = "2006-01-02"

func validateActorID(actorID int64) error {
	if actorID <= 0 {
		return fmt.Errorf("%w: invalid actor id", core_errors.ErrInvalidArgument)
	}
	return nil
}

func requireCalendarReadRole(actorRole string) error {
	switch actorRole {
	case core_domain.RoleClient.String(), core_domain.RoleLogist.String(), core_domain.RoleAdmin.String():
		return nil
	default:
		return fmt.Errorf("%w: only client, logist or admin can read pickup calendar", core_errors.ErrForbidden)
	}
}

func requireCalendarManageRole(actorRole string) error {
	switch actorRole {
	case core_domain.RoleLogist.String(), core_domain.RoleAdmin.String():
		return nil
	default:
		return fmt.Errorf("%w: only logist or admin can manage pickup calendar", core_errors.ErrForbidden)
	}
}

func validateWarehouseID(warehouseID int64) error {
	if warehouseID <= 0 {
		return fmt.Errorf("%w: invalid warehouse id", core_errors.ErrInvalidArgument)
	}
	return nil
}

func validateBlockID(blockID int64) error {
	if blockID <= 0 {
		return fmt.Errorf("%w: invalid calendar block id", core_errors.ErrInvalidArgument)
	}
	return nil
}

func validateDate(value string, field string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%w: %s is required", core_errors.ErrInvalidArgument, field)
	}
	if _, err := time.Parse(dateLayout, value); err != nil {
		return fmt.Errorf("%w: invalid %s format, expected YYYY-MM-DD", core_errors.ErrInvalidArgument, field)
	}
	return nil
}

func validateDateRange(dateFrom string, dateTo string) (string, string, error) {
	from := strings.TrimSpace(dateFrom)
	to := strings.TrimSpace(dateTo)

	if from == "" && to == "" {
		fromTime := time.Now().UTC()
		toTime := fromTime.AddDate(0, 0, 30)
		return fromTime.Format(dateLayout), toTime.Format(dateLayout), nil
	}

	if from == "" || to == "" {
		return "", "", fmt.Errorf("%w: date_from and date_to must be set together", core_errors.ErrInvalidArgument)
	}

	if err := validateDate(from, "date_from"); err != nil {
		return "", "", err
	}
	if err := validateDate(to, "date_to"); err != nil {
		return "", "", err
	}

	parsedFrom, _ := time.Parse(dateLayout, from)
	parsedTo, _ := time.Parse(dateLayout, to)
	if parsedTo.Before(parsedFrom) {
		return "", "", fmt.Errorf("%w: date_to must be after or equal date_from", core_errors.ErrInvalidArgument)
	}
	if parsedTo.Sub(parsedFrom).Hours()/24 > 92 {
		return "", "", fmt.Errorf("%w: pickup calendar range cannot be longer than 92 days", core_errors.ErrInvalidArgument)
	}
	return from, to, nil
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
