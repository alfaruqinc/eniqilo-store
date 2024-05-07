CREATE TABLE IF NOT EXISTS user_admins (
  id uuid PRIMARY KEY,
  phone_number varchar NOT NULL,
  name varchar NOT NULL,
  password varchar NOT NULL,
  role varchar NOT NULL,
  created_at timestamptz NOT NULL
);
