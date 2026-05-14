BEGIN;

-- =========================================================
-- 000013 STATUS HISTORY
-- История статусов заявок и грузовых мест.
-- =========================================================

CREATE TABLE order_status_history (
    id BIGSERIAL PRIMARY KEY,

    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,

    old_status VARCHAR(50),
    new_status VARCHAR(50) NOT NULL,

    changed_by BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    comment TEXT,

    changed_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT chk_order_status_history_old_status
        CHECK (
            old_status IS NULL
                OR old_status IN (
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

    CONSTRAINT chk_order_status_history_new_status
        CHECK (
            new_status IN (
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
        )
);

CREATE TABLE cargo_status_history (
    id BIGSERIAL PRIMARY KEY,

    cargo_item_id BIGINT NOT NULL REFERENCES cargo_items(id) ON DELETE CASCADE,

    old_status VARCHAR(50),
    new_status VARCHAR(50) NOT NULL,

    changed_by BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    comment TEXT,

    changed_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT chk_cargo_status_history_old_status
        CHECK (
            old_status IS NULL
                OR old_status IN (
                    'accepted',
                    'stored',
                    'ready_to_ship',
                    'shipped',
                    'lost',
                    'damaged',
                    'cancelled'
                )
        ),

    CONSTRAINT chk_cargo_status_history_new_status
        CHECK (
            new_status IN (
                'accepted',
                'stored',
                'ready_to_ship',
                'shipped',
                'lost',
                'damaged',
                'cancelled'
            )
        )
);

COMMIT;
