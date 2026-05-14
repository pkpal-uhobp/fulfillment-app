BEGIN;

-- =========================================================
-- 000001 USERS
-- =========================================================

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,

    email VARCHAR(255) NOT NULL,
    password_hash TEXT NOT NULL,

    full_name VARCHAR(255) NOT NULL,
    phone VARCHAR(50),

    role VARCHAR(50) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_blocked BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT chk_users_role
        CHECK (role IN ('client', 'logist', 'worker', 'admin')),

    CONSTRAINT chk_users_email_not_empty
        CHECK (length(trim(email)) > 0),

    CONSTRAINT chk_users_email_format
        CHECK (
            email ~* '^[A-Z0-9._%+-]+@[A-Z0-9.-]+[.][A-Z]{2,63}$'
        ),

    CONSTRAINT chk_users_password_hash_not_empty
        CHECK (length(trim(password_hash)) > 0),

    CONSTRAINT chk_users_full_name_not_empty
        CHECK (length(trim(full_name)) > 0),

    CONSTRAINT chk_users_full_name_format
        CHECK (
            full_name ~ '^[A-Za-zА-Яа-яЁё]+([ .''-][A-Za-zА-Яа-яЁё]+)*$'
        ),

    CONSTRAINT chk_users_phone_format
        CHECK (
            phone IS NULL
                OR phone ~ '^[+]?[0-9]{10,15}$'
        )
);

COMMIT;
