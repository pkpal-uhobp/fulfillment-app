BEGIN;

-- =========================================================
-- DOWN MIGRATION
-- Удаляем таблицы в обратном порядке создания.
-- Индексы, CHECK, UNIQUE, FK и PRIMARY KEY удаляются вместе
-- с таблицами.
-- =========================================================

DROP TABLE IF EXISTS cargo_status_history;
DROP TABLE IF EXISTS order_status_history;

DROP TABLE IF EXISTS shipment_items;
DROP TABLE IF EXISTS shipments;

DROP TABLE IF EXISTS pickup_capacity;
DROP TABLE IF EXISTS pickup_calendar_blocks;

DROP TABLE IF EXISTS cargo_receipt_discrepancies;
DROP TABLE IF EXISTS cargo_items;

DROP TABLE IF EXISTS pickup_requests;
DROP TABLE IF EXISTS order_cargo_places;
DROP TABLE IF EXISTS orders;

DROP TABLE IF EXISTS cargo_place_types;
DROP TABLE IF EXISTS product_types;

DROP TABLE IF EXISTS gates;
DROP TABLE IF EXISTS storage_zones;
DROP TABLE IF EXISTS worker_profiles;
DROP TABLE IF EXISTS warehouses;

DROP TABLE IF EXISTS issued_tokens;
DROP TABLE IF EXISTS users;

COMMIT;