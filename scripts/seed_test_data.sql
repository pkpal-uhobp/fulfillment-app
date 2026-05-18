-- ============================================================
-- Fulfillment App: расширенные реалистичные демонстрационные данные
-- Файл: scripts/seed_test_data.sql
--
-- Запуск локально:
--   psql "postgres://postgres:postgres@localhost:5433/fulfillment-app?sslmode=disable" -f scripts/seed_test_data.sql
--
-- Пароль для всех демо-аккаунтов:
--   Password123!
--
-- Что добавлено:
--   1) У каждого склада автоматически создаются зоны хранения и зоны/ворота отгрузки.
--   2) Больше клиентов, заявок и грузовых мест.
--   3) В каждой новой заявке создаётся от 1 до 7 грузовых мест.
--   4) Для принятых/хранящихся/отгружаемых заявок создаются QR-грузовые места.
--   5) Данные идемпотентные: повторный запуск не должен плодить дубли.
-- ============================================================

BEGIN;

DO $$
DECLARE
    password_hash TEXT := '$2a$12$SHdXyODTjMTP7luGwEUDAeV9l4kYdiM1fJpeJ8qJQ5LTtvv4W4rz6';

    admin_id BIGINT;

    logist_msk_id BIGINT;
    logist_kzn_id BIGINT;
    logist_spb_id BIGINT;

    worker_msk_id BIGINT;
    worker_kzn_id BIGINT;
    worker_spb_id BIGINT;

    client_nordtex_id BIGINT;
    client_cosmoline_id BIGINT;
    client_homedrop_id BIGINT;
    client_techline_id BIGINT;
    client_babybox_id BIGINT;
    client_sportway_id BIGINT;

    wh_msk_id BIGINT;
    wh_kzn_id BIGINT;
    wh_spb_id BIGINT;
    wh_wb_id BIGINT;
    wh_ozon_id BIGINT;
    wh_yandex_id BIGINT;

    product_clothes_id BIGINT;
    product_cosmetics_id BIGINT;
    product_electronics_id BIGINT;
    product_home_id BIGINT;
    product_kids_id BIGINT;
    product_sport_id BIGINT;

    box_type_id BIGINT;
    pallet_type_id BIGINT;
    bag_type_id BIGINT;
    crate_type_id BIGINT;

    wh_rec RECORD;

    client_ids BIGINT[];
    receiving_wh_ids BIGINT[];
    destination_wh_ids BIGINT[];
    product_type_ids BIGINT[];
    cargo_type_ids BIGINT[];
    logist_ids BIGINT[];
    worker_ids BIGINT[];

    i INT;
    j INT;
    cargo_count INT;

    v_order_id BIGINT;
    place_id BIGINT;
    v_cargo_item_id BIGINT;
    shipment_id BIGINT;

    order_comment TEXT;
    place_comment TEXT;
    pickup_comment TEXT;
    cargo_comment TEXT;
    qr TEXT;

    selected_client_id BIGINT;
    selected_receiving_wh_id BIGINT;
    selected_destination_wh_id BIGINT;
    selected_product_type_id BIGINT;
    selected_cargo_type_id BIGINT;
    selected_logist_id BIGINT;
    selected_worker_id BIGINT;

    handover TEXT;
    order_status TEXT;
    pickup_status TEXT;
    cargo_status TEXT;

    delivery_date DATE;
    time_from TIME;
    time_to TIME;

    storage_zone_id BIGINT;
    gate_id BIGINT;

    item_weight NUMERIC(10, 2);
    item_length NUMERIC(10, 2);
    item_width NUMERIC(10, 2);
    item_height NUMERIC(10, 2);

    shipment_rec RECORD;
