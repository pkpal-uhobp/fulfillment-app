BEGIN;

DROP INDEX IF EXISTS idx_cargo_status_history_changed_at;
DROP INDEX IF EXISTS idx_cargo_status_history_changed_by;
DROP INDEX IF EXISTS idx_cargo_status_history_cargo_item_id;

DROP INDEX IF EXISTS idx_order_status_history_changed_at;
DROP INDEX IF EXISTS idx_order_status_history_changed_by;
DROP INDEX IF EXISTS idx_order_status_history_order_id;

DROP INDEX IF EXISTS idx_shipment_items_shipment_id;

DROP INDEX IF EXISTS idx_shipments_created_by;
DROP INDEX IF EXISTS idx_shipments_planned_departure_at;
DROP INDEX IF EXISTS idx_shipments_status;
DROP INDEX IF EXISTS idx_shipments_gate_id;
DROP INDEX IF EXISTS idx_shipments_destination_warehouse_id;

DROP INDEX IF EXISTS idx_pickup_capacity_pickup_date;
DROP INDEX IF EXISTS idx_pickup_capacity_warehouse_id;

DROP INDEX IF EXISTS idx_pickup_calendar_blocks_blocked_date;
DROP INDEX IF EXISTS idx_pickup_calendar_blocks_warehouse_id;

DROP INDEX IF EXISTS idx_discrepancies_created_by;
DROP INDEX IF EXISTS idx_discrepancies_order_cargo_place_id;
DROP INDEX IF EXISTS idx_discrepancies_order_id;

DROP INDEX IF EXISTS idx_cargo_items_created_at;
DROP INDEX IF EXISTS idx_cargo_items_shipped_by;
DROP INDEX IF EXISTS idx_cargo_items_received_by;
DROP INDEX IF EXISTS idx_cargo_items_gate_id;
DROP INDEX IF EXISTS idx_cargo_items_storage_zone_id;
DROP INDEX IF EXISTS idx_cargo_items_status;
DROP INDEX IF EXISTS idx_cargo_items_type_id;
DROP INDEX IF EXISTS idx_cargo_items_order_cargo_place_id;
DROP INDEX IF EXISTS idx_cargo_items_order_id;

DROP INDEX IF EXISTS idx_pickup_requests_assigned_logist_id;
DROP INDEX IF EXISTS idx_pickup_requests_pickup_date;
DROP INDEX IF EXISTS idx_pickup_requests_status;
DROP INDEX IF EXISTS idx_pickup_requests_order_id;

DROP INDEX IF EXISTS idx_order_cargo_places_type_id;
DROP INDEX IF EXISTS idx_order_cargo_places_order_id;

DROP INDEX IF EXISTS idx_orders_created_at;
DROP INDEX IF EXISTS idx_orders_product_type_id;
DROP INDEX IF EXISTS idx_orders_destination_warehouse_id;
DROP INDEX IF EXISTS idx_orders_receiving_warehouse_id;
DROP INDEX IF EXISTS idx_orders_self_delivery_date;
DROP INDEX IF EXISTS idx_orders_handover_type;
DROP INDEX IF EXISTS idx_orders_status;
DROP INDEX IF EXISTS idx_orders_client_id;

DROP INDEX IF EXISTS idx_cargo_place_types_active;
DROP INDEX IF EXISTS idx_product_types_active;

DROP INDEX IF EXISTS idx_gates_active;
DROP INDEX IF EXISTS idx_gates_warehouse_id;

DROP INDEX IF EXISTS idx_storage_zones_active;
DROP INDEX IF EXISTS idx_storage_zones_warehouse_id;

DROP INDEX IF EXISTS idx_warehouses_active;
DROP INDEX IF EXISTS idx_warehouses_city;
DROP INDEX IF EXISTS idx_warehouses_type;

DROP INDEX IF EXISTS idx_worker_profiles_warehouse_id;

DROP INDEX IF EXISTS idx_issued_tokens_expires_at;
DROP INDEX IF EXISTS idx_issued_tokens_active;
DROP INDEX IF EXISTS idx_issued_tokens_user_device_type;
DROP INDEX IF EXISTS idx_issued_tokens_device_id;
DROP INDEX IF EXISTS idx_issued_tokens_user_id;

DROP INDEX IF EXISTS idx_users_blocked;
DROP INDEX IF EXISTS idx_users_active;
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS uq_users_email_lower;

COMMIT;
