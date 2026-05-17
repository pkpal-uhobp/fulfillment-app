-- ============================================================
-- Fulfillment App: реалистичные демонстрационные данные
-- Файл: fulfillment_demo_realistic_seed.sql
--
-- Запуск локально:
--   psql "postgres://postgres:postgres@localhost:5433/fulfillment-app?sslmode=disable" -f fulfillment_demo_realistic_seed.sql
--
-- Запуск на Railway:
--   psql "$DATABASE_URL" -f fulfillment_demo_realistic_seed.sql
--
-- Пароль для всех демо-аккаунтов:
--   Password123!
--
-- Основные аккаунты:
--   admin@transitpro.local              Администратор
--   marina.krylova@transitpro.local     Логист
--   dmitry.sokolov@transitpro.local     Логист
--   anton.orlov@transitpro.local        Рабочий склада
--   ilya.zakharov@transitpro.local      Рабочий склада
--   elena.morozova@nordtex.local        Клиент
--   pavel.danilov@cosmoline.local       Клиент
--   olga.belova@homedrop.local          Клиент
--
-- Внутри:
--   пользователи, склады, зоны, гейты, типы товаров/мест,
--   календарь приёмки, заявки, грузовые места, QR, отгрузки и истории статусов.
--
-- Скрипт идемпотентный: повторный запуск не должен создавать дубликаты
-- для записей, где есть проверки WHERE NOT EXISTS.
-- ============================================================

-- Realistic demo data for fulfillment-app.
-- Run after all migrations:
-- psql -h localhost -p 5433 -U postgres -d fulfillment-app -f scripts/seed_test_data.sql
--
-- Demo password for all users: Password123!

BEGIN;

