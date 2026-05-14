BEGIN;

-- =========================================================
-- 000012 SHIPMENTS
-- Отгрузки на склад назначения и связь с QR-грузовыми местами.
-- =========================================================

CREATE TABLE shipments (
    id BIGSERIAL PRIMARY KEY,

    destination_warehouse_id BIGINT NOT NULL REFERENCES warehouses(id) ON DELETE RESTRICT,
    gate_id BIGINT NOT NULL REFERENCES gates(id) ON DELETE RESTRICT,

    planned_departure_at TIMESTAMPTZ NOT NULL,
    actual_departure_at TIMESTAMPTZ,

    status VARCHAR(50) NOT NULL DEFAULT 'planned',

    created_by BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT chk_shipments_status
        CHECK (
            status IN (
                'planned',
                'loading',
                'shipped',
                'completed',
                'cancelled'
            )
        ),

    CONSTRAINT chk_shipments_actual_departure
        CHECK (
            actual_departure_at IS NULL
                OR actual_departure_at >= planned_departure_at
        )
);

CREATE TABLE shipment_items (
    id BIGSERIAL PRIMARY KEY,

    shipment_id BIGINT NOT NULL REFERENCES shipments(id) ON DELETE CASCADE,
    cargo_item_id BIGINT NOT NULL REFERENCES cargo_items(id) ON DELETE RESTRICT,

    CONSTRAINT uq_shipment_items_shipment_cargo
        UNIQUE (shipment_id, cargo_item_id),

    CONSTRAINT uq_shipment_items_cargo_item
        UNIQUE (cargo_item_id)
);

COMMIT;
