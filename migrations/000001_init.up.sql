BEGIN;

-- =========================================================
-- 1. USERS
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


-- =========================================================
-- 2. JWT TOKENS
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


-- =========================================================
-- 3. WAREHOUSES
-- =========================================================

CREATE TABLE warehouses (
                            id BIGSERIAL PRIMARY KEY,

                            name VARCHAR(255) NOT NULL,

                            warehouse_type VARCHAR(50) NOT NULL,
                            marketplace VARCHAR(100),

                            city VARCHAR(100) NOT NULL,
                            address TEXT NOT NULL,

                            is_active BOOLEAN NOT NULL DEFAULT TRUE,

                            created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

                            CONSTRAINT chk_warehouses_type
                                CHECK (warehouse_type IN ('receiving', 'destination', 'both')),

                            CONSTRAINT chk_warehouses_name_not_empty
                                CHECK (length(trim(name)) > 0),

                            CONSTRAINT chk_warehouses_name_format
                                CHECK (
                                    name ~ '^[A-Za-zА-Яа-яЁё0-9 .,''"№#()/_-]+$'
                                    ),

                            CONSTRAINT chk_warehouses_marketplace_format
                                CHECK (
                                    marketplace IS NULL
                                        OR marketplace ~ '^[A-Za-zА-Яа-яЁё0-9 ._/-]+$'
                                    ),

                            CONSTRAINT chk_warehouses_city_not_empty
                                CHECK (length(trim(city)) > 0),

                            CONSTRAINT chk_warehouses_city_format
                                CHECK (
                                    city ~ '^[A-Za-zА-Яа-яЁё]+([ -][A-Za-zА-Яа-яЁё]+)*$'
                                    ),

                            CONSTRAINT chk_warehouses_address_not_empty
                                CHECK (length(trim(address)) > 0),

                            CONSTRAINT chk_warehouses_address_format
                                CHECK (
                                    length(trim(address)) BETWEEN 5 AND 500
                                        AND address !~ '[[:cntrl:]]'
                                    )
);


-- =========================================================
-- 4. WORKER PROFILES
-- =========================================================
-- Подроли сотрудника убраны.
-- Сотрудник определяется через users.role = 'worker',
-- а его склад хранится в worker_profiles.warehouse_id.

