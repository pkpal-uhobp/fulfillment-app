BEGIN;

-- =========================================================
-- 000004 WAREHOUSE STRUCTURE
-- worker_profiles, storage_zones, gates
-- =========================================================

CREATE TABLE worker_profiles (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    warehouse_id BIGINT NOT NULL REFERENCES warehouses(id) ON DELETE RESTRICT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE storage_zones (
    id BIGSERIAL PRIMARY KEY,

    warehouse_id BIGINT NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,

    name VARCHAR(100) NOT NULL,
    description TEXT,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    CONSTRAINT chk_storage_zones_name_not_empty
        CHECK (length(trim(name)) > 0),

    CONSTRAINT chk_storage_zones_name_format
        CHECK (
            name ~ '^[A-Za-zА-Яа-яЁё0-9 ._/-]+$'
        ),

    CONSTRAINT uq_storage_zones_warehouse_name
        UNIQUE (warehouse_id, name)
);

CREATE TABLE gates (
    id BIGSERIAL PRIMARY KEY,

    warehouse_id BIGINT NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,

    name VARCHAR(100) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    CONSTRAINT chk_gates_name_not_empty
        CHECK (length(trim(name)) > 0),

    CONSTRAINT chk_gates_name_format
        CHECK (
            name ~ '^[A-Za-zА-Яа-яЁё0-9 ._/-]+$'
        ),

    CONSTRAINT uq_gates_warehouse_name
        UNIQUE (warehouse_id, name)
);

COMMIT;
