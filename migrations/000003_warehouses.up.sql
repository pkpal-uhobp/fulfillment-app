BEGIN;

-- =========================================================
-- 000003 WAREHOUSES
-- =========================================================

CREATE TABLE warehouses (
    id BIGSERIAL PRIMARY KEY,

    name VARCHAR(255) NOT NULL,

    warehouse_type VARCHAR(50) NOT NULL,
    marketplace VARCHAR(100),

    city VARCHAR(100) NOT NULL,
    address TEXT NOT NULL,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT chk_warehouses_type
        CHECK (warehouse_type IN ('receiving', 'destination', 'both')),

    CONSTRAINT chk_warehouses_name_not_empty
        CHECK (length(trim(name)) > 0),

    CONSTRAINT chk_warehouses_name_format
        CHECK (
            name ~ '^[A-Za-zА-Яа-яЁё0-9 .,''"№#()/_-]+$'
        ),

    CONSTRAINT chk_warehouses_marketplace_format
        CHECK (
            marketplace IS NULL
                OR marketplace ~ '^[A-Za-zА-Яа-яЁё0-9 ._/-]+$'
        ),

    CONSTRAINT chk_warehouses_city_not_empty
        CHECK (length(trim(city)) > 0),

    CONSTRAINT chk_warehouses_city_format
        CHECK (
            city ~ '^[A-Za-zА-Яа-яЁё]+([ -][A-Za-zА-Яа-яЁё]+)*$'
        ),

    CONSTRAINT chk_warehouses_address_not_empty
        CHECK (length(trim(address)) > 0),

    CONSTRAINT chk_warehouses_address_format
        CHECK (
            length(trim(address)) BETWEEN 5 AND 500
                AND address !~ '[[:cntrl:]]'
        )
);

COMMIT;
