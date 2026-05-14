BEGIN;

-- =========================================================
-- 000005 CATALOGS
-- product_types, cargo_place_types
-- =========================================================

CREATE TABLE product_types (
    id BIGSERIAL PRIMARY KEY,

    name VARCHAR(150) NOT NULL UNIQUE,
    description TEXT,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    CONSTRAINT chk_product_types_name_not_empty
        CHECK (length(trim(name)) > 0),

    CONSTRAINT chk_product_types_name_format
        CHECK (
            name ~ '^[A-Za-zА-Яа-яЁё0-9 ._/-]+$'
        )
);

CREATE TABLE cargo_place_types (
    id BIGSERIAL PRIMARY KEY,

    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    CONSTRAINT chk_cargo_place_types_name_not_empty
        CHECK (length(trim(name)) > 0),

    CONSTRAINT chk_cargo_place_types_name_format
        CHECK (
            name ~ '^[A-Za-zА-Яа-яЁё0-9 ._/-]+$'
        )
);

COMMIT;
