BEGIN;

-- =========================================================
-- 000002 JWT ISSUED TOKENS
-- =========================================================
-- jti хранится в JWT и в БД.
-- По нему можно проверить:
-- был ли токен выдан системой и не был ли он отозван.

CREATE TABLE issued_tokens (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    jti UUID NOT NULL UNIQUE,

    token_type VARCHAR(20) NOT NULL,
    device_id UUID NOT NULL,

    revoked BOOLEAN NOT NULL DEFAULT FALSE,

    issued_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMPTZ NOT NULL,

    revoked_at TIMESTAMPTZ,
    revoked_reason TEXT,

    CONSTRAINT chk_issued_tokens_type
        CHECK (token_type IN ('access', 'refresh')),

    CONSTRAINT chk_issued_tokens_expires_after_issued
        CHECK (expires_at > issued_at),

    CONSTRAINT chk_issued_tokens_revoked_data
        CHECK (
            (
                revoked = FALSE
                    AND revoked_at IS NULL
                    AND revoked_reason IS NULL
            )
                OR
            (
                revoked = TRUE
                    AND revoked_at IS NOT NULL
            )
        )
);

COMMIT;
