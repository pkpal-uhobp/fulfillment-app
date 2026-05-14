BEGIN;

-- =========================================================
-- 000010 CARGO RECEIPT DISCREPANCIES
-- Расхождения при приёмке:
-- заявлено 5 коробок, фактически принято 4.
-- =========================================================

CREATE TABLE cargo_receipt_discrepancies (
    id BIGSERIAL PRIMARY KEY,

    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,

    order_cargo_place_id BIGINT NOT NULL,
    cargo_place_type_id BIGINT NOT NULL REFERENCES cargo_place_types(id) ON DELETE RESTRICT,

    declared_quantity INT NOT NULL DEFAULT 0,
    actual_quantity INT NOT NULL DEFAULT 0,

    comment TEXT,
    created_by BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_discrepancies_declared_place_order
        FOREIGN KEY (order_cargo_place_id, order_id)
            REFERENCES order_cargo_places(id, order_id)
            ON DELETE RESTRICT,

    CONSTRAINT fk_discrepancies_declared_place_type
        FOREIGN KEY (order_cargo_place_id, cargo_place_type_id)
            REFERENCES order_cargo_places(id, cargo_place_type_id)
            ON DELETE RESTRICT,

    CONSTRAINT chk_discrepancies_quantities
        CHECK (
            declared_quantity >= 0
                AND actual_quantity >= 0
                AND declared_quantity <> actual_quantity
        )
);

COMMIT;
