-- Test data for fulfillment-app.
-- Run after all migrations:
-- psql -h localhost -p 5433 -U postgres -d fulfillment-app -f scripts/seed_test_data.sql

BEGIN;

DO $$
<<seed>>
DECLARE
    admin_id BIGINT;
    logist_id BIGINT;
    worker_id BIGINT;
    client_id BIGINT;

    receiving_warehouse_id BIGINT;
    destination_warehouse_id BIGINT;
    storage_zone_id BIGINT;
    gate_id BIGINT;
    product_type_id BIGINT;
    cargo_place_type_id BIGINT;

    self_delivery_order_id BIGINT;
    pickup_order_id BIGINT;
    cargo_place_id BIGINT;
    pickup_cargo_place_id BIGINT;
    cargo_item_1_id BIGINT;
    cargo_item_2_id BIGINT;
    shipment_id BIGINT;

    -- Password for all test users: Password123!
    -- Bcrypt hash generated for Password123!.
    password_hash TEXT := '$2a$12$SHdXyODTjMTP7luGwEUDAeV9l4kYdiM1fJpeJ8qJQ5LTtvv4W4rz6';
BEGIN
    -- Users
    INSERT INTO users (email, password_hash, full_name, phone, role, is_active, is_blocked)
    SELECT 'admin@example.com', password_hash, 'Admin User', '+77000000001', 'admin', TRUE, FALSE
    WHERE NOT EXISTS (SELECT 1 FROM users WHERE lower(email) = lower('admin@example.com'));

    INSERT INTO users (email, password_hash, full_name, phone, role, is_active, is_blocked)
    SELECT 'logist@example.com', password_hash, 'Logist User', '+77000000002', 'logist', TRUE, FALSE
    WHERE NOT EXISTS (SELECT 1 FROM users WHERE lower(email) = lower('logist@example.com'));

    INSERT INTO users (email, password_hash, full_name, phone, role, is_active, is_blocked)
    SELECT 'worker@example.com', password_hash, 'Worker User', '+77000000003', 'worker', TRUE, FALSE
    WHERE NOT EXISTS (SELECT 1 FROM users WHERE lower(email) = lower('worker@example.com'));

    INSERT INTO users (email, password_hash, full_name, phone, role, is_active, is_blocked)
    SELECT 'client@example.com', password_hash, 'Client User', '+77000000004', 'client', TRUE, FALSE
    WHERE NOT EXISTS (SELECT 1 FROM users WHERE lower(email) = lower('client@example.com'));

    SELECT id INTO admin_id FROM users WHERE lower(email) = lower('admin@example.com');
    SELECT id INTO logist_id FROM users WHERE lower(email) = lower('logist@example.com');
    SELECT id INTO worker_id FROM users WHERE lower(email) = lower('worker@example.com');
    SELECT id INTO client_id FROM users WHERE lower(email) = lower('client@example.com');

    -- Catalogs, in case seed catalog migration was not applied.
    INSERT INTO product_types (name, description)
    SELECT 'Одежда', 'Одежда и текстильные изделия'
    WHERE NOT EXISTS (SELECT 1 FROM product_types WHERE name = 'Одежда');

    INSERT INTO cargo_place_types (name, description)
    SELECT 'Коробка', 'Товар в коробке'
    WHERE NOT EXISTS (SELECT 1 FROM cargo_place_types WHERE name = 'Коробка');

    SELECT id INTO product_type_id FROM product_types WHERE name = 'Одежда' LIMIT 1;
    SELECT id INTO cargo_place_type_id FROM cargo_place_types WHERE name = 'Коробка' LIMIT 1;

    -- Warehouses
    INSERT INTO warehouses (name, warehouse_type, marketplace, city, address, is_active)
    SELECT 'Приемный склад Алматы', 'receiving', NULL, 'Алматы', 'ул. Абая, 10', TRUE
    WHERE NOT EXISTS (SELECT 1 FROM warehouses WHERE name = 'Приемный склад Алматы');

    INSERT INTO warehouses (name, warehouse_type, marketplace, city, address, is_active)
    SELECT 'Ozon Астана', 'destination', 'Ozon', 'Астана', 'пр. Кабанбай Батыра, 20', TRUE
    WHERE NOT EXISTS (SELECT 1 FROM warehouses WHERE name = 'Ozon Астана');

    SELECT id INTO receiving_warehouse_id FROM warehouses WHERE name = 'Приемный склад Алматы';
    SELECT id INTO destination_warehouse_id FROM warehouses WHERE name = 'Ozon Астана';

    -- Warehouse structure
    INSERT INTO storage_zones (warehouse_id, name, description, is_active)
    SELECT receiving_warehouse_id, 'A-01', 'Тестовая зона хранения', TRUE
    WHERE NOT EXISTS (
        SELECT 1 FROM storage_zones WHERE warehouse_id = receiving_warehouse_id AND name = 'A-01'
    );

    INSERT INTO gates (warehouse_id, name, is_active)
    SELECT receiving_warehouse_id, 'Gate-1', TRUE
    WHERE NOT EXISTS (
        SELECT 1 FROM gates WHERE warehouse_id = receiving_warehouse_id AND name = 'Gate-1'
    );

    SELECT id INTO storage_zone_id FROM storage_zones WHERE warehouse_id = receiving_warehouse_id AND name = 'A-01';
    SELECT id INTO gate_id FROM gates WHERE warehouse_id = receiving_warehouse_id AND name = 'Gate-1';

    INSERT INTO worker_profiles (user_id, warehouse_id)
    SELECT worker_id, receiving_warehouse_id
    WHERE NOT EXISTS (SELECT 1 FROM worker_profiles WHERE user_id = worker_id);

    -- Pickup calendar test dates.
    INSERT INTO pickup_capacity (warehouse_id, pickup_date, max_orders, current_orders, is_closed)
    SELECT receiving_warehouse_id, CURRENT_DATE + INTERVAL '1 day', 10, 2, FALSE
    WHERE NOT EXISTS (
        SELECT 1 FROM pickup_capacity
        WHERE warehouse_id = receiving_warehouse_id
          AND pickup_date = CURRENT_DATE + INTERVAL '1 day'
    );

    INSERT INTO pickup_calendar_blocks (warehouse_id, blocked_date, reason, created_by)
    SELECT receiving_warehouse_id, CURRENT_DATE + INTERVAL '5 days', 'Тестовая закрытая дата', logist_id
    WHERE NOT EXISTS (
        SELECT 1 FROM pickup_calendar_blocks
        WHERE warehouse_id = receiving_warehouse_id
          AND blocked_date = CURRENT_DATE + INTERVAL '5 days'
    );

    -- Self delivery order.
    INSERT INTO orders (
        client_id,
        receiving_warehouse_id,
        destination_warehouse_id,
        product_type_id,
        handover_type,
        self_delivery_date,
        self_delivery_time_from,
        self_delivery_time_to,
        status,
        comment
    )
    SELECT
        client_id,
        receiving_warehouse_id,
        destination_warehouse_id,
        product_type_id,
        'self_delivery',
        CURRENT_DATE + INTERVAL '1 day',
        '10:00',
        '12:00',
        'received',
        'seed:self_delivery_order'
    WHERE NOT EXISTS (SELECT 1 FROM orders WHERE comment = 'seed:self_delivery_order');

    SELECT id INTO self_delivery_order_id FROM orders WHERE comment = 'seed:self_delivery_order';

    -- Pickup order.
    INSERT INTO orders (
        client_id,
        receiving_warehouse_id,
        destination_warehouse_id,
        product_type_id,
        handover_type,
        status,
        comment
    )
    SELECT
        client_id,
        receiving_warehouse_id,
        destination_warehouse_id,
        product_type_id,
        'pickup',
        'waiting_pickup',
        'seed:pickup_order'
    WHERE NOT EXISTS (SELECT 1 FROM orders WHERE comment = 'seed:pickup_order');

    SELECT id INTO pickup_order_id FROM orders WHERE comment = 'seed:pickup_order';

    INSERT INTO pickup_requests (
        order_id,
        pickup_address,
        pickup_date,
        pickup_time_from,
        pickup_time_to,
        contact_name,
        contact_phone,
        status,
        assigned_logist_id,
        comment
    )
    SELECT
        pickup_order_id,
        'ул. Толе би, 15',
        CURRENT_DATE + INTERVAL '2 days',
        '14:00',
        '16:00',
        'Client User',
        '+77000000004',
        'assigned',
        logist_id,
        'seed:pickup_request'
    WHERE NOT EXISTS (SELECT 1 FROM pickup_requests WHERE order_id = pickup_order_id);

    -- Declared cargo places.
    INSERT INTO order_cargo_places (
        order_id,
        cargo_place_type_id,
        quantity,
        weight_per_place_kg,
        length_cm,
        width_cm,
        height_cm,
        comment
    )
    SELECT
        self_delivery_order_id,
        cargo_place_type_id,
        2,
        5.50,
        40,
        30,
        25,
        'seed:self_delivery_cargo_place'
    WHERE NOT EXISTS (
        SELECT 1 FROM order_cargo_places
        WHERE order_id = self_delivery_order_id AND comment = 'seed:self_delivery_cargo_place'
    );

    INSERT INTO order_cargo_places (
        order_id,
        cargo_place_type_id,
        quantity,
        weight_per_place_kg,
        length_cm,
        width_cm,
        height_cm,
        comment
    )
    SELECT
        pickup_order_id,
        cargo_place_type_id,
        1,
        3.20,
        35,
        25,
        20,
        'seed:pickup_cargo_place'
    WHERE NOT EXISTS (
        SELECT 1 FROM order_cargo_places
        WHERE order_id = pickup_order_id AND comment = 'seed:pickup_cargo_place'
    );

    SELECT id INTO cargo_place_id
    FROM order_cargo_places
    WHERE order_id = self_delivery_order_id AND comment = 'seed:self_delivery_cargo_place'
    LIMIT 1;

    SELECT id INTO pickup_cargo_place_id
    FROM order_cargo_places
    WHERE order_id = pickup_order_id AND comment = 'seed:pickup_cargo_place'
    LIMIT 1;

    -- Cargo items for QR testing.
    INSERT INTO cargo_items (
        order_id,
        order_cargo_place_id,
        cargo_place_type_id,
        qr_code,
        status,
        storage_zone_id,
        gate_id,
        received_by,
        received_at,
        comment
    )
    SELECT
        self_delivery_order_id,
        cargo_place_id,
        cargo_place_type_id,
        'QR-TEST-001',
        'stored',
        storage_zone_id,
        NULL,
        worker_id,
        NOW(),
        'seed:stored_cargo_item'
    WHERE NOT EXISTS (SELECT 1 FROM cargo_items WHERE qr_code = 'QR-TEST-001');

    INSERT INTO cargo_items (
        order_id,
        order_cargo_place_id,
        cargo_place_type_id,
        qr_code,
        status,
        storage_zone_id,
        gate_id,
        received_by,
        received_at,
        comment
    )
    SELECT
        self_delivery_order_id,
        cargo_place_id,
        cargo_place_type_id,
        'QR-TEST-002',
        'ready_to_ship',
        storage_zone_id,
        gate_id,
        worker_id,
        NOW(),
        'seed:ready_to_ship_cargo_item'
    WHERE NOT EXISTS (SELECT 1 FROM cargo_items WHERE qr_code = 'QR-TEST-002');

    INSERT INTO cargo_items (
        order_id,
        order_cargo_place_id,
        cargo_place_type_id,
        qr_code,
        status,
        storage_zone_id,
        gate_id,
        received_by,
        received_at,
        comment
    )
    SELECT
        pickup_order_id,
        pickup_cargo_place_id,
        cargo_place_type_id,
        'QR-PICKUP-001',
        'accepted',
        NULL,
        NULL,
        worker_id,
        NOW(),
        'seed:accepted_pickup_cargo_item'
    WHERE NOT EXISTS (SELECT 1 FROM cargo_items WHERE qr_code = 'QR-PICKUP-001');

    SELECT id INTO cargo_item_1_id FROM cargo_items WHERE qr_code = 'QR-TEST-001';
    SELECT id INTO cargo_item_2_id FROM cargo_items WHERE qr_code = 'QR-TEST-002';

    -- Status history samples.
    INSERT INTO order_status_history (order_id, old_status, new_status, changed_by, comment)
    SELECT self_delivery_order_id, NULL, 'created', client_id, 'seed:order_created'
    WHERE NOT EXISTS (
        SELECT 1 FROM order_status_history
        WHERE order_id = self_delivery_order_id AND comment = 'seed:order_created'
    );

    INSERT INTO order_status_history (order_id, old_status, new_status, changed_by, comment)
    SELECT self_delivery_order_id, 'created', 'received', worker_id, 'seed:order_received'
    WHERE NOT EXISTS (
        SELECT 1 FROM order_status_history
        WHERE order_id = self_delivery_order_id AND comment = 'seed:order_received'
    );

    INSERT INTO cargo_status_history (cargo_item_id, old_status, new_status, changed_by, comment)
    SELECT cargo_item_1_id, NULL, 'accepted', worker_id, 'seed:cargo_accepted'
    WHERE NOT EXISTS (
        SELECT 1 FROM cargo_status_history
        WHERE cargo_item_id = cargo_item_1_id AND comment = 'seed:cargo_accepted'
    );

    INSERT INTO cargo_status_history (cargo_item_id, old_status, new_status, changed_by, comment)
    SELECT cargo_item_1_id, 'accepted', 'stored', worker_id, 'seed:cargo_stored'
    WHERE NOT EXISTS (
        SELECT 1 FROM cargo_status_history
        WHERE cargo_item_id = cargo_item_1_id AND comment = 'seed:cargo_stored'
    );

    INSERT INTO cargo_status_history (cargo_item_id, old_status, new_status, changed_by, comment)
    SELECT cargo_item_2_id, 'stored', 'ready_to_ship', logist_id, 'seed:cargo_ready_to_ship'
    WHERE NOT EXISTS (
        SELECT 1 FROM cargo_status_history
        WHERE cargo_item_id = cargo_item_2_id AND comment = 'seed:cargo_ready_to_ship'
    );

    -- Shipment with one ready-to-ship cargo item.
    INSERT INTO shipments (
        destination_warehouse_id,
        gate_id,
        planned_departure_at,
        status,
        created_by
    )
    SELECT
        destination_warehouse_id,
        gate_id,
        NOW() + INTERVAL '2 days',
        'planned',
        logist_id
    WHERE NOT EXISTS (
        SELECT 1 FROM shipments s
        WHERE s.destination_warehouse_id = seed.destination_warehouse_id
          AND s.gate_id = seed.gate_id
          AND s.created_by = seed.logist_id
          AND s.status = 'planned'
    );

    SELECT id INTO shipment_id
    FROM shipments s
    WHERE s.gate_id = seed.gate_id
      AND s.created_by = seed.logist_id
      AND s.status = 'planned'
    ORDER BY id DESC
    LIMIT 1;

    INSERT INTO shipment_items (shipment_id, cargo_item_id)
    SELECT seed.shipment_id, seed.cargo_item_2_id
    WHERE NOT EXISTS (SELECT 1 FROM shipment_items si WHERE si.cargo_item_id = seed.cargo_item_2_id);
END seed;
$$;

COMMIT;
