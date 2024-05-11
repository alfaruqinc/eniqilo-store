BEGIN;

CREATE TABLE IF NOT EXISTS product_checkouts (
  id uuid PRIMARY KEY,
  sequence_id SERIAL UNIQUE  NOT NULL,
  product_id uuid NOT NULL,
  quantity int NOT NULL,
  checkout_id uuid NOT NULL
);

ALTER TABLE product_checkouts ADD CONSTRAINT fk_product_id_product_checkouts FOREIGN KEY (product_id) REFERENCES products (id);

ALTER TABLE product_checkouts ADD CONSTRAINT fk_checkout_id_product_checkouts FOREIGN KEY (checkout_id) REFERENCES checkouts (id);

COMMIT;