CREATE TABLE worker_profiles (
                                 id BIGSERIAL PRIMARY KEY,

                                 user_id BIGINT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
                                 warehouse_id BIGINT NOT NULL REFERENCES warehouses(id) ON DELETE RESTRICT,

                                 created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- =========================================================
-- 5. STORAGE ZONES
-- =========================================================

CREATE TABLE storage_zones (
                               id BIGSERIAL PRIMARY KEY,

                               warehouse_id BIGINT NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,

                               name VARCHAR(100) NOT NULL,
                               description TEXT,

                               is_active BOOLEAN NOT NULL DEFAULT TRUE,

                               CONSTRAINT chk_storage_zones_name_not_empty
                                   CHECK (length(trim(name)) > 0),

                               CONSTRAINT chk_storage_zones_name_format
                                   CHECK (
                                       name ~ '^[A-Za-zА-Яа-яЁё0-9 ._/-]+$'
                                       ),

                               CONSTRAINT uq_storage_zones_warehouse_name
                                   UNIQUE (warehouse_id, name)
);


-- =========================================================
-- 6. GATES
-- =========================================================

CREATE TABLE gates (
                       id BIGSERIAL PRIMARY KEY,

                       warehouse_id BIGINT NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,

                       name VARCHAR(100) NOT NULL,
                       is_active BOOLEAN NOT NULL DEFAULT TRUE,

                       CONSTRAINT chk_gates_name_not_empty
                           CHECK (length(trim(name)) > 0),

                       CONSTRAINT chk_gates_name_format
                           CHECK (
                               name ~ '^[A-Za-zА-Яа-яЁё0-9 ._/-]+$'
                               ),

                       CONSTRAINT uq_gates_warehouse_name
                           UNIQUE (warehouse_id, name)
);


-- =========================================================
-- 7. PRODUCT TYPES
-- =========================================================
-- Что внутри груза: одежда, косметика, электроника и т.д.

CREATE TABLE product_types (
                               id BIGSERIAL PRIMARY KEY,

                               name VARCHAR(150) NOT NULL UNIQUE,
                               description TEXT,

                               is_active BOOLEAN NOT NULL DEFAULT TRUE,

                               CONSTRAINT chk_product_types_name_not_empty
                                   CHECK (length(trim(name)) > 0),

                               CONSTRAINT chk_product_types_name_format
                                   CHECK (
                                       name ~ '^[A-Za-zА-Яа-яЁё0-9 ._/-]+$'
                                       )
);


-- =========================================================
-- 8. CARGO PLACE TYPES
-- =========================================================
-- Во что упакован груз: паллета, коробка, мешок и т.д.

CREATE TABLE cargo_place_types (
                                   id BIGSERIAL PRIMARY KEY,

                                   name VARCHAR(100) NOT NULL UNIQUE,
                                   description TEXT,

                                   is_active BOOLEAN NOT NULL DEFAULT TRUE,

                                   CONSTRAINT chk_cargo_place_types_name_not_empty
                                       CHECK (length(trim(name)) > 0),

                                   CONSTRAINT chk_cargo_place_types_name_format
                                       CHECK (
                                           name ~ '^[A-Za-zА-Яа-яЁё0-9 ._/-]+$'
                                           )
);


-- =========================================================
-- 9. ORDERS
-- =========================================================

CREATE TABLE orders (
                        id BIGSERIAL PRIMARY KEY,

                        client_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,

                        receiving_warehouse_id BIGINT NOT NULL REFERENCES warehouses(id) ON DELETE RESTRICT,
                        destination_warehouse_id BIGINT NOT NULL REFERENCES warehouses(id) ON DELETE RESTRICT,

                        product_type_id BIGINT NOT NULL REFERENCES product_types(id) ON DELETE RESTRICT,

                        handover_type VARCHAR(50) NOT NULL,

                        status VARCHAR(50) NOT NULL DEFAULT 'created',

                        comment TEXT,

                        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

                        CONSTRAINT chk_orders_handover_type
                            CHECK (handover_type IN ('pickup', 'self_delivery')),

                        CONSTRAINT chk_orders_status
                            CHECK (
                                status IN (
                                           'created',
                                           'waiting_pickup',
                                           'waiting_delivery',
                                           'received',
                                           'stored',
                                           'assigned_to_shipping',
                                           'shipped',
                                           'delivered',
                                           'cancelled'
                                    )
                                ),

                        CONSTRAINT chk_orders_status_matches_handover
                            CHECK (
                                (handover_type = 'pickup' AND status <> 'waiting_delivery')
                                    OR
                                (handover_type = 'self_delivery' AND status <> 'waiting_pickup')
                                )
);


-- =========================================================
-- 10. ORDER CARGO PLACES
-- =========================================================
-- Заявленный клиентом состав груза:
-- например, 2 паллеты, 5 коробок, 3 мешка.

CREATE TABLE order_cargo_places (
                                    id BIGSERIAL PRIMARY KEY,

                                    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
                                    cargo_place_type_id BIGINT NOT NULL REFERENCES cargo_place_types(id) ON DELETE RESTRICT,

                                    quantity INT NOT NULL,

                                    weight_per_place_kg NUMERIC(10, 2),
                                    length_cm NUMERIC(10, 2),
                                    width_cm NUMERIC(10, 2),
                                    height_cm NUMERIC(10, 2),

                                    comment TEXT,

                                    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

                                    CONSTRAINT chk_order_cargo_places_quantity
                                        CHECK (quantity > 0),

                                    CONSTRAINT chk_order_cargo_places_weight
                                        CHECK (weight_per_place_kg IS NULL OR weight_per_place_kg > 0),

                                    CONSTRAINT chk_order_cargo_places_length
                                        CHECK (length_cm IS NULL OR length_cm > 0),

                                    CONSTRAINT chk_order_cargo_places_width
                                        CHECK (width_cm IS NULL OR width_cm > 0),

                                    CONSTRAINT chk_order_cargo_places_height
                                        CHECK (height_cm IS NULL OR height_cm > 0),

                                    CONSTRAINT uq_order_cargo_places_id_order_id
                                        UNIQUE (id, order_id),

                                    CONSTRAINT uq_order_cargo_places_id_type_id
                                        UNIQUE (id, cargo_place_type_id)
);


-- =========================================================
-- 11. PICKUP REQUESTS
-- =========================================================
-- Заполняется только если orders.handover_type = 'pickup'.
-- Проверку соответствия handover_type лучше делать в backend
-- или отдельным триггером, так как CHECK не может ссылаться
-- на другую таблицу.

CREATE TABLE pickup_requests (
                                 id BIGSERIAL PRIMARY KEY,

                                 order_id BIGINT NOT NULL UNIQUE REFERENCES orders(id) ON DELETE CASCADE,

                                 pickup_address TEXT NOT NULL,
                                 pickup_date DATE NOT NULL,
                                 pickup_time_from TIME,
                                 pickup_time_to TIME,

                                 contact_name VARCHAR(255),
                                 contact_phone VARCHAR(50),

                                 status VARCHAR(50) NOT NULL DEFAULT 'created',

                                 assigned_logist_id BIGINT REFERENCES users(id) ON DELETE SET NULL,

                                 comment TEXT,

                                 created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                 updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

                                 CONSTRAINT chk_pickup_requests_status
                                     CHECK (
                                         status IN (
                                                    'created',
                                                    'approved',
                                                    'assigned',
                                                    'in_progress',
                                                    'picked_up',
                                                    'cancelled'
                                             )
                                         ),

                                 CONSTRAINT chk_pickup_requests_address_not_empty
                                     CHECK (length(trim(pickup_address)) > 0),

                                 CONSTRAINT chk_pickup_requests_address_format
                                     CHECK (
                                         length(trim(pickup_address)) BETWEEN 5 AND 500
                                             AND pickup_address !~ '[[:cntrl:]]'
                                         ),

                                 CONSTRAINT chk_pickup_requests_contact_name_format
                                     CHECK (
                                         contact_name IS NULL
                                             OR contact_name ~ '^[A-Za-zА-Яа-яЁё]+([ .''-][A-Za-zА-Яа-яЁё]+)*$'
                                         ),

                                 CONSTRAINT chk_pickup_requests_contact_phone_format
                                     CHECK (
                                         contact_phone IS NULL
                                             OR contact_phone ~ '^[+]?[0-9]{10,15}$'
                                         ),

                                 CONSTRAINT chk_pickup_requests_time_range
                                     CHECK (
                                         (pickup_time_from IS NULL AND pickup_time_to IS NULL)
                                             OR
                                         (
                                             pickup_time_from IS NOT NULL
                                                 AND pickup_time_to IS NOT NULL
                                                 AND pickup_time_to > pickup_time_from
                                             )
                                         )
);


-- =========================================================
-- 12. CARGO ITEMS
-- =========================================================
-- Фактически принятые грузовые места.
-- Каждая запись = одно физическое место с уникальным QR-кодом.

CREATE TABLE cargo_items (
                             id BIGSERIAL PRIMARY KEY,

                             order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,

                             order_cargo_place_id BIGINT NOT NULL,
                             cargo_place_type_id BIGINT NOT NULL REFERENCES cargo_place_types(id) ON DELETE RESTRICT,

                             qr_code VARCHAR(255) NOT NULL UNIQUE,

                             status VARCHAR(50) NOT NULL DEFAULT 'accepted',

                             storage_zone_id BIGINT REFERENCES storage_zones(id) ON DELETE SET NULL,
                             gate_id BIGINT REFERENCES gates(id) ON DELETE SET NULL,

                             received_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
                             shipped_by BIGINT REFERENCES users(id) ON DELETE SET NULL,

                             received_at TIMESTAMPTZ,
                             shipped_at TIMESTAMPTZ,

                             comment TEXT,

                             created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
                             updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

                             CONSTRAINT fk_cargo_items_declared_place_order
                                 FOREIGN KEY (order_cargo_place_id, order_id)
                                     REFERENCES order_cargo_places(id, order_id)
                                     ON DELETE RESTRICT,

                             CONSTRAINT fk_cargo_items_declared_place_type
                                 FOREIGN KEY (order_cargo_place_id, cargo_place_type_id)
                                     REFERENCES order_cargo_places(id, cargo_place_type_id)
                                     ON DELETE RESTRICT,

                             CONSTRAINT chk_cargo_items_status
                                 CHECK (
                                     status IN (
                                                'accepted',
                                                'stored',
                                                'ready_to_ship',
                                                'shipped',
                                                'lost',
                                                'damaged',
                                                'cancelled'
                                         )
                                     ),

                             CONSTRAINT chk_cargo_items_qr_code_not_empty
                                 CHECK (length(trim(qr_code)) > 0),

                             CONSTRAINT chk_cargo_items_qr_code_format
                                 CHECK (
                                     length(trim(qr_code)) BETWEEN 3 AND 255
                                         AND qr_code ~ '^[-A-Za-z0-9._:/]+$'
                                     ),

                             CONSTRAINT chk_cargo_items_shipped_data
                                 CHECK (
                                     status <> 'shipped'
                                         OR
                                     shipped_at IS NOT NULL
                                     )
);


-- =========================================================
-- 13. CARGO RECEIPT DISCREPANCIES
-- =========================================================
-- Расхождения при приёмке:
-- заявлено 5 коробок, фактически принято 4.

CREATE TABLE cargo_receipt_discrepancies (
                                             id BIGSERIAL PRIMARY KEY,

                                             order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,

                                             order_cargo_place_id BIGINT NOT NULL,
                                             cargo_place_type_id BIGINT NOT NULL REFERENCES cargo_place_types(id) ON DELETE RESTRICT,

                                             declared_quantity INT NOT NULL DEFAULT 0,
                                             actual_quantity INT NOT NULL DEFAULT 0,

                                             comment TEXT,
                                             created_by BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,

                                             created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

                                             CONSTRAINT fk_discrepancies_declared_place_order
                                                 FOREIGN KEY (order_cargo_place_id, order_id)
                                                     REFERENCES order_cargo_places(id, order_id)
                                                     ON DELETE RESTRICT,

                                             CONSTRAINT fk_discrepancies_declared_place_type
                                                 FOREIGN KEY (order_cargo_place_id, cargo_place_type_id)
                                                     REFERENCES order_cargo_places(id, cargo_place_type_id)
                                                     ON DELETE RESTRICT,

                                             CONSTRAINT chk_discrepancies_quantities
                                                 CHECK (
                                                     declared_quantity >= 0
                                                         AND actual_quantity >= 0
                                                         AND declared_quantity <> actual_quantity
                                                     )
);


-- =========================================================
-- 14. PICKUP CALENDAR BLOCKS
-- =========================================================
-- Закрытые даты для приёма/вывоза.

CREATE TABLE pickup_calendar_blocks (
                                        id BIGSERIAL PRIMARY KEY,

                                        warehouse_id BIGINT NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,

                                        blocked_date DATE NOT NULL,
                                        reason TEXT,

                                        created_by BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
                                        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

                                        CONSTRAINT uq_pickup_calendar_blocks_warehouse_date
                                            UNIQUE (warehouse_id, blocked_date)
);


-- =========================================================
-- 15. PICKUP CAPACITY
-- =========================================================
-- Лимиты заявок на дату.

CREATE TABLE pickup_capacity (
                                 id BIGSERIAL PRIMARY KEY,

                                 warehouse_id BIGINT NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,

                                 pickup_date DATE NOT NULL,

                                 max_orders INT NOT NULL DEFAULT 0,
                                 current_orders INT NOT NULL DEFAULT 0,
                                 is_closed BOOLEAN NOT NULL DEFAULT FALSE,

                                 CONSTRAINT uq_pickup_capacity_warehouse_date
                                     UNIQUE (warehouse_id, pickup_date),

                                 CONSTRAINT chk_pickup_capacity_max_orders
                                     CHECK (max_orders >= 0),

                                 CONSTRAINT chk_pickup_capacity_current_orders
                                     CHECK (current_orders >= 0),

                                 CONSTRAINT chk_pickup_capacity_current_not_greater_max
                                     CHECK (current_orders <= max_orders)
);


-- =========================================================
-- 16. SHIPMENTS
-- =========================================================
-- Отгрузки на склад назначения.

CREATE TABLE shipments (
                           id BIGSERIAL PRIMARY KEY,

                           destination_warehouse_id BIGINT NOT NULL REFERENCES warehouses(id) ON DELETE RESTRICT,
                           gate_id BIGINT NOT NULL REFERENCES gates(id) ON DELETE RESTRICT,

                           planned_departure_at TIMESTAMPTZ NOT NULL,
                           actual_departure_at TIMESTAMPTZ,

                           status VARCHAR(50) NOT NULL DEFAULT 'planned',

                           created_by BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,

                           created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

                           CONSTRAINT chk_shipments_status
                               CHECK (
                                   status IN (
                                              'planned',
                                              'loading',
                                              'shipped',
                                              'completed',
                                              'cancelled'
                                       )
                                   ),

                           CONSTRAINT chk_shipments_actual_departure
                               CHECK (
                                   actual_departure_at IS NULL
                                       OR actual_departure_at >= planned_departure_at
                                   )
);


-- =========================================================
-- 17. SHIPMENT ITEMS
-- =========================================================
-- Связь отгрузки с конкретными QR-грузовыми местами.

CREATE TABLE shipment_items (
                                id BIGSERIAL PRIMARY KEY,

                                shipment_id BIGINT NOT NULL REFERENCES shipments(id) ON DELETE CASCADE,
                                cargo_item_id BIGINT NOT NULL REFERENCES cargo_items(id) ON DELETE RESTRICT,

                                CONSTRAINT uq_shipment_items_shipment_cargo
                                    UNIQUE (shipment_id, cargo_item_id),

                                CONSTRAINT uq_shipment_items_cargo_item
                                    UNIQUE (cargo_item_id)
);


-- =========================================================
-- 18. ORDER STATUS HISTORY
-- =========================================================

CREATE TABLE order_status_history (
                                      id BIGSERIAL PRIMARY KEY,

                                      order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,

                                      old_status VARCHAR(50),
                                      new_status VARCHAR(50) NOT NULL,

                                      changed_by BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
                                      comment TEXT,

                                      changed_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

                                      CONSTRAINT chk_order_status_history_old_status
                                          CHECK (
                                              old_status IS NULL
                                                  OR old_status IN (
                                                                    'created',
                                                                    'waiting_pickup',
                                                                    'waiting_delivery',
                                                                    'received',
                                                                    'stored',
                                                                    'assigned_to_shipping',
                                                                    'shipped',
                                                                    'delivered',
                                                                    'cancelled'
                                                  )
                                              ),

                                      CONSTRAINT chk_order_status_history_new_status
                                          CHECK (
                                              new_status IN (
                                                             'created',
                                                             'waiting_pickup',
                                                             'waiting_delivery',
                                                             'received',
                                                             'stored',
                                                             'assigned_to_shipping',
                                                             'shipped',
                                                             'delivered',
                                                             'cancelled'
                                                  )
                                              )
);


-- =========================================================
-- 19. CARGO STATUS HISTORY
-- =========================================================

CREATE TABLE cargo_status_history (
                                      id BIGSERIAL PRIMARY KEY,

                                      cargo_item_id BIGINT NOT NULL REFERENCES cargo_items(id) ON DELETE CASCADE,

                                      old_status VARCHAR(50),
                                      new_status VARCHAR(50) NOT NULL,

                                      changed_by BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
                                      comment TEXT,

                                      changed_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

                                      CONSTRAINT chk_cargo_status_history_old_status
                                          CHECK (
                                              old_status IS NULL
                                                  OR old_status IN (
                                                                    'accepted',
                                                                    'stored',
                                                                    'ready_to_ship',
                                                                    'shipped',
                                                                    'lost',
                                                                    'damaged',
                                                                    'cancelled'
                                                  )
                                              ),

                                      CONSTRAINT chk_cargo_status_history_new_status
                                          CHECK (
                                              new_status IN (
                                                             'accepted',
                                                             'stored',
                                                             'ready_to_ship',
                                                             'shipped',
                                                             'lost',
                                                             'damaged',
                                                             'cancelled'
                                                  )
                                              )
);


-- =========================================================
-- 20. UPDATED_AT TRIGGERS
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


-- =========================================================
-- 21. INDEXES
-- =========================================================
-- UNIQUE и PRIMARY KEY уже создают индексы автоматически.
-- Ниже — дополнительные индексы для частых JOIN, WHERE и фильтров.

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


-- =========================================================
-- 22. SEED DATA
-- =========================================================

INSERT INTO product_types (name, description) VALUES
                                                  ('Одежда', 'Одежда и текстильные изделия'),
                                                  ('Обувь', 'Обувные товары'),
                                                  ('Косметика', 'Косметические товары'),
                                                  ('Электроника', 'Электронные устройства и аксессуары'),
                                                  ('Игрушки', 'Детские товары и игрушки'),
                                                  ('Товары для дома', 'Бытовые и хозяйственные товары'),
                                                  ('Другое', 'Иной тип товара')
ON CONFLICT (name) DO NOTHING;


INSERT INTO cargo_place_types (name, description) VALUES
                                                      ('Паллета', 'Грузовое место на поддоне'),
                                                      ('Коробка', 'Товар в коробке'),
                                                      ('Мешок', 'Товар в мешке'),
                                                      ('Ящик', 'Товар в ящике'),
                                                      ('Бочка', 'Товар в бочке'),
                                                      ('Биг-бэг', 'Крупный мягкий контейнер'),
                                                      ('Другое', 'Иной тип грузового места')
ON CONFLICT (name) DO NOTHING;


COMMIT;
