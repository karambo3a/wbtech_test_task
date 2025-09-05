CREATE TABLE IF NOT EXISTS deliveries
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(255) NOT NULL,
    zip VARCHAR(255) NOT NULL,
    city VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    region VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL CHECK (position('@' IN email) > 0),
    UNIQUE (name, phone, zip, city, address, region, email)
);

CREATE TABLE IF NOT EXISTS payments
(
    id SERIAL PRIMARY KEY,
    transaction VARCHAR(255) NOT NULL,
    request_id VARCHAR(255) NOT NULL DEFAULT '',
    currency VARCHAR(255) NOT NULL,
    provider VARCHAR(255) NOT NULL,
    amount INTEGER NOT NULL,
    payment_dt BIGINT NOT NULL,
    bank VARCHAR(255) NOT NULL,
    delivery_cost INTEGER NOT NULL,
    goods_total INTEGER NOT NULL,
    custom_fee INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS orders
(
    order_uid VARCHAR(255) PRIMARY KEY,
    track_number VARCHAR(255) NOT NULL,
    entry VARCHAR(255) NOT NULL,
    delivery_id INTEGER NOT NULL REFERENCES deliveries(id) ON DELETE CASCADE,
    payment_id INTEGER NOT NULL REFERENCES payments(id) ON DELETE CASCADE,
    locale VARCHAR(255) NOT NULL,
    internal_signature VARCHAR(255) NOT NULL DEFAULT '',
    customer_id VARCHAR(255) NOT NULL,
    delivery_service VARCHAR(255) NOT NULL,
    shardkey VARCHAR(255) NOT NULL,
    sm_id INTEGER NOT NULL,
    date_created TIMESTAMPTZ NOT NULL,
    oof_shard VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS items
(
    id SERIAL PRIMARY KEY,
    chrt_id INTEGER NOT NULL,
    track_number VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL,
    rid VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    sale INTEGER NOT NULL CHECK (sale >= 0 AND sale <= 100),
    size VARCHAR(10) NOT NULL,
    total_price INTEGER NOT NULL,
    nm_id INTEGER NOT NULL,
    brand VARCHAR(255) NOT NULL,
    status INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS orders_x_items
(
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(255) NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
    item_id INTEGER NOT NULL REFERENCES items(id) ON DELETE CASCADE
);


INSERT INTO deliveries (name, phone, zip, city, address, region, email)
VALUES ('Test Testov', '+9720000000', '2639809', 'Kiryat Mozkin', 'Ploshad Mira 15', 'Kraiot', 'test@gmail.com');

INSERT INTO payments (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
VALUES ('b563feb7b2b84b6test', '', 'USD', 'wbpay', 1817, 1637907727, 'alpha', 1500, 317, 0);

INSERT INTO orders (order_uid, track_number, entry, delivery_id, payment_id, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
VALUES ('b563feb7b2b84b6test', 'WBILMTESTTRACK', 'WBIL', 1, 1, 'en', '', 'test', 'meest', '9', 99, '2021-11-26T06:22:19Z', '1');

INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
VALUES (9934930, 'WBILMTESTTRACK', 453, 'ab4219087a764ae0btest', 'Mascaras', 30, '0', 317, 2389212, 'Vivienne Sabo', 202);

INSERT INTO orders_x_items (order_uid, item_id)
VALUES ('b563feb7b2b84b6test', 1);

INSERT INTO deliveries (name, phone, zip, city, address, region, email)
VALUES ('John Smith', '+442012345678', 'SW1A 1AA', 'London', '10 Downing Street', 'Greater London', 'john.smith@email.com');

INSERT INTO payments (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
VALUES ('order_1700000001_1234', 'req_0001', 'GBP', 'stripe', 5420, 1672534891, 'barclays', 500, 4920, 0);

INSERT INTO orders (order_uid, track_number, entry, delivery_id, payment_id, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
VALUES ('order_1700000001_1234', 'WBILTRACK000001', 'WBIL', 2, 2, 'en', 'signature_001', 'customer_001', 'dhl', '1', 42, '2023-01-15T14:28:31Z', '1');

INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
VALUES (87654321, 'WBILTRACK000001', 2460, 'rid_0001', 'Sports Sneakers', 15, '42', 2091, 6543210, 'Nike', 200);

INSERT INTO orders_x_items (order_uid, item_id)
VALUES ('order_1700000001_1234', 2);

INSERT INTO deliveries (name, phone, zip, city, address, region, email)
VALUES ('Emma Johnson', '+13125551234', '10001', 'New York', '350 5th Avenue', 'New York', 'emma.johnson@email.com');

INSERT INTO payments (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
VALUES ('order_1700000002_5678', 'req_0002', 'USD', 'paypal', 8900, 1672621291, 'chase', 300, 8600, 0);

INSERT INTO orders (order_uid, track_number, entry, delivery_id, payment_id, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
VALUES ('order_1700000002_5678', 'WBILTRACK000002', 'WBIL', 3, 3, 'en', 'signature_002', 'customer_002', 'fedex', '2', 43, '2023-01-16T10:15:22Z', '2');

INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
VALUES (87654322, 'WBILTRACK000002', 4300, 'rid_0002', 'Smartphone', 10, '0', 3870, 6543211, 'Apple', 200);

INSERT INTO orders_x_items (order_uid, item_id)
VALUES ('order_1700000002_5678', 3);

INSERT INTO deliveries (name, phone, zip, city, address, region, email)
VALUES ('Michael Brown', '+61391234567', '2000', 'Sydney', '1 Macquarie Street', 'NSW', 'michael.brown@email.com');

INSERT INTO payments (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
VALUES ('order_1700000003_9012', 'req_0003', 'AUD', 'afterpay', 12500, 1672707691, 'commonwealth', 700, 11800, 0);

INSERT INTO orders (order_uid, track_number, entry, delivery_id, payment_id, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
VALUES ('order_1700000003_9012', 'WBILTRACK000003', 'WBIL', 4, 4, 'en', 'signature_003', 'customer_003', 'australia post', '3', 44, '2023-01-17T16:45:18Z', '3');

INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
VALUES (87654323, 'WBILTRACK000003', 5900, 'rid_0003', 'Laptop', 5, '15.6', 5605, 6543212, 'Dell', 200);

INSERT INTO orders_x_items (order_uid, item_id)
VALUES ('order_1700000003_9012', 4);

INSERT INTO deliveries (name, phone, zip, city, address, region, email)
VALUES ('Sarah Wilson', '+498912345678', '10115', 'Berlin', 'Unter den Linden 77', 'Berlin', 'sarah.wilson@email.com');

INSERT INTO payments (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
VALUES ('order_1700000004_3456', 'req_0004', 'EUR', 'klarna', 7800, 1672794091, 'deutsche', 400, 7400, 0);

INSERT INTO orders (order_uid, track_number, entry, delivery_id, payment_id, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
VALUES ('order_1700000004_3456', 'WBILTRACK000004', 'WBIL', 5, 5, 'en', 'signature_004', 'customer_004', 'dhl', '4', 45, '2023-01-18T09:30:45Z', '4');

INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
VALUES (87654324, 'WBILTRACK000004', 3700, 'rid_0004', 'Tablet', 20, '10.1', 2960, 6543213, 'Samsung', 200);

INSERT INTO orders_x_items (order_uid, item_id)
VALUES ('order_1700000004_3456', 5);

INSERT INTO deliveries (name, phone, zip, city, address, region, email)
VALUES ('David Taylor', '+14165551234', 'M5V 2T6', 'Toronto', '1 Dundas Street West', 'Ontario', 'david.taylor@email.com');

INSERT INTO payments (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
VALUES ('order_1700000005_7890', 'req_0005', 'CAD', 'shopify', 15600, 1672880491, 'royal bank', 600, 15000, 0);

INSERT INTO orders (order_uid, track_number, entry, delivery_id, payment_id, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
VALUES ('order_1700000005_7890', 'WBILTRACK000005', 'WBIL', 6, 6, 'en', 'signature_005', 'customer_005', 'canada post', '5', 46, '2023-01-19T14:20:33Z', '5');

INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
VALUES (87654325, 'WBILTRACK000005', 7500, 'rid_0005', 'Digital Camera', 25, '0', 5625, 6543214, 'Canon', 200);

INSERT INTO orders_x_items (order_uid, item_id)
VALUES ('order_1700000005_7890', 6);
