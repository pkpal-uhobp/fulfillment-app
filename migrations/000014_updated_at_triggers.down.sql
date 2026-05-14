BEGIN;

DROP TRIGGER IF EXISTS trg_cargo_items_set_updated_at ON cargo_items;
DROP TRIGGER IF EXISTS trg_pickup_requests_set_updated_at ON pickup_requests;
DROP TRIGGER IF EXISTS trg_orders_set_updated_at ON orders;
DROP TRIGGER IF EXISTS trg_users_set_updated_at ON users;

DROP FUNCTION IF EXISTS set_updated_at();

COMMIT;
