CREATE TABLE checkouts (
  id uuid PRIMARY KEY,
  user_customer_id uuid NOT NULL,
  product_id uuid NOT NULL,
  paid int NOT NULL,
  change int NOT NULL,
  created_at timestamptz
);

ALTER TABLE checkouts ADD CONSTRAINT fk_user_customer_id_checkouts FOREIGN KEY (user_customer_id) REFERENCES user_customers (id);

ALTER TABLE checkouts ADD CONSTRAINT fk_product_id_checkouts FOREIGN KEY (product_id) REFERENCES products (id);