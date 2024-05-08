CREATE TABLE IF NOT EXISTS user_customers (
  id uuid PRIMARY KEY,
  phone_number varchar NOT NULL UNIQUE,
  name varchar NOT NULL,
  created_at timestamptz NOT NULL
);
