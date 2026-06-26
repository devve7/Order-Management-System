ALTER TABLE orders
ADD CONSTRAINT fk_orders_customer_user
FOREIGN KEY (customer_id)
REFERENCES users(id);