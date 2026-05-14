BEGIN;

-- =========================================================
-- 000007 ORDER CARGO PLACES
-- Заявленный клиентом состав груза:
-- например, 2 паллеты, 5 коробок, 3 мешка.
-- =========================================================

CREATE TABLE order_cargo_places (
    id BIGSERIAL PRIMARY KEY,

    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    cargo_place_type_id BIGINT NOT NULL REFERENCES cargo_place_types(id) ON DELETE RESTRICT,

    quantity INT NOT NULL,

    weight_per_place_kg NUMERIC(10, 2),
    length_cm NUMERIC(10, 2),
    width_cm NUMERIC(10, 2),
    height_cm NUMERIC(10, 2),

    comment TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT chk_order_cargo_places_quantity
        CHECK (quantity > 0),

    CONSTRAINT chk_order_cargo_places_weight
        CHECK (weight_per_place_kg IS NULL OR weight_per_place_kg > 0),

    CONSTRAINT chk_order_cargo_places_length
        CHECK (length_cm IS NULL OR length_cm > 0),

    CONSTRAINT chk_order_cargo_places_width
        CHECK (width_cm IS NULL OR width_cm > 0),

    CONSTRAINT chk_order_cargo_places_height
        CHECK (height_cm IS NULL OR height_cm > 0),

    CONSTRAINT uq_order_cargo_places_id_order_id
        UNIQUE (id, order_id),

    CONSTRAINT uq_order_cargo_places_id_type_id
        UNIQUE (id, cargo_place_type_id)
);

COMMIT;
