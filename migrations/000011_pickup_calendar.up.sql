BEGIN;

-- =========================================================
-- 000011 PICKUP CALENDAR
-- Закрытые даты и лимиты заявок на дату.
-- =========================================================

CREATE TABLE pickup_calendar_blocks (
    id BIGSERIAL PRIMARY KEY,

    warehouse_id BIGINT NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,

    blocked_date DATE NOT NULL,
    reason TEXT,

    created_by BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_pickup_calendar_blocks_warehouse_date
        UNIQUE (warehouse_id, blocked_date)
);

CREATE TABLE pickup_capacity (
    id BIGSERIAL PRIMARY KEY,

    warehouse_id BIGINT NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,

    pickup_date DATE NOT NULL,

    max_orders INT NOT NULL DEFAULT 0,
    current_orders INT NOT NULL DEFAULT 0,
    is_closed BOOLEAN NOT NULL DEFAULT FALSE,

    CONSTRAINT uq_pickup_capacity_warehouse_date
        UNIQUE (warehouse_id, pickup_date),

    CONSTRAINT chk_pickup_capacity_max_orders
        CHECK (max_orders >= 0),

    CONSTRAINT chk_pickup_capacity_current_orders
        CHECK (current_orders >= 0),

    CONSTRAINT chk_pickup_capacity_current_not_greater_max
        CHECK (current_orders <= max_orders)
);

COMMIT;
