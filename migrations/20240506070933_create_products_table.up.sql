CREATE TABLE products (
  id uuid PRIMARY KEY,
  name varchar NOT NULL,
  sku varchar NOT NULL,
  category varchar NOT NULL,
  image_url varchar NOT NULL,
  notes varchar NOT NULL,
  price numeric NOT NULL,
  stock int NOT NULL,
  location varchar NOT NULL,
  is_available bool NOT NULL,
  created_at timestamptz NOT NULL
);
