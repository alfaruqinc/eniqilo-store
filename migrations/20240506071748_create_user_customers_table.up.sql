CREATE TABLE IF NOT EXISTS user_customers (
  id uuid PRIMARY KEY,
  sid serial,
  phone_number varchar NOT NULL UNIQUE,
  name varchar NOT NULL,
  created_at timestamptz NOT NULL
);
