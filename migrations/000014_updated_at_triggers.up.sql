BEGIN;

-- =========================================================
-- 000014 UPDATED_AT TRIGGERS
-- =========================================================

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_users_set_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_orders_set_updated_at
BEFORE UPDATE ON orders
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_pickup_requests_set_updated_at
BEFORE UPDATE ON pickup_requests
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_cargo_items_set_updated_at
BEFORE UPDATE ON cargo_items
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

COMMIT;
