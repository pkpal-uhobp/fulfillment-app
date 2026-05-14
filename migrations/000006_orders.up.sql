BEGIN;

-- =========================================================
-- 000006 ORDERS
-- Заявки клиентов.
-- Для pickup дата хранится в pickup_requests.
-- Для self_delivery дата самостоятельной сдачи хранится прямо в orders.
-- =========================================================

CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,

    client_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,

    receiving_warehouse_id BIGINT NOT NULL REFERENCES warehouses(id) ON DELETE RESTRICT,
    destination_warehouse_id BIGINT NOT NULL REFERENCES warehouses(id) ON DELETE RESTRICT,

    product_type_id BIGINT NOT NULL REFERENCES product_types(id) ON DELETE RESTRICT,

    handover_type VARCHAR(50) NOT NULL,

    self_delivery_date DATE,
    self_delivery_time_from TIME,
    self_delivery_time_to TIME,

    status VARCHAR(50) NOT NULL DEFAULT 'created',

    comment TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT chk_orders_handover_type
        CHECK (handover_type IN ('pickup', 'self_delivery')),

    CONSTRAINT chk_orders_status
        CHECK (
            status IN (
                'created',
                'waiting_pickup',
                'waiting_delivery',
                'received',
                'stored',
                'assigned_to_shipping',
                'shipped',
                'delivered',
                'cancelled'
            )
        ),

    CONSTRAINT chk_orders_status_matches_handover
        CHECK (
            (handover_type = 'pickup' AND status <> 'waiting_delivery')
                OR
            (handover_type = 'self_delivery' AND status <> 'waiting_pickup')
        ),

    CONSTRAINT chk_orders_self_delivery_schedule
        CHECK (
            (
                handover_type = 'self_delivery'
                    AND self_delivery_date IS NOT NULL
            )
                OR
            (
                handover_type = 'pickup'
                    AND self_delivery_date IS NULL
                    AND self_delivery_time_from IS NULL
                    AND self_delivery_time_to IS NULL
            )
        ),

    CONSTRAINT chk_orders_self_delivery_time_range
        CHECK (
            (self_delivery_time_from IS NULL AND self_delivery_time_to IS NULL)
                OR
            (
                self_delivery_time_from IS NOT NULL
                    AND self_delivery_time_to IS NOT NULL
                    AND self_delivery_time_to > self_delivery_time_from
            )
        )
);

COMMIT;
