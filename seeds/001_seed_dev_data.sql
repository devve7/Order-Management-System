-- Products

INSERT INTO products (id, name, price_cents, stock, active)
VALUES
    (1, 'Laptop', 120000, 10, TRUE),
    (2, 'Phone', 70000, 25, TRUE),
    (3, 'Headphones', 15000, 40, TRUE),
    (4, 'Keyboard', 8000, 30, TRUE),
    (5, 'Mouse', 5000, 50, TRUE),
    (6, 'Inactive Monitor', 30000, 0, FALSE)
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    price_cents = EXCLUDED.price_cents,
    stock = EXCLUDED.stock,
    active = EXCLUDED.active;

SELECT setval('products_id_seq', (SELECT MAX(id) FROM products));

-- Orders

INSERT INTO orders (id, customer_id, status, created_at, next_item_id, version)
VALUES
    (1, 101, 'created', NOW() - INTERVAL '2 days', 3, 1),
    (2, 102, 'paid', NOW() - INTERVAL '1 day', 2, 1),
    (3, 103, 'cancelled', NOW(), 2, 1)
ON CONFLICT (id) DO UPDATE SET
    customer_id = EXCLUDED.customer_id,
    status = EXCLUDED.status,
    created_at = EXCLUDED.created_at,
    next_item_id = EXCLUDED.next_item_id,
    version = EXCLUDED.version;

SELECT setval('orders_id_seq', (SELECT MAX(id) FROM orders));

-- Order Items

INSERT INTO order_items (order_id, item_id, product_id, name, price_cents, quantity)
VALUES
    (1, 1, 1, 'Laptop', 120000, 1),
    (1, 2, 3, 'Headphones', 15000, 2),
    (2, 1, 2, 'Phone', 70000, 1),
    (3, 1, 5, 'Mouse', 5000, 1)
ON CONFLICT (order_id, item_id) DO UPDATE SET
    product_id = EXCLUDED.product_id,
    name = EXCLUDED.name,
    price_cents = EXCLUDED.price_cents,
    quantity = EXCLUDED.quantity;