CREATE TABLE user_customers (
  id uuid PRIMARY KEY,
  phone_number varchar NOT NULL,
  name varchar NOT NULL,
  password varchar NOT NULL,
  created_at timestamptz NOT NULL
);