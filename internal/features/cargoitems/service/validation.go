package cargoitems_service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

const (
	defaultServiceCargoItemsPage  = 1
	defaultServiceCargoItemsLimit = 20
	maxServiceCargoItemsLimit     = 100
)

var qrCodeRegexp = regexp.MustCompile(`^[-A-Za-z0-9._:/]+$`)

func normalizePageLimit(page int, limit int) (int, int) {
	if page <= 0 {
		page = defaultServiceCargoItemsPage
	}
	if limit <= 0 {
		limit = defaultServiceCargoItemsLimit
	}
	if limit > maxServiceCargoItemsLimit {
		limit = maxServiceCargoItemsLimit
	}
	return page, limit
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

func normalizeQRCode(value *string) (string, error) {
	if value == nil || strings.TrimSpace(*value) == "" {
		return generateQRCode()
	}

	qrCode := strings.TrimSpace(*value)
	if len(qrCode) < 3 || len(qrCode) > 255 || !qrCodeRegexp.MatchString(qrCode) {
		return "", fmt.Errorf("%w: invalid qr code", core_errors.ErrInvalidArgument)
	}

	return qrCode, nil
}

func generateQRCode() (string, error) {
	buffer := make([]byte, 8)
	if _, err := rand.Read(buffer); err != nil {
		return "", fmt.Errorf("generate qr code: %w", err)
	}
	return "CI-" + strings.ToUpper(hex.EncodeToString(buffer)), nil
}

func validateCargoItemStatus(status string) (core_domain.CargoItemStatus, error) {
	value := core_domain.CargoItemStatus(strings.TrimSpace(status))
	if !value.IsValid() {
		return "", fmt.Errorf("%w: invalid cargo item status", core_errors.ErrInvalidArgument)
	}
	return value, nil
}

func validateCargoStatusTransition(
	current core_domain.CargoItem,
	next core_domain.CargoItemStatus,
) error {
	if current.Status == next {
		return nil
	}

	if current.Status.IsTerminal() {
		return fmt.Errorf("%w: terminal cargo item status cannot be changed", core_errors.ErrConflict)
	}

	allowed := map[core_domain.CargoItemStatus][]core_domain.CargoItemStatus{
		core_domain.CargoItemStatusAccepted: {
			core_domain.CargoItemStatusStored,
			core_domain.CargoItemStatusLost,
			core_domain.CargoItemStatusDamaged,
			core_domain.CargoItemStatusCancelled,
		},
		core_domain.CargoItemStatusStored: {
			core_domain.CargoItemStatusReadyToShip,
			core_domain.CargoItemStatusLost,
			core_domain.CargoItemStatusDamaged,
			core_domain.CargoItemStatusCancelled,
		},
		core_domain.CargoItemStatusReadyToShip: {
			core_domain.CargoItemStatusShipped,
			core_domain.CargoItemStatusLost,
			core_domain.CargoItemStatusDamaged,
			core_domain.CargoItemStatusCancelled,
		},
	}

	for _, candidate := range allowed[current.Status] {
		if candidate == next {
			return validateCargoStatusRequirements(current, next)
		}
	}

	return fmt.Errorf(
		"%w: invalid cargo item status transition %s -> %s",
		core_errors.ErrConflict,
		current.Status.String(),
		next.String(),
	)
}

func validateCargoStatusRequirements(
	current core_domain.CargoItem,
	next core_domain.CargoItemStatus,
) error {
	switch next {
	case core_domain.CargoItemStatusStored:
		if current.StorageZoneID == nil {
			return fmt.Errorf("%w: storage zone must be assigned before stored status", core_errors.ErrConflict)
		}
	case core_domain.CargoItemStatusReadyToShip, core_domain.CargoItemStatusShipped:
		if current.GateID == nil {
			return fmt.Errorf("%w: gate must be assigned before shipping status", core_errors.ErrConflict)
		}
	}

	return nil
}

func ensureCargoItemEditable(item core_domain.CargoItem) error {
	if item.Status.IsTerminal() {
		return fmt.Errorf("%w: terminal cargo item cannot be changed", core_errors.ErrConflict)
	}
	return nil
}
