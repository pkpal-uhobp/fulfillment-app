BEGIN;

-- =========================================================
-- 000008 PICKUP REQUESTS
-- Заполняется только если orders.handover_type = 'pickup'.
-- Проверку соответствия handover_type лучше делать в backend
-- или отдельным триггером, так как CHECK не может ссылаться
-- на другую таблицу.
-- =========================================================

CREATE TABLE pickup_requests (
    id BIGSERIAL PRIMARY KEY,

    order_id BIGINT NOT NULL UNIQUE REFERENCES orders(id) ON DELETE CASCADE,

    pickup_address TEXT NOT NULL,
    pickup_date DATE NOT NULL,
    pickup_time_from TIME,
    pickup_time_to TIME,

    contact_name VARCHAR(255),
    contact_phone VARCHAR(50),

    status VARCHAR(50) NOT NULL DEFAULT 'created',

    assigned_logist_id BIGINT REFERENCES users(id) ON DELETE SET NULL,

    comment TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT chk_pickup_requests_status
        CHECK (
            status IN (
                'created',
                'approved',
                'assigned',
                'in_progress',
                'picked_up',
                'cancelled'
            )
        ),

    CONSTRAINT chk_pickup_requests_address_not_empty
        CHECK (length(trim(pickup_address)) > 0),

    CONSTRAINT chk_pickup_requests_address_format
        CHECK (
            length(trim(pickup_address)) BETWEEN 5 AND 500
                AND pickup_address !~ '[[:cntrl:]]'
        ),

    CONSTRAINT chk_pickup_requests_contact_name_format
        CHECK (
            contact_name IS NULL
                OR contact_name ~ '^[A-Za-zА-Яа-яЁё]+([ .''-][A-Za-zА-Яа-яЁё]+)*$'
        ),

    CONSTRAINT chk_pickup_requests_contact_phone_format
        CHECK (
            contact_phone IS NULL
                OR contact_phone ~ '^[+]?[0-9]{10,15}$'
        ),

    CONSTRAINT chk_pickup_requests_time_range
        CHECK (
            (pickup_time_from IS NULL AND pickup_time_to IS NULL)
                OR
            (
                pickup_time_from IS NOT NULL
                    AND pickup_time_to IS NOT NULL
                    AND pickup_time_to > pickup_time_from
            )
        )
);

COMMIT;
