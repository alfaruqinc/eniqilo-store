CREATE TABLE IF NOT EXISTS checkouts (
  id uuid PRIMARY KEY,
  sid serial,
  user_customer_id uuid NOT NULL,
  paid int NOT NULL,
  change int NOT NULL,
  created_at timestamptz
);

ALTER TABLE checkouts ADD CONSTRAINT fk_user_customer_id_checkouts FOREIGN KEY (user_customer_id) REFERENCES user_customers (id);
