package pickupcalendar_repository_postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

func (r *PickupCalendarRepository) ListCalendar(
	ctx context.Context,
	filter core_domain.PickupCalendarFilter,
) ([]core_domain.PickupCalendarDay, error) {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()
	q := r.tx.Querier(ctx)

	const query = `
		WITH days AS (
			SELECT generate_series($2::date, $3::date, interval '1 day')::date AS day
		)
		SELECT
			$1::bigint AS warehouse_id,
			d.day,
			COALESCE(c.max_orders, 0) AS max_orders,
			COALESCE(c.current_orders, 0) AS current_orders,
			(
				COALESCE(c.is_closed, false)
				OR b.id IS NOT NULL
				OR (COALESCE(c.max_orders, 0) > 0 AND COALESCE(c.current_orders, 0) >= COALESCE(c.max_orders, 0))
			) AS is_closed,
			b.id,
			b.reason,
			b.created_by,
			b.created_at
		FROM days d
		LEFT JOIN pickup_capacity c
			ON c.warehouse_id = $1 AND c.pickup_date = d.day
		LEFT JOIN pickup_calendar_blocks b
			ON b.warehouse_id = $1 AND b.blocked_date = d.day
		ORDER BY d.day;
	`

	rows, err := q.Query(ctx, query, filter.WarehouseID, filter.DateFrom, filter.DateTo)
	if err != nil {
		return nil, fmt.Errorf("list pickup calendar: %w", err)
	}
	defer rows.Close()

	days := make([]core_domain.PickupCalendarDay, 0)
	for rows.Next() {
		var (
			day       core_domain.PickupCalendarDay
			date      time.Time
			blockID   sql.NullInt64
			reason    sql.NullString
			createdBy sql.NullInt64
			createdAt sql.NullTime
		)
		if err := rows.Scan(
			&day.WarehouseID,
			&date,
			&day.MaxOrders,
			&day.CurrentOrders,
			&day.IsClosed,
			&blockID,
			&reason,
			&createdBy,
			&createdAt,
		); err != nil {
			return nil, fmt.Errorf("scan pickup calendar day: %w", err)
		}
		day.Date = date.Format(dateLayout)
		if blockID.Valid {
			block := core_domain.PickupCalendarBlock{
				ID:          blockID.Int64,
				WarehouseID: day.WarehouseID,
				BlockedDate: day.Date,
			}
			if reason.Valid {
				block.Reason = &reason.String
			}
			if createdBy.Valid {
				block.CreatedBy = createdBy.Int64
			}
			if createdAt.Valid {
				block.CreatedAt = createdAt.Time
			}
			day.Block = &block
		}
		days = append(days, day)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate pickup calendar: %w", err)
	}
	return days, nil
}
