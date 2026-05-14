BEGIN;

-- =========================================================
-- 000015 INDEXES
-- UNIQUE и PRIMARY KEY уже создают индексы автоматически.
-- Ниже — дополнительные индексы для частых JOIN, WHERE и фильтров.
-- =========================================================

-- users
CREATE UNIQUE INDEX uq_users_email_lower
    ON users (lower(email));

CREATE INDEX idx_users_role
    ON users(role);

CREATE INDEX idx_users_active
    ON users(is_active);

CREATE INDEX idx_users_blocked
    ON users(is_blocked);

-- issued_tokens
CREATE INDEX idx_issued_tokens_user_id
    ON issued_tokens(user_id);

CREATE INDEX idx_issued_tokens_device_id
    ON issued_tokens(device_id);

CREATE INDEX idx_issued_tokens_user_device_type
    ON issued_tokens(user_id, device_id, token_type);

CREATE INDEX idx_issued_tokens_active
    ON issued_tokens(user_id, token_type, device_id)
    WHERE revoked = FALSE;

CREATE INDEX idx_issued_tokens_expires_at
    ON issued_tokens(expires_at);

-- worker_profiles
CREATE INDEX idx_worker_profiles_warehouse_id
    ON worker_profiles(warehouse_id);

-- warehouses
CREATE INDEX idx_warehouses_type
    ON warehouses(warehouse_type);

CREATE INDEX idx_warehouses_city
    ON warehouses(city);

CREATE INDEX idx_warehouses_active
    ON warehouses(is_active);

-- storage_zones
CREATE INDEX idx_storage_zones_warehouse_id
    ON storage_zones(warehouse_id);

CREATE INDEX idx_storage_zones_active
    ON storage_zones(is_active);

-- gates
CREATE INDEX idx_gates_warehouse_id
    ON gates(warehouse_id);

CREATE INDEX idx_gates_active
    ON gates(is_active);

-- product_types
CREATE INDEX idx_product_types_active
    ON product_types(is_active);

-- cargo_place_types
CREATE INDEX idx_cargo_place_types_active
    ON cargo_place_types(is_active);

-- orders
CREATE INDEX idx_orders_client_id
    ON orders(client_id);

CREATE INDEX idx_orders_status
    ON orders(status);

CREATE INDEX idx_orders_handover_type
    ON orders(handover_type);

CREATE INDEX idx_orders_self_delivery_date
    ON orders(self_delivery_date);

CREATE INDEX idx_orders_receiving_warehouse_id
    ON orders(receiving_warehouse_id);

CREATE INDEX idx_orders_destination_warehouse_id
    ON orders(destination_warehouse_id);

CREATE INDEX idx_orders_product_type_id
    ON orders(product_type_id);

CREATE INDEX idx_orders_created_at
    ON orders(created_at);

-- order_cargo_places
CREATE INDEX idx_order_cargo_places_order_id
    ON order_cargo_places(order_id);

CREATE INDEX idx_order_cargo_places_type_id
    ON order_cargo_places(cargo_place_type_id);

-- pickup_requests
CREATE INDEX idx_pickup_requests_order_id
    ON pickup_requests(order_id);

CREATE INDEX idx_pickup_requests_status
    ON pickup_requests(status);

CREATE INDEX idx_pickup_requests_pickup_date
    ON pickup_requests(pickup_date);

CREATE INDEX idx_pickup_requests_assigned_logist_id
    ON pickup_requests(assigned_logist_id);

-- cargo_items
CREATE INDEX idx_cargo_items_order_id
    ON cargo_items(order_id);

CREATE INDEX idx_cargo_items_order_cargo_place_id
    ON cargo_items(order_cargo_place_id);

CREATE INDEX idx_cargo_items_type_id
    ON cargo_items(cargo_place_type_id);

CREATE INDEX idx_cargo_items_status
    ON cargo_items(status);

CREATE INDEX idx_cargo_items_storage_zone_id
    ON cargo_items(storage_zone_id);

CREATE INDEX idx_cargo_items_gate_id
    ON cargo_items(gate_id);

CREATE INDEX idx_cargo_items_received_by
    ON cargo_items(received_by);

CREATE INDEX idx_cargo_items_shipped_by
    ON cargo_items(shipped_by);

CREATE INDEX idx_cargo_items_created_at
    ON cargo_items(created_at);

-- discrepancies
CREATE INDEX idx_discrepancies_order_id
    ON cargo_receipt_discrepancies(order_id);

CREATE INDEX idx_discrepancies_order_cargo_place_id
    ON cargo_receipt_discrepancies(order_cargo_place_id);

CREATE INDEX idx_discrepancies_created_by
    ON cargo_receipt_discrepancies(created_by);

-- calendar
CREATE INDEX idx_pickup_calendar_blocks_warehouse_id
    ON pickup_calendar_blocks(warehouse_id);

CREATE INDEX idx_pickup_calendar_blocks_blocked_date
    ON pickup_calendar_blocks(blocked_date);

-- capacity
CREATE INDEX idx_pickup_capacity_warehouse_id
    ON pickup_capacity(warehouse_id);

CREATE INDEX idx_pickup_capacity_pickup_date
    ON pickup_capacity(pickup_date);

-- shipments
CREATE INDEX idx_shipments_destination_warehouse_id
    ON shipments(destination_warehouse_id);

CREATE INDEX idx_shipments_gate_id
    ON shipments(gate_id);

CREATE INDEX idx_shipments_status
    ON shipments(status);

CREATE INDEX idx_shipments_planned_departure_at
    ON shipments(planned_departure_at);

CREATE INDEX idx_shipments_created_by
    ON shipments(created_by);

-- shipment_items
CREATE INDEX idx_shipment_items_shipment_id
    ON shipment_items(shipment_id);

-- order_status_history
CREATE INDEX idx_order_status_history_order_id
    ON order_status_history(order_id);

CREATE INDEX idx_order_status_history_changed_by
    ON order_status_history(changed_by);

CREATE INDEX idx_order_status_history_changed_at
    ON order_status_history(changed_at);

-- cargo_status_history
CREATE INDEX idx_cargo_status_history_cargo_item_id
    ON cargo_status_history(cargo_item_id);

CREATE INDEX idx_cargo_status_history_changed_by
    ON cargo_status_history(changed_by);

CREATE INDEX idx_cargo_status_history_changed_at
    ON cargo_status_history(changed_at);

COMMIT;
