BEGIN;

-- =========================================================
-- 000009 CARGO ITEMS
-- Фактически принятые грузовые места.
-- Каждая запись = одно физическое место с уникальным QR-кодом.
-- =========================================================

CREATE TABLE cargo_items (
    id BIGSERIAL PRIMARY KEY,

    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,

    order_cargo_place_id BIGINT NOT NULL,
    cargo_place_type_id BIGINT NOT NULL REFERENCES cargo_place_types(id) ON DELETE RESTRICT,

    qr_code VARCHAR(255) NOT NULL UNIQUE,

    status VARCHAR(50) NOT NULL DEFAULT 'accepted',

    storage_zone_id BIGINT REFERENCES storage_zones(id) ON DELETE SET NULL,
    gate_id BIGINT REFERENCES gates(id) ON DELETE SET NULL,

    received_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    shipped_by BIGINT REFERENCES users(id) ON DELETE SET NULL,

    received_at TIMESTAMPTZ,
    shipped_at TIMESTAMPTZ,

    comment TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_cargo_items_declared_place_order
        FOREIGN KEY (order_cargo_place_id, order_id)
            REFERENCES order_cargo_places(id, order_id)
            ON DELETE RESTRICT,

    CONSTRAINT fk_cargo_items_declared_place_type
        FOREIGN KEY (order_cargo_place_id, cargo_place_type_id)
            REFERENCES order_cargo_places(id, cargo_place_type_id)
            ON DELETE RESTRICT,

    CONSTRAINT chk_cargo_items_status
        CHECK (
            status IN (
                'accepted',
                'stored',
                'ready_to_ship',
                'shipped',
                'lost',
                'damaged',
                'cancelled'
            )
        ),

    CONSTRAINT chk_cargo_items_qr_code_not_empty
        CHECK (length(trim(qr_code)) > 0),

    CONSTRAINT chk_cargo_items_qr_code_format
        CHECK (
            length(trim(qr_code)) BETWEEN 3 AND 255
                AND qr_code ~ '^[-A-Za-z0-9._:/]+$'
        ),

    CONSTRAINT chk_cargo_items_shipped_data
        CHECK (
            status <> 'shipped'
                OR
            shipped_at IS NOT NULL
        )
);

COMMIT;
