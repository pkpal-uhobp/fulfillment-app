package pickupcalendar_repository_postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

const dateLayout = "2006-01-02"

func scanBlockRow(row pgx.Row) (core_domain.PickupCalendarBlock, error) {
	var (
		block       core_domain.PickupCalendarBlock
		blockedDate time.Time
		reason      sql.NullString
	)

	err := row.Scan(
		&block.ID,
		&block.WarehouseID,
		&blockedDate,
		&reason,
		&block.CreatedBy,
		&block.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core_domain.PickupCalendarBlock{}, fmt.Errorf("%w: pickup calendar block not found", core_errors.ErrNotFound)
		}
		return core_domain.PickupCalendarBlock{}, fmt.Errorf("scan pickup calendar block: %w", err)
	}
	block.BlockedDate = blockedDate.Format(dateLayout)
	if reason.Valid {
		block.Reason = &reason.String
	}
	return block, nil
}

func scanCapacityRow(row pgx.Row) (core_domain.PickupCapacity, error) {
	var (
		capacity   core_domain.PickupCapacity
		pickupDate time.Time
	)

	err := row.Scan(
		&capacity.ID,
		&capacity.WarehouseID,
		&pickupDate,
		&capacity.MaxOrders,
		&capacity.CurrentOrders,
		&capacity.IsClosed,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core_domain.PickupCapacity{}, fmt.Errorf("%w: pickup capacity not found", core_errors.ErrNotFound)
		}
		return core_domain.PickupCapacity{}, fmt.Errorf("scan pickup capacity: %w", err)
	}
	capacity.PickupDate = pickupDate.Format(dateLayout)
	return capacity, nil
}