BEGIN
    -- ========================================================
    -- Users.
    -- ========================================================
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
    SELECT 'anna.romanova@transitpro.local', password_hash, 'Анна Романова', '+79951004582', 'logist', TRUE, FALSE
    WHERE NOT EXISTS (SELECT 1 FROM users WHERE lower(email) = lower('anna.romanova@transitpro.local'));

    INSERT INTO users (email, password_hash, full_name, phone, role, is_active, is_blocked)
    SELECT 'anton.orlov@transitpro.local', password_hash, 'Антон Орлов', '+79951004580', 'worker', TRUE, FALSE
    WHERE NOT EXISTS (SELECT 1 FROM users WHERE lower(email) = lower('anton.orlov@transitpro.local'));

    INSERT INTO users (email, password_hash, full_name, phone, role, is_active, is_blocked)
    SELECT 'ilya.zakharov@transitpro.local', password_hash, 'Илья Захаров', '+79951004581', 'worker', TRUE, FALSE
    WHERE NOT EXISTS (SELECT 1 FROM users WHERE lower(email) = lower('ilya.zakharov@transitpro.local'));

    INSERT INTO users (email, password_hash, full_name, phone, role, is_active, is_blocked)
    SELECT 'egor.melnikov@transitpro.local', password_hash, 'Егор Мельников', '+79951004583', 'worker', TRUE, FALSE
    WHERE NOT EXISTS (SELECT 1 FROM users WHERE lower(email) = lower('egor.melnikov@transitpro.local'));

    INSERT INTO users (email, password_hash, full_name, phone, role, is_active, is_blocked)
    SELECT 'elena.morozova@nordtex.local', password_hash, 'Елена Морозова', '+79953331122', 'client', TRUE, FALSE
    WHERE NOT EXISTS (SELECT 1 FROM users WHERE lower(email) = lower('elena.morozova@nordtex.local'));

    INSERT INTO users (email, password_hash, full_name, phone, role, is_active, is_blocked)
    SELECT 'pavel.danilov@cosmoline.local', password_hash, 'Павел Данилов', '+79954442233', 'client', TRUE, FALSE
    WHERE NOT EXISTS (SELECT 1 FROM users WHERE lower(email) = lower('pavel.danilov@cosmoline.local'));

    INSERT INTO users (email, password_hash, full_name, phone, role, is_active, is_blocked)
    SELECT 'olga.belova@homedrop.local', password_hash, 'Ольга Белова', '+79955553344', 'client', TRUE, FALSE
    WHERE NOT EXISTS (SELECT 1 FROM users WHERE lower(email) = lower('olga.belova@homedrop.local'));

    INSERT INTO users (email, password_hash, full_name, phone, role, is_active, is_blocked)
    SELECT 'mikhail.fedorov@techline.local', password_hash, 'Михаил Федоров', '+79956664455', 'client', TRUE, FALSE
    WHERE NOT EXISTS (SELECT 1 FROM users WHERE lower(email) = lower('mikhail.fedorov@techline.local'));

    INSERT INTO users (email, password_hash, full_name, phone, role, is_active, is_blocked)
    SELECT 'irina.nikitina@babybox.local', password_hash, 'Ирина Никитина', '+79957775566', 'client', TRUE, FALSE
    WHERE NOT EXISTS (SELECT 1 FROM users WHERE lower(email) = lower('irina.nikitina@babybox.local'));

    INSERT INTO users (email, password_hash, full_name, phone, role, is_active, is_blocked)
    SELECT 'roman.karpov@sportway.local', password_hash, 'Роман Карпов', '+79958886677', 'client', TRUE, FALSE
    WHERE NOT EXISTS (SELECT 1 FROM users WHERE lower(email) = lower('roman.karpov@sportway.local'));

    SELECT id INTO admin_id FROM users WHERE lower(email) = lower('admin@transitpro.local');
    SELECT id INTO logist_msk_id FROM users WHERE lower(email) = lower('marina.krylova@transitpro.local');
    SELECT id INTO logist_kzn_id FROM users WHERE lower(email) = lower('dmitry.sokolov@transitpro.local');
    SELECT id INTO logist_spb_id FROM users WHERE lower(email) = lower('anna.romanova@transitpro.local');
    SELECT id INTO worker_msk_id FROM users WHERE lower(email) = lower('anton.orlov@transitpro.local');
    SELECT id INTO worker_kzn_id FROM users WHERE lower(email) = lower('ilya.zakharov@transitpro.local');
    SELECT id INTO worker_spb_id FROM users WHERE lower(email) = lower('egor.melnikov@transitpro.local');

    SELECT id INTO client_nordtex_id FROM users WHERE lower(email) = lower('elena.morozova@nordtex.local');
    SELECT id INTO client_cosmoline_id FROM users WHERE lower(email) = lower('pavel.danilov@cosmoline.local');
    SELECT id INTO client_homedrop_id FROM users WHERE lower(email) = lower('olga.belova@homedrop.local');
    SELECT id INTO client_techline_id FROM users WHERE lower(email) = lower('mikhail.fedorov@techline.local');
    SELECT id INTO client_babybox_id FROM users WHERE lower(email) = lower('irina.nikitina@babybox.local');
    SELECT id INTO client_sportway_id FROM users WHERE lower(email) = lower('roman.karpov@sportway.local');

    -- ========================================================
    -- Catalogs.
    -- ========================================================
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

    INSERT INTO product_types (name, description)
    SELECT 'Детские товары', 'Игрушки, текстиль и товары для детей'
    WHERE NOT EXISTS (SELECT 1 FROM product_types WHERE name = 'Детские товары');

    INSERT INTO product_types (name, description)
    SELECT 'Спорттовары', 'Спортивный инвентарь, одежда и аксессуары'
    WHERE NOT EXISTS (SELECT 1 FROM product_types WHERE name = 'Спорттовары');

    INSERT INTO cargo_place_types (name, description)
    SELECT 'Коробка', 'Картонная коробка с единицами товара'
    WHERE NOT EXISTS (SELECT 1 FROM cargo_place_types WHERE name = 'Коробка');

    INSERT INTO cargo_place_types (name, description)
    SELECT 'Паллет', 'Паллетированное грузовое место'
    WHERE NOT EXISTS (SELECT 1 FROM cargo_place_types WHERE name = 'Паллет');

    INSERT INTO cargo_place_types (name, description)
    SELECT 'Мешок', 'Мягкая упаковка или транспортный мешок'
    WHERE NOT EXISTS (SELECT 1 FROM cargo_place_types WHERE name = 'Мешок');

    INSERT INTO cargo_place_types (name, description)
    SELECT 'Ящик', 'Пластиковый или деревянный транспортный ящик'
    WHERE NOT EXISTS (SELECT 1 FROM cargo_place_types WHERE name = 'Ящик');

    SELECT id INTO product_clothes_id FROM product_types WHERE name = 'Одежда';
    SELECT id INTO product_cosmetics_id FROM product_types WHERE name = 'Косметика';
    SELECT id INTO product_electronics_id FROM product_types WHERE name = 'Электроника';
    SELECT id INTO product_home_id FROM product_types WHERE name = 'Товары для дома';
    SELECT id INTO product_kids_id FROM product_types WHERE name = 'Детские товары';
    SELECT id INTO product_sport_id FROM product_types WHERE name = 'Спорттовары';

    SELECT id INTO box_type_id FROM cargo_place_types WHERE name = 'Коробка';
    SELECT id INTO pallet_type_id FROM cargo_place_types WHERE name = 'Паллет';
    SELECT id INTO bag_type_id FROM cargo_place_types WHERE name = 'Мешок';
    SELECT id INTO crate_type_id FROM cargo_place_types WHERE name = 'Ящик';

    -- ========================================================
    -- Warehouses.
    -- ========================================================
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

    -- ========================================================
    -- Warehouse structure.
    -- Важно: цикл проходит по ВСЕМ складам в базе, поэтому зоны
    -- хранения и зоны/ворота отгрузки появятся у каждого склада.
    -- ========================================================
    FOR wh_rec IN SELECT id, name FROM warehouses LOOP
        INSERT INTO storage_zones (warehouse_id, name, description, is_active)
        SELECT wh_rec.id, 'Хранение A', 'Основная зона хранения коробок и мелких партий', TRUE
        WHERE NOT EXISTS (
            SELECT 1 FROM storage_zones WHERE warehouse_id = wh_rec.id AND name = 'Хранение A'
        );

        INSERT INTO storage_zones (warehouse_id, name, description, is_active)
        SELECT wh_rec.id, 'Хранение B', 'Паллетное хранение и средние партии', TRUE
        WHERE NOT EXISTS (
            SELECT 1 FROM storage_zones WHERE warehouse_id = wh_rec.id AND name = 'Хранение B'
        );

        INSERT INTO storage_zones (warehouse_id, name, description, is_active)
        SELECT wh_rec.id, 'Буфер приемки', 'Буферная зона после приемки и перед сортировкой', TRUE
        WHERE NOT EXISTS (
            SELECT 1 FROM storage_zones WHERE warehouse_id = wh_rec.id AND name = 'Буфер приемки'
        );

        INSERT INTO gates (warehouse_id, name, is_active)
        SELECT wh_rec.id, 'Отгрузка 1', TRUE
        WHERE NOT EXISTS (
            SELECT 1 FROM gates WHERE warehouse_id = wh_rec.id AND name = 'Отгрузка 1'
        );

        INSERT INTO gates (warehouse_id, name, is_active)
        SELECT wh_rec.id, 'Отгрузка 2', TRUE
        WHERE NOT EXISTS (
            SELECT 1 FROM gates WHERE warehouse_id = wh_rec.id AND name = 'Отгрузка 2'
        );

        INSERT INTO gates (warehouse_id, name, is_active)
        SELECT wh_rec.id, 'Приемка 1', TRUE
        WHERE NOT EXISTS (
            SELECT 1 FROM gates WHERE warehouse_id = wh_rec.id AND name = 'Приемка 1'
        );
    END LOOP;

    -- Worker profiles.
    INSERT INTO worker_profiles (user_id, warehouse_id)
    SELECT worker_msk_id, wh_msk_id
    WHERE NOT EXISTS (SELECT 1 FROM worker_profiles WHERE user_id = worker_msk_id);

    INSERT INTO worker_profiles (user_id, warehouse_id)
    SELECT worker_kzn_id, wh_kzn_id
    WHERE NOT EXISTS (SELECT 1 FROM worker_profiles WHERE user_id = worker_kzn_id);

    INSERT INTO worker_profiles (user_id, warehouse_id)
    SELECT worker_spb_id, wh_spb_id
    WHERE NOT EXISTS (SELECT 1 FROM worker_profiles WHERE user_id = worker_spb_id);

    -- Pickup calendar.
    FOR wh_rec IN
        SELECT id FROM warehouses WHERE warehouse_type IN ('receiving', 'both') ORDER BY id
    LOOP
        FOR i IN 1..7 LOOP
            INSERT INTO pickup_capacity (warehouse_id, pickup_date, max_orders, current_orders, is_closed)
            SELECT wh_rec.id, CURRENT_DATE + i, 12 + ((wh_rec.id + i) % 8), ((wh_rec.id + i) % 5), FALSE
            WHERE NOT EXISTS (
                SELECT 1
                FROM pickup_capacity
                WHERE warehouse_id = wh_rec.id
                  AND pickup_date = CURRENT_DATE + i
            );
        END LOOP;
    END LOOP;

    INSERT INTO pickup_calendar_blocks (warehouse_id, blocked_date, reason, created_by)
    SELECT wh_msk_id, CURRENT_DATE + 5, 'Плановая инвентаризация зоны приемки', logist_msk_id
    WHERE NOT EXISTS (
        SELECT 1 FROM pickup_calendar_blocks
        WHERE warehouse_id = wh_msk_id AND blocked_date = CURRENT_DATE + 5
    );

    -- ========================================================
    -- Orders and cargo places.
    -- Создаём 28 заявок; в каждой заявке от 1 до 7 грузовых мест.
    -- ========================================================
    client_ids := ARRAY[
        client_nordtex_id,
        client_cosmoline_id,
        client_homedrop_id,
        client_techline_id,
        client_babybox_id,
        client_sportway_id
    ];

    receiving_wh_ids := ARRAY[wh_msk_id, wh_kzn_id, wh_spb_id];
    destination_wh_ids := ARRAY[wh_wb_id, wh_ozon_id, wh_yandex_id];
    product_type_ids := ARRAY[
        product_clothes_id,
        product_cosmetics_id,
        product_electronics_id,
        product_home_id,
        product_kids_id,
        product_sport_id
    ];
    cargo_type_ids := ARRAY[box_type_id, pallet_type_id, bag_type_id, crate_type_id];
    logist_ids := ARRAY[logist_msk_id, logist_kzn_id, logist_spb_id];
    worker_ids := ARRAY[worker_msk_id, worker_kzn_id, worker_spb_id];

    FOR i IN 1..28 LOOP
        order_comment := format('seed:v2:order:%s', lpad(i::TEXT, 2, '0'));
        place_comment := format('seed:v2:order:%s:declared-cargo', lpad(i::TEXT, 2, '0'));
        pickup_comment := format('seed:v2:order:%s:pickup', lpad(i::TEXT, 2, '0'));

        selected_client_id := client_ids[((i - 1) % array_length(client_ids, 1)) + 1];
        selected_receiving_wh_id := receiving_wh_ids[((i - 1) % array_length(receiving_wh_ids, 1)) + 1];
        selected_destination_wh_id := destination_wh_ids[((i - 1) % array_length(destination_wh_ids, 1)) + 1];
        selected_product_type_id := product_type_ids[((i - 1) % array_length(product_type_ids, 1)) + 1];
        selected_cargo_type_id := cargo_type_ids[((i - 1) % array_length(cargo_type_ids, 1)) + 1];
        selected_logist_id := logist_ids[((i - 1) % array_length(logist_ids, 1)) + 1];
        selected_worker_id := worker_ids[((i - 1) % array_length(worker_ids, 1)) + 1];

        cargo_count := ((i * 5) % 7) + 1;

        IF i % 2 = 0 THEN
            handover := 'pickup';
            order_status := (ARRAY[
                'waiting_pickup',
                'received',
                'stored',
                'assigned_to_shipping',
                'shipped',
                'delivered'
            ])[((i - 1) % 6) + 1];
        ELSE
            handover := 'self_delivery';
            order_status := (ARRAY[
                'waiting_delivery',
                'received',
                'stored',
                'assigned_to_shipping',
                'shipped',
                'delivered'
            ])[((i - 1) % 6) + 1];
        END IF;

        delivery_date := CURRENT_DATE + ((i % 11) - 4);
        time_from := make_time(9 + (i % 7), 0, 0);
        time_to := time_from + INTERVAL '2 hours';

        IF selected_cargo_type_id = pallet_type_id THEN
            item_weight := 72 + (i % 6) * 8;
            item_length := 120;
            item_width := 80;
            item_height := 90 + (i % 5) * 10;
        ELSIF selected_cargo_type_id = bag_type_id THEN
            item_weight := 6 + (i % 5) * 1.5;
            item_length := 70;
            item_width := 45;
            item_height := 30;
        ELSIF selected_cargo_type_id = crate_type_id THEN
            item_weight := 12 + (i % 4) * 2.2;
            item_length := 60;
            item_width := 40;
            item_height := 45;
        ELSE
            item_weight := 4 + (i % 7) * 1.35;
            item_length := 45 + (i % 3) * 10;
            item_width := 35 + (i % 2) * 5;
            item_height := 25 + (i % 4) * 5;
        END IF;

        IF NOT EXISTS (SELECT 1 FROM orders WHERE comment = order_comment) THEN
            IF handover = 'self_delivery' THEN
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
                VALUES (
                    selected_client_id,
                    selected_receiving_wh_id,
                    selected_destination_wh_id,
                    selected_product_type_id,
                    handover,
                    delivery_date,
                    time_from,
                    time_to,
                    order_status,
                    order_comment
                )
                RETURNING id INTO v_order_id;
            ELSE
                INSERT INTO orders (
                    client_id,
                    receiving_warehouse_id,
                    destination_warehouse_id,
                    product_type_id,
                    handover_type,
                    status,
                    comment
                )
                VALUES (
                    selected_client_id,
                    selected_receiving_wh_id,
                    selected_destination_wh_id,
                    selected_product_type_id,
                    handover,
                    order_status,
                    order_comment
                )
                RETURNING id INTO v_order_id;
            END IF;
        ELSE
            SELECT id INTO v_order_id FROM orders WHERE comment = order_comment LIMIT 1;
        END IF;

        INSERT INTO order_status_history (order_id, old_status, new_status, changed_by, comment)
        SELECT v_order_id, NULL, 'created', selected_client_id, order_comment || ':created'
        WHERE NOT EXISTS (
            SELECT 1 FROM order_status_history
            WHERE order_id = v_order_id AND comment = order_comment || ':created'
        );

        IF order_status <> 'created' THEN
            INSERT INTO order_status_history (order_id, old_status, new_status, changed_by, comment)
            SELECT v_order_id, 'created', order_status, selected_logist_id, order_comment || ':' || order_status
            WHERE NOT EXISTS (
                SELECT 1 FROM order_status_history
                WHERE order_id = v_order_id AND comment = order_comment || ':' || order_status
            );
        END IF;

        IF handover = 'pickup' THEN
            pickup_status := CASE
                WHEN order_status = 'waiting_pickup' THEN 'assigned'
                WHEN order_status IN ('received', 'stored', 'assigned_to_shipping', 'shipped', 'delivered') THEN 'picked_up'
                ELSE 'approved'
            END;

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
                v_order_id,
                CASE ((i - 1) % 6) + 1
                    WHEN 1 THEN 'Москва, ул. Рябиновая, 43Б, склад NordTex'
                    WHEN 2 THEN 'Казань, ул. Восстания, 100, ворота 3'
                    WHEN 3 THEN 'Санкт-Петербург, Пулковское шоссе, 42к6'
                    WHEN 4 THEN 'Москва, ул. Дорожная, 60Б'
                    WHEN 5 THEN 'Казань, ул. Техническая, 10'
                    ELSE 'Санкт-Петербург, Софийская ул., 91'
                END,
                CURRENT_DATE + ((i % 8) - 2),
                make_time(10 + (i % 4), 0, 0),
                make_time(12 + (i % 4), 30, 0),
                (SELECT full_name FROM users WHERE id = selected_client_id),
                (SELECT phone FROM users WHERE id = selected_client_id),
                pickup_status,
                selected_logist_id,
                pickup_comment
            WHERE NOT EXISTS (SELECT 1 FROM pickup_requests WHERE order_id = v_order_id);
        END IF;

        IF NOT EXISTS (
            SELECT 1 FROM order_cargo_places
            WHERE order_id = v_order_id AND comment = place_comment
        ) THEN
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
            VALUES (
                v_order_id,
                selected_cargo_type_id,
                cargo_count,
                item_weight,
                item_length,
                item_width,
                item_height,
                place_comment
            )
            RETURNING id INTO place_id;
        ELSE
            SELECT id INTO place_id
            FROM order_cargo_places
            WHERE order_id = v_order_id AND comment = place_comment
            LIMIT 1;
        END IF;

        SELECT id INTO storage_zone_id
        FROM storage_zones
        WHERE warehouse_id = selected_receiving_wh_id
        ORDER BY id
        LIMIT 1 OFFSET (i % 3);

        SELECT id INTO gate_id
        FROM gates
        WHERE warehouse_id = selected_receiving_wh_id
        ORDER BY id
        LIMIT 1 OFFSET (i % 2);

        cargo_status := CASE
            WHEN order_status IN ('waiting_pickup', 'waiting_delivery', 'received') THEN 'accepted'
            WHEN order_status = 'stored' THEN 'stored'
            WHEN order_status = 'assigned_to_shipping' THEN 'ready_to_ship'
            WHEN order_status IN ('shipped', 'delivered') THEN 'shipped'
            ELSE 'accepted'
        END;

        FOR j IN 1..cargo_count LOOP
            qr := format(
                'QR-TPRO-DEMO-%s-%s',
                lpad(i::TEXT, 3, '0'),
                lpad(j::TEXT, 2, '0')
            );
            cargo_comment := format(
                'seed:v2:order:%s:cargo:%s',
                lpad(i::TEXT, 2, '0'),
                lpad(j::TEXT, 2, '0')
            );

            IF NOT EXISTS (SELECT 1 FROM cargo_items WHERE qr_code = qr) THEN
                INSERT INTO cargo_items (
                    order_id,
                    order_cargo_place_id,
                    cargo_place_type_id,
                    qr_code,
                    status,
                    storage_zone_id,
                    gate_id,
                    received_by,
                    shipped_by,
                    received_at,
                    shipped_at,
                    comment
                )
                VALUES (
                    v_order_id,
                    place_id,
                    selected_cargo_type_id,
                    qr,
                    cargo_status,
                    CASE WHEN cargo_status IN ('stored', 'ready_to_ship', 'shipped') THEN storage_zone_id ELSE NULL END,
                    CASE WHEN cargo_status IN ('ready_to_ship', 'shipped') THEN gate_id ELSE NULL END,
                    CASE WHEN cargo_status IN ('accepted', 'stored', 'ready_to_ship', 'shipped') THEN selected_worker_id ELSE NULL END,
                    CASE WHEN cargo_status = 'shipped' THEN selected_worker_id ELSE NULL END,
                    CASE WHEN cargo_status IN ('accepted', 'stored', 'ready_to_ship', 'shipped') THEN NOW() - ((30 + i + j)::TEXT || ' hours')::INTERVAL ELSE NULL END,
                    CASE WHEN cargo_status = 'shipped' THEN NOW() - ((6 + j)::TEXT || ' hours')::INTERVAL ELSE NULL END,
                    cargo_comment
                )
                RETURNING id INTO v_cargo_item_id;
            ELSE
                SELECT id INTO v_cargo_item_id FROM cargo_items WHERE qr_code = qr LIMIT 1;
            END IF;

            INSERT INTO cargo_status_history (cargo_item_id, old_status, new_status, changed_by, comment)
            SELECT v_cargo_item_id, NULL, 'accepted', selected_worker_id, cargo_comment || ':accepted'
            WHERE NOT EXISTS (
                SELECT 1 FROM cargo_status_history
                WHERE cargo_item_id = v_cargo_item_id AND comment = cargo_comment || ':accepted'
            );

            IF cargo_status <> 'accepted' THEN
                INSERT INTO cargo_status_history (cargo_item_id, old_status, new_status, changed_by, comment)
                SELECT v_cargo_item_id, 'accepted', cargo_status, selected_worker_id, cargo_comment || ':' || cargo_status
                WHERE NOT EXISTS (
                    SELECT 1 FROM cargo_status_history
                    WHERE cargo_item_id = v_cargo_item_id AND comment = cargo_comment || ':' || cargo_status
                );
            END IF;
        END LOOP;
    END LOOP;

    -- ========================================================
    -- Shipments for ready/shipped cargo.
    -- ========================================================
    FOR shipment_rec IN
        SELECT
            o.destination_warehouse_id,
            ci.gate_id,
            MIN(o.receiving_warehouse_id) AS receiving_warehouse_id,
            COUNT(*) AS cargo_count,
            BOOL_OR(ci.status = 'shipped') AS has_shipped
        FROM cargo_items ci
        JOIN orders o ON o.id = ci.order_id
        WHERE ci.comment LIKE 'seed:v2:%'
          AND ci.gate_id IS NOT NULL
          AND ci.status IN ('ready_to_ship', 'shipped')
        GROUP BY o.destination_warehouse_id, ci.gate_id
    LOOP
        IF NOT EXISTS (
            SELECT 1
            FROM shipments s
            WHERE s.destination_warehouse_id = shipment_rec.destination_warehouse_id
              AND s.gate_id = shipment_rec.gate_id
              AND s.planned_departure_at = date_trunc('day', CURRENT_TIMESTAMP) + INTERVAL '18 hours'
        ) THEN
            INSERT INTO shipments (
                destination_warehouse_id,
                gate_id,
                planned_departure_at,
                actual_departure_at,
                status,
                created_by
            )
            VALUES (
                shipment_rec.destination_warehouse_id,
                shipment_rec.gate_id,
                date_trunc('day', CURRENT_TIMESTAMP) + INTERVAL '18 hours',
                CASE WHEN shipment_rec.has_shipped THEN date_trunc('day', CURRENT_TIMESTAMP) + INTERVAL '19 hours' ELSE NULL END,
                CASE WHEN shipment_rec.has_shipped THEN 'shipped' ELSE 'loading' END,
                logist_msk_id
            )
            RETURNING id INTO shipment_id;
        ELSE
            SELECT id INTO shipment_id
            FROM shipments s
            WHERE s.destination_warehouse_id = shipment_rec.destination_warehouse_id
              AND s.gate_id = shipment_rec.gate_id
              AND s.planned_departure_at = date_trunc('day', CURRENT_TIMESTAMP) + INTERVAL '18 hours'
            LIMIT 1;
        END IF;

        INSERT INTO shipment_items (shipment_id, cargo_item_id)
        SELECT shipment_id, ci.id
        FROM cargo_items ci
        JOIN orders o ON o.id = ci.order_id
        WHERE ci.comment LIKE 'seed:v2:%'
          AND ci.gate_id = shipment_rec.gate_id
          AND o.destination_warehouse_id = shipment_rec.destination_warehouse_id
          AND ci.status IN ('ready_to_ship', 'shipped')
        ON CONFLICT DO NOTHING;
    END LOOP;
END $$;

COMMIT;

-- Быстрая проверка после запуска:
-- SELECT COUNT(*) AS warehouses FROM warehouses;
-- SELECT w.name, COUNT(DISTINCT z.id) AS storage_zones, COUNT(DISTINCT g.id) AS gates
-- FROM warehouses w
-- LEFT JOIN storage_zones z ON z.warehouse_id = w.id
-- LEFT JOIN gates g ON g.warehouse_id = w.id
-- GROUP BY w.id, w.name
-- ORDER BY w.name;
--
-- SELECT o.id, o.comment, COALESCE(SUM(ocp.quantity), 0) AS declared_places
-- FROM orders o
-- LEFT JOIN order_cargo_places ocp ON ocp.order_id = o.id
-- WHERE o.comment LIKE 'seed:v2:%'
-- GROUP BY o.id, o.comment
-- ORDER BY o.id;