DO $$
    <<seed>>
        DECLARE
        admin_id BIGINT;
        logist_msk_id BIGINT;
        logist_kzn_id BIGINT;
        worker_msk_id BIGINT;
        worker_kzn_id BIGINT;
        client_nordtex_id BIGINT;
        client_cosmoline_id BIGINT;
        client_homedrop_id BIGINT;

        wh_msk_id BIGINT;
        wh_kzn_id BIGINT;
        wh_spb_id BIGINT;
        wh_wb_id BIGINT;
        wh_ozon_id BIGINT;
        wh_yandex_id BIGINT;

        zone_msk_a_id BIGINT;
        zone_msk_b_id BIGINT;
        zone_msk_c_id BIGINT;
        zone_kzn_a_id BIGINT;
        gate_msk_1_id BIGINT;
        gate_msk_2_id BIGINT;
        gate_kzn_1_id BIGINT;

        product_clothes_id BIGINT;
        product_cosmetics_id BIGINT;
        product_electronics_id BIGINT;
        product_home_id BIGINT;
        box_type_id BIGINT;
        pallet_type_id BIGINT;
        bag_type_id BIGINT;

        order_1_id BIGINT;
        order_2_id BIGINT;
        order_3_id BIGINT;
        order_4_id BIGINT;
        place_1_id BIGINT;
        place_2_id BIGINT;
        place_3_id BIGINT;
        place_4_id BIGINT;
        cargo_1_id BIGINT;
        cargo_2_id BIGINT;
        cargo_3_id BIGINT;
        cargo_4_id BIGINT;
        cargo_5_id BIGINT;
        shipment_planned_id BIGINT;
        shipment_shipped_id BIGINT;

        password_hash TEXT := '$2a$12$SHdXyODTjMTP7luGwEUDAeV9l4kYdiM1fJpeJ8qJQ5LTtvv4W4rz6';
    BEGIN
        -- Users: company team + real-looking clients.
        INSERT INTO users (email, password_hash, full_name, phone, role, is_active, is_blocked)
        SELECT 'admin@transitpro.local', password_hash, 'Сергей Волков', '+79951004577', 'admin', TRUE, FALSE
        WHERE NOT EXISTS (SELECT 1 FROM users WHERE lower(email) = lower('admin@transitpro.local'));

        INSERT INTO users (email, password_hash, full_name, phone, role, is_active, is_blocked)
        SELECT 'marina.krylova@transitpro.local', password_hash, 'Марина Крылова', '+79951004578', 'logist', TRUE, FALSE
        WHERE NOT EXISTS (SELECT 1 FROM users WHERE lower(email) = lower('marina.krylova@transitpro.local'));

        INSERT INTO users (email, password_hash, full_name, phone, role, is_active, is_blocked)
        SELECT 'dmitry.sokolov@transitpro.local', password_hash, 'Дмитрий Соколов', '+79951004579', 'logist', TRUE, FALSE
        WHERE NOT EXISTS (SELECT 1 FROM users WHERE lower(email) = lower('dmitry.sokolov@transitpro.local'));

        INSERT INTO users (email, password_hash, full_name, phone, role, is_active, is_blocked)
        SELECT 'anton.orlov@transitpro.local', password_hash, 'Антон Орлов', '+79951004580', 'worker', TRUE, FALSE
        WHERE NOT EXISTS (SELECT 1 FROM users WHERE lower(email) = lower('anton.orlov@transitpro.local'));

        INSERT INTO users (email, password_hash, full_name, phone, role, is_active, is_blocked)
        SELECT 'ilya.zakharov@transitpro.local', password_hash, 'Илья Захаров', '+79951004581', 'worker', TRUE, FALSE
        WHERE NOT EXISTS (SELECT 1 FROM users WHERE lower(email) = lower('ilya.zakharov@transitpro.local'));

        INSERT INTO users (email, password_hash, full_name, phone, role, is_active, is_blocked)
        SELECT 'elena.morozova@nordtex.local', password_hash, 'Елена Морозова', '+79953331122', 'client', TRUE, FALSE
        WHERE NOT EXISTS (SELECT 1 FROM users WHERE lower(email) = lower('elena.morozova@nordtex.local'));

        INSERT INTO users (email, password_hash, full_name, phone, role, is_active, is_blocked)
        SELECT 'pavel.danilov@cosmoline.local', password_hash, 'Павел Данилов', '+79954442233', 'client', TRUE, FALSE
        WHERE NOT EXISTS (SELECT 1 FROM users WHERE lower(email) = lower('pavel.danilov@cosmoline.local'));

        INSERT INTO users (email, password_hash, full_name, phone, role, is_active, is_blocked)
        SELECT 'olga.belova@homedrop.local', password_hash, 'Ольга Белова', '+79955553344', 'client', TRUE, FALSE
        WHERE NOT EXISTS (SELECT 1 FROM users WHERE lower(email) = lower('olga.belova@homedrop.local'));

        SELECT id INTO admin_id FROM users WHERE lower(email) = lower('admin@transitpro.local');
        SELECT id INTO logist_msk_id FROM users WHERE lower(email) = lower('marina.krylova@transitpro.local');
        SELECT id INTO logist_kzn_id FROM users WHERE lower(email) = lower('dmitry.sokolov@transitpro.local');
        SELECT id INTO worker_msk_id FROM users WHERE lower(email) = lower('anton.orlov@transitpro.local');
        SELECT id INTO worker_kzn_id FROM users WHERE lower(email) = lower('ilya.zakharov@transitpro.local');
        SELECT id INTO client_nordtex_id FROM users WHERE lower(email) = lower('elena.morozova@nordtex.local');
        SELECT id INTO client_cosmoline_id FROM users WHERE lower(email) = lower('pavel.danilov@cosmoline.local');
        SELECT id INTO client_homedrop_id FROM users WHERE lower(email) = lower('olga.belova@homedrop.local');

        -- Catalogs.
        INSERT INTO product_types (name, description)
        SELECT 'Одежда', 'Одежда, текстиль и аксессуары для маркетплейсов'
        WHERE NOT EXISTS (SELECT 1 FROM product_types WHERE name = 'Одежда');

        INSERT INTO product_types (name, description)
        SELECT 'Косметика', 'Косметика, уходовые средства и парфюмерия'
        WHERE NOT EXISTS (SELECT 1 FROM product_types WHERE name = 'Косметика');

        INSERT INTO product_types (name, description)
        SELECT 'Электроника', 'Мелкая электроника и аксессуары'
        WHERE NOT EXISTS (SELECT 1 FROM product_types WHERE name = 'Электроника');

        INSERT INTO product_types (name, description)
        SELECT 'Товары для дома', 'Посуда, декор, бытовые товары и хранение'
        WHERE NOT EXISTS (SELECT 1 FROM product_types WHERE name = 'Товары для дома');

        INSERT INTO cargo_place_types (name, description)
        SELECT 'Коробка', 'Картонная коробка с единицами товара'
        WHERE NOT EXISTS (SELECT 1 FROM cargo_place_types WHERE name = 'Коробка');

        INSERT INTO cargo_place_types (name, description)
        SELECT 'Паллет', 'Паллетированное грузовое место'
        WHERE NOT EXISTS (SELECT 1 FROM cargo_place_types WHERE name = 'Паллет');

        INSERT INTO cargo_place_types (name, description)
        SELECT 'Мешок', 'Мягкая упаковка или транспортный мешок'
        WHERE NOT EXISTS (SELECT 1 FROM cargo_place_types WHERE name = 'Мешок');

        SELECT id INTO product_clothes_id FROM product_types WHERE name = 'Одежда';
        SELECT id INTO product_cosmetics_id FROM product_types WHERE name = 'Косметика';
        SELECT id INTO product_electronics_id FROM product_types WHERE name = 'Электроника';
        SELECT id INTO product_home_id FROM product_types WHERE name = 'Товары для дома';
        SELECT id INTO box_type_id FROM cargo_place_types WHERE name = 'Коробка';
        SELECT id INTO pallet_type_id FROM cargo_place_types WHERE name = 'Паллет';
        SELECT id INTO bag_type_id FROM cargo_place_types WHERE name = 'Мешок';

        -- Warehouses.
        INSERT INTO warehouses (name, warehouse_type, marketplace, city, address, is_active)
        SELECT 'TransitPro Москва Север', 'both', NULL, 'Москва', 'Дмитровское шоссе, 163А, строение 2', TRUE
        WHERE NOT EXISTS (SELECT 1 FROM warehouses WHERE name = 'TransitPro Москва Север');

        INSERT INTO warehouses (name, warehouse_type, marketplace, city, address, is_active)
        SELECT 'TransitPro Казань Восток', 'receiving', NULL, 'Казань', 'ул. Техническая, 52, склад 4', TRUE
        WHERE NOT EXISTS (SELECT 1 FROM warehouses WHERE name = 'TransitPro Казань Восток');

        INSERT INTO warehouses (name, warehouse_type, marketplace, city, address, is_active)
        SELECT 'TransitPro Санкт-Петербург Юг', 'receiving', NULL, 'Санкт-Петербург', 'Московское шоссе, 13к2, терминал B', TRUE
        WHERE NOT EXISTS (SELECT 1 FROM warehouses WHERE name = 'TransitPro Санкт-Петербург Юг');

        INSERT INTO warehouses (name, warehouse_type, marketplace, city, address, is_active)
        SELECT 'WB Коледино', 'destination', 'Wildberries', 'Москва', 'деревня Коледино, складской комплекс 1', TRUE
        WHERE NOT EXISTS (SELECT 1 FROM warehouses WHERE name = 'WB Коледино');

        INSERT INTO warehouses (name, warehouse_type, marketplace, city, address, is_active)
        SELECT 'Ozon Хоругвино', 'destination', 'Ozon', 'Москва', 'село Хоругвино, индустриальный парк Север', TRUE
        WHERE NOT EXISTS (SELECT 1 FROM warehouses WHERE name = 'Ozon Хоругвино');

        INSERT INTO warehouses (name, warehouse_type, marketplace, city, address, is_active)
        SELECT 'Яндекс Маркет Софьино', 'destination', 'Яндекс Маркет', 'Москва', 'деревня Софьино, логистический центр 3', TRUE
        WHERE NOT EXISTS (SELECT 1 FROM warehouses WHERE name = 'Яндекс Маркет Софьино');

        SELECT id INTO wh_msk_id FROM warehouses WHERE name = 'TransitPro Москва Север';
        SELECT id INTO wh_kzn_id FROM warehouses WHERE name = 'TransitPro Казань Восток';
        SELECT id INTO wh_spb_id FROM warehouses WHERE name = 'TransitPro Санкт-Петербург Юг';
        SELECT id INTO wh_wb_id FROM warehouses WHERE name = 'WB Коледино';
        SELECT id INTO wh_ozon_id FROM warehouses WHERE name = 'Ozon Хоругвино';
        SELECT id INTO wh_yandex_id FROM warehouses WHERE name = 'Яндекс Маркет Софьино';

        -- Warehouse structure.
        INSERT INTO storage_zones (warehouse_id, name, description, is_active)
        SELECT wh_msk_id, 'A-01', 'Быстрооборачиваемые коробки и мелкие партии', TRUE
        WHERE NOT EXISTS (SELECT 1 FROM storage_zones WHERE warehouse_id = wh_msk_id AND name = 'A-01');

        INSERT INTO storage_zones (warehouse_id, name, description, is_active)
        SELECT wh_msk_id, 'B-02', 'Паллетное хранение и крупногабаритные поставки', TRUE
        WHERE NOT EXISTS (SELECT 1 FROM storage_zones WHERE warehouse_id = wh_msk_id AND name = 'B-02');

        INSERT INTO storage_zones (warehouse_id, name, description, is_active)
        SELECT wh_msk_id, 'C-03', 'Зона предпродажной сортировки и маркировки', TRUE
        WHERE NOT EXISTS (SELECT 1 FROM storage_zones WHERE warehouse_id = wh_msk_id AND name = 'C-03');

        INSERT INTO storage_zones (warehouse_id, name, description, is_active)
        SELECT wh_kzn_id, 'KZN-A01', 'Приемка региональных партий Казань', TRUE
        WHERE NOT EXISTS (SELECT 1 FROM storage_zones WHERE warehouse_id = wh_kzn_id AND name = 'KZN-A01');

        INSERT INTO gates (warehouse_id, name, is_active)
        SELECT wh_msk_id, 'MSK-G01', TRUE
        WHERE NOT EXISTS (SELECT 1 FROM gates WHERE warehouse_id = wh_msk_id AND name = 'MSK-G01');

        INSERT INTO gates (warehouse_id, name, is_active)
        SELECT wh_msk_id, 'MSK-G02', TRUE
        WHERE NOT EXISTS (SELECT 1 FROM gates WHERE warehouse_id = wh_msk_id AND name = 'MSK-G02');

        INSERT INTO gates (warehouse_id, name, is_active)
        SELECT wh_kzn_id, 'KZN-G01', TRUE
        WHERE NOT EXISTS (SELECT 1 FROM gates WHERE warehouse_id = wh_kzn_id AND name = 'KZN-G01');

        SELECT id INTO zone_msk_a_id FROM storage_zones WHERE warehouse_id = wh_msk_id AND name = 'A-01';
        SELECT id INTO zone_msk_b_id FROM storage_zones WHERE warehouse_id = wh_msk_id AND name = 'B-02';
        SELECT id INTO zone_msk_c_id FROM storage_zones WHERE warehouse_id = wh_msk_id AND name = 'C-03';
        SELECT id INTO zone_kzn_a_id FROM storage_zones WHERE warehouse_id = wh_kzn_id AND name = 'KZN-A01';
        SELECT id INTO gate_msk_1_id FROM gates WHERE warehouse_id = wh_msk_id AND name = 'MSK-G01';
        SELECT id INTO gate_msk_2_id FROM gates WHERE warehouse_id = wh_msk_id AND name = 'MSK-G02';
        SELECT id INTO gate_kzn_1_id FROM gates WHERE warehouse_id = wh_kzn_id AND name = 'KZN-G01';

        INSERT INTO worker_profiles (user_id, warehouse_id)
        SELECT worker_msk_id, wh_msk_id
        WHERE NOT EXISTS (SELECT 1 FROM worker_profiles WHERE user_id = worker_msk_id);

        INSERT INTO worker_profiles (user_id, warehouse_id)
        SELECT worker_kzn_id, wh_kzn_id
        WHERE NOT EXISTS (SELECT 1 FROM worker_profiles WHERE user_id = worker_kzn_id);

        -- Pickup calendar: available days, busy days and one closed day.
        INSERT INTO pickup_capacity (warehouse_id, pickup_date, max_orders, current_orders, is_closed)
        SELECT wh_msk_id, CURRENT_DATE + INTERVAL '1 day', 18, 7, FALSE
        WHERE NOT EXISTS (SELECT 1 FROM pickup_capacity WHERE warehouse_id = wh_msk_id AND pickup_date = CURRENT_DATE + INTERVAL '1 day');

        INSERT INTO pickup_capacity (warehouse_id, pickup_date, max_orders, current_orders, is_closed)
        SELECT wh_msk_id, CURRENT_DATE + INTERVAL '2 days', 18, 16, FALSE
        WHERE NOT EXISTS (SELECT 1 FROM pickup_capacity WHERE warehouse_id = wh_msk_id AND pickup_date = CURRENT_DATE + INTERVAL '2 days');

        INSERT INTO pickup_capacity (warehouse_id, pickup_date, max_orders, current_orders, is_closed)
        SELECT wh_kzn_id, CURRENT_DATE + INTERVAL '2 days', 10, 4, FALSE
        WHERE NOT EXISTS (SELECT 1 FROM pickup_capacity WHERE warehouse_id = wh_kzn_id AND pickup_date = CURRENT_DATE + INTERVAL '2 days');

        INSERT INTO pickup_calendar_blocks (warehouse_id, blocked_date, reason, created_by)
        SELECT wh_msk_id, CURRENT_DATE + INTERVAL '5 days', 'Плановая инвентаризация зоны приемки', logist_msk_id
        WHERE NOT EXISTS (SELECT 1 FROM pickup_calendar_blocks WHERE warehouse_id = wh_msk_id AND blocked_date = CURRENT_DATE + INTERVAL '5 days');

        -- Orders.
        INSERT INTO orders (
            client_id, receiving_warehouse_id, destination_warehouse_id, product_type_id,
            handover_type, self_delivery_date, self_delivery_time_from, self_delivery_time_to,
            status, comment
        )
        SELECT
            client_nordtex_id, wh_msk_id, wh_wb_id, product_clothes_id,
            'self_delivery', CURRENT_DATE + INTERVAL '1 day', '10:00', '12:00',
            'received', 'seed:TransitPro:NordTex self delivery to WB'
        WHERE NOT EXISTS (SELECT 1 FROM orders WHERE comment = 'seed:TransitPro:NordTex self delivery to WB');

        INSERT INTO orders (
            client_id, receiving_warehouse_id, destination_warehouse_id, product_type_id,
            handover_type, status, comment
        )
        SELECT
            client_cosmoline_id, wh_msk_id, wh_ozon_id, product_cosmetics_id,
            'pickup', 'waiting_pickup', 'seed:TransitPro:CosmoLine pickup to Ozon'
        WHERE NOT EXISTS (SELECT 1 FROM orders WHERE comment = 'seed:TransitPro:CosmoLine pickup to Ozon');

        INSERT INTO orders (
            client_id, receiving_warehouse_id, destination_warehouse_id, product_type_id,
            handover_type, self_delivery_date, self_delivery_time_from, self_delivery_time_to,
            status, comment
        )
        SELECT
            client_homedrop_id, wh_msk_id, wh_yandex_id, product_home_id,
            'self_delivery', CURRENT_DATE - INTERVAL '1 day', '13:00', '15:00',
            'stored', 'seed:TransitPro:HomeDrop stored order to Yandex'
        WHERE NOT EXISTS (SELECT 1 FROM orders WHERE comment = 'seed:TransitPro:HomeDrop stored order to Yandex');

        INSERT INTO orders (
            client_id, receiving_warehouse_id, destination_warehouse_id, product_type_id,
            handover_type, status, comment
        )
        SELECT
            client_nordtex_id, wh_kzn_id, wh_wb_id, product_clothes_id,
            'pickup', 'assigned_to_shipping', 'seed:TransitPro:NordTex regional shipment'
        WHERE NOT EXISTS (SELECT 1 FROM orders WHERE comment = 'seed:TransitPro:NordTex regional shipment');

        SELECT id INTO order_1_id FROM orders WHERE comment = 'seed:TransitPro:NordTex self delivery to WB';
        SELECT id INTO order_2_id FROM orders WHERE comment = 'seed:TransitPro:CosmoLine pickup to Ozon';
        SELECT id INTO order_3_id FROM orders WHERE comment = 'seed:TransitPro:HomeDrop stored order to Yandex';
        SELECT id INTO order_4_id FROM orders WHERE comment = 'seed:TransitPro:NordTex regional shipment';

        -- Pickup requests.
        INSERT INTO pickup_requests (
            order_id, pickup_address, pickup_date, pickup_time_from, pickup_time_to,
            contact_name, contact_phone, status, assigned_logist_id, comment
        )
        SELECT
            order_2_id, 'Москва, ул. Рябиновая, 43Б, склад CosmoLine', CURRENT_DATE + INTERVAL '2 days',
            '14:00', '17:00', 'Павел Данилов', '+79954442233', 'assigned', logist_msk_id,
            'seed:TransitPro:CosmoLine pickup request'
        WHERE NOT EXISTS (SELECT 1 FROM pickup_requests WHERE order_id = order_2_id);

        INSERT INTO pickup_requests (
            order_id, pickup_address, pickup_date, pickup_time_from, pickup_time_to,
            contact_name, contact_phone, status, assigned_logist_id, comment
        )
        SELECT
            order_4_id, 'Казань, ул. Восстания, 100, ворота 3', CURRENT_DATE - INTERVAL '1 day',
            '09:00', '11:00', 'Елена Морозова', '+79953331122', 'picked_up', logist_kzn_id,
            'seed:TransitPro:NordTex Kazan pickup request'
        WHERE NOT EXISTS (SELECT 1 FROM pickup_requests WHERE order_id = order_4_id);

        -- Declared cargo places.
        INSERT INTO order_cargo_places (order_id, cargo_place_type_id, quantity, weight_per_place_kg, length_cm, width_cm, height_cm, comment)
        SELECT order_1_id, box_type_id, 2, 7.40, 60, 40, 35, 'seed:TransitPro:order1 boxes'
        WHERE NOT EXISTS (SELECT 1 FROM order_cargo_places WHERE order_id = order_1_id AND comment = 'seed:TransitPro:order1 boxes');

        INSERT INTO order_cargo_places (order_id, cargo_place_type_id, quantity, weight_per_place_kg, length_cm, width_cm, height_cm, comment)
        SELECT order_2_id, box_type_id, 3, 4.80, 45, 35, 25, 'seed:TransitPro:order2 cosmetics boxes'
        WHERE NOT EXISTS (SELECT 1 FROM order_cargo_places WHERE order_id = order_2_id AND comment = 'seed:TransitPro:order2 cosmetics boxes');

        INSERT INTO order_cargo_places (order_id, cargo_place_type_id, quantity, weight_per_place_kg, length_cm, width_cm, height_cm, comment)
        SELECT order_3_id, pallet_type_id, 1, 86.00, 120, 80, 135, 'seed:TransitPro:order3 pallet'
        WHERE NOT EXISTS (SELECT 1 FROM order_cargo_places WHERE order_id = order_3_id AND comment = 'seed:TransitPro:order3 pallet');

        INSERT INTO order_cargo_places (order_id, cargo_place_type_id, quantity, weight_per_place_kg, length_cm, width_cm, height_cm, comment)
        SELECT order_4_id, bag_type_id, 2, 9.20, 70, 45, 30, 'seed:TransitPro:order4 bags'
        WHERE NOT EXISTS (SELECT 1 FROM order_cargo_places WHERE order_id = order_4_id AND comment = 'seed:TransitPro:order4 bags');

        SELECT id INTO place_1_id FROM order_cargo_places WHERE order_id = order_1_id AND comment = 'seed:TransitPro:order1 boxes' LIMIT 1;
        SELECT id INTO place_2_id FROM order_cargo_places WHERE order_id = order_2_id AND comment = 'seed:TransitPro:order2 cosmetics boxes' LIMIT 1;
        SELECT id INTO place_3_id FROM order_cargo_places WHERE order_id = order_3_id AND comment = 'seed:TransitPro:order3 pallet' LIMIT 1;
        SELECT id INTO place_4_id FROM order_cargo_places WHERE order_id = order_4_id AND comment = 'seed:TransitPro:order4 bags' LIMIT 1;

        -- Cargo items with QR codes for frontend scan page.
        INSERT INTO cargo_items (order_id, order_cargo_place_id, cargo_place_type_id, qr_code, status, storage_zone_id, gate_id, received_by, received_at, comment)
        SELECT order_1_id, place_1_id, box_type_id, 'QR-TPRO-MSK-240001', 'accepted', NULL, NULL, worker_msk_id, NOW() - INTERVAL '3 hours', 'seed:TransitPro:accepted box NordTex'
        WHERE NOT EXISTS (SELECT 1 FROM cargo_items WHERE qr_code = 'QR-TPRO-MSK-240001');

        INSERT INTO cargo_items (order_id, order_cargo_place_id, cargo_place_type_id, qr_code, status, storage_zone_id, gate_id, received_by, received_at, comment)
        SELECT order_1_id, place_1_id, box_type_id, 'QR-TPRO-MSK-240002', 'stored', zone_msk_a_id, NULL, worker_msk_id, NOW() - INTERVAL '2 hours', 'seed:TransitPro:stored box NordTex'
        WHERE NOT EXISTS (SELECT 1 FROM cargo_items WHERE qr_code = 'QR-TPRO-MSK-240002');

        INSERT INTO cargo_items (order_id, order_cargo_place_id, cargo_place_type_id, qr_code, status, storage_zone_id, gate_id, received_by, received_at, comment)
        SELECT order_3_id, place_3_id, pallet_type_id, 'QR-TPRO-MSK-240003', 'ready_to_ship', zone_msk_b_id, gate_msk_1_id, worker_msk_id, NOW() - INTERVAL '1 day', 'seed:TransitPro:ready pallet HomeDrop'
        WHERE NOT EXISTS (SELECT 1 FROM cargo_items WHERE qr_code = 'QR-TPRO-MSK-240003');

        INSERT INTO cargo_items (order_id, order_cargo_place_id, cargo_place_type_id, qr_code, status, storage_zone_id, gate_id, received_by, shipped_by, received_at, shipped_at, comment)
        SELECT order_4_id, place_4_id, bag_type_id, 'QR-TPRO-KZN-240001', 'shipped', zone_kzn_a_id, gate_kzn_1_id, worker_kzn_id, worker_kzn_id, NOW() - INTERVAL '2 days', NOW() - INTERVAL '6 hours', 'seed:TransitPro:shipped bag NordTex'
        WHERE NOT EXISTS (SELECT 1 FROM cargo_items WHERE qr_code = 'QR-TPRO-KZN-240001');

        INSERT INTO cargo_items (order_id, order_cargo_place_id, cargo_place_type_id, qr_code, status, storage_zone_id, gate_id, received_by, received_at, comment)
        SELECT order_2_id, place_2_id, box_type_id, 'QR-TPRO-MSK-240004', 'accepted', NULL, NULL, worker_msk_id, NOW() - INTERVAL '30 minutes', 'seed:TransitPro:pickup cosmetics accepted'
        WHERE NOT EXISTS (SELECT 1 FROM cargo_items WHERE qr_code = 'QR-TPRO-MSK-240004');

        SELECT id INTO cargo_1_id FROM cargo_items WHERE qr_code = 'QR-TPRO-MSK-240001';
        SELECT id INTO cargo_2_id FROM cargo_items WHERE qr_code = 'QR-TPRO-MSK-240002';
        SELECT id INTO cargo_3_id FROM cargo_items WHERE qr_code = 'QR-TPRO-MSK-240003';
        SELECT id INTO cargo_4_id FROM cargo_items WHERE qr_code = 'QR-TPRO-KZN-240001';
        SELECT id INTO cargo_5_id FROM cargo_items WHERE qr_code = 'QR-TPRO-MSK-240004';

        -- Status history.
        INSERT INTO order_status_history (order_id, old_status, new_status, changed_by, comment)
        SELECT order_1_id, NULL, 'created', client_nordtex_id, 'seed:TransitPro:order1 created'
        WHERE NOT EXISTS (SELECT 1 FROM order_status_history WHERE order_id = order_1_id AND comment = 'seed:TransitPro:order1 created');

        INSERT INTO order_status_history (order_id, old_status, new_status, changed_by, comment)
        SELECT order_1_id, 'created', 'received', worker_msk_id, 'seed:TransitPro:order1 received'
        WHERE NOT EXISTS (SELECT 1 FROM order_status_history WHERE order_id = order_1_id AND comment = 'seed:TransitPro:order1 received');

        INSERT INTO order_status_history (order_id, old_status, new_status, changed_by, comment)
        SELECT order_3_id, 'received', 'stored', worker_msk_id, 'seed:TransitPro:order3 stored'
        WHERE NOT EXISTS (SELECT 1 FROM order_status_history WHERE order_id = order_3_id AND comment = 'seed:TransitPro:order3 stored');

        INSERT INTO order_status_history (order_id, old_status, new_status, changed_by, comment)
        SELECT order_4_id, 'stored', 'assigned_to_shipping', logist_kzn_id, 'seed:TransitPro:order4 assigned to shipping'
        WHERE NOT EXISTS (SELECT 1 FROM order_status_history WHERE order_id = order_4_id AND comment = 'seed:TransitPro:order4 assigned to shipping');

        INSERT INTO cargo_status_history (cargo_item_id, old_status, new_status, changed_by, comment)
        SELECT cargo_1_id, NULL, 'accepted', worker_msk_id, 'seed:TransitPro:cargo1 accepted'
        WHERE NOT EXISTS (SELECT 1 FROM cargo_status_history WHERE cargo_item_id = cargo_1_id AND comment = 'seed:TransitPro:cargo1 accepted');

        INSERT INTO cargo_status_history (cargo_item_id, old_status, new_status, changed_by, comment)
        SELECT cargo_2_id, 'accepted', 'stored', worker_msk_id, 'seed:TransitPro:cargo2 stored'
        WHERE NOT EXISTS (SELECT 1 FROM cargo_status_history WHERE cargo_item_id = cargo_2_id AND comment = 'seed:TransitPro:cargo2 stored');

        INSERT INTO cargo_status_history (cargo_item_id, old_status, new_status, changed_by, comment)
        SELECT cargo_3_id, 'stored', 'ready_to_ship', logist_msk_id, 'seed:TransitPro:cargo3 ready'
        WHERE NOT EXISTS (SELECT 1 FROM cargo_status_history WHERE cargo_item_id = cargo_3_id AND comment = 'seed:TransitPro:cargo3 ready');

        INSERT INTO cargo_status_history (cargo_item_id, old_status, new_status, changed_by, comment)
        SELECT cargo_4_id, 'ready_to_ship', 'shipped', worker_kzn_id, 'seed:TransitPro:cargo4 shipped'
        WHERE NOT EXISTS (SELECT 1 FROM cargo_status_history WHERE cargo_item_id = cargo_4_id AND comment = 'seed:TransitPro:cargo4 shipped');

        -- Shipments.
        INSERT INTO shipments (destination_warehouse_id, gate_id, planned_departure_at, status, created_by)
        SELECT wh_yandex_id, gate_msk_1_id, NOW() + INTERVAL '1 day', 'planned', logist_msk_id
        WHERE NOT EXISTS (
            SELECT 1 FROM shipments
            WHERE destination_warehouse_id = wh_yandex_id
              AND gate_id = gate_msk_1_id
              AND created_by = logist_msk_id
              AND status = 'planned'
        );

        SELECT id INTO shipment_planned_id
        FROM shipments
        WHERE destination_warehouse_id = wh_yandex_id
          AND gate_id = gate_msk_1_id
          AND created_by = logist_msk_id
          AND status = 'planned'
        ORDER BY id DESC
        LIMIT 1;

        INSERT INTO shipment_items (shipment_id, cargo_item_id)
        SELECT shipment_planned_id, cargo_3_id
        WHERE NOT EXISTS (SELECT 1 FROM shipment_items WHERE cargo_item_id = cargo_3_id);

        INSERT INTO shipments (destination_warehouse_id, gate_id, planned_departure_at, actual_departure_at, status, created_by)
        SELECT wh_wb_id, gate_kzn_1_id, NOW() - INTERVAL '8 hours', NOW() - INTERVAL '6 hours', 'shipped', logist_kzn_id
        WHERE NOT EXISTS (
            SELECT 1 FROM shipments
            WHERE destination_warehouse_id = wh_wb_id
              AND gate_id = gate_kzn_1_id
              AND created_by = logist_kzn_id
              AND status = 'shipped'
        );

        SELECT id INTO shipment_shipped_id
        FROM shipments
        WHERE destination_warehouse_id = wh_wb_id
          AND gate_id = gate_kzn_1_id
          AND created_by = logist_kzn_id
          AND status = 'shipped'
        ORDER BY id DESC
        LIMIT 1;

        INSERT INTO shipment_items (shipment_id, cargo_item_id)
        SELECT shipment_shipped_id, cargo_4_id
        WHERE NOT EXISTS (SELECT 1 FROM shipment_items WHERE cargo_item_id = cargo_4_id);
    END seed;
$$;

COMMIT;
