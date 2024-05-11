CREATE TABLE IF NOT EXISTS user_admins (
  id uuid PRIMARY KEY,
  sid serial,
  phone_number varchar NOT NULL UNIQUE,
  name varchar NOT NULL,
  password varchar NOT NULL,
  role varchar NOT NULL,
  created_at timestamptz NOT NULL
);
