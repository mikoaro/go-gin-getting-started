DROP TABLE IF EXISTS products;
CREATE TABLE products (
  id            SERIAL PRIMARY KEY,
  name          VARCHAR(128) NOT NULL,
  description   VARCHAR(255) NOT NULL,
  image         VARCHAR(128) NOT NULL,
  category      VARCHAR(128) NOT NULL,
  price         DECIMAL(5,2) NOT NULL
);

INSERT INTO products
  (name, description, image, category, price)
VALUES
  ('ELECTRA X2', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit', 'https://s3-us-west-2.amazonaws.com/dev-or-devrl-s3-bucket/sample-apps/ebikes/electrax2.jpg', 'Mountain', 56.99),
  ('ELECTRA X3', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit', 'https://s3-us-west-2.amazonaws.com/dev-or-devrl-s3-bucket/sample-apps/ebikes/electrax3.jpg', 'Mountain', 63.99),
  ('ELECTRA X1', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit', 'https://s3-us-west-2.amazonaws.com/dev-or-devrl-s3-bucket/sample-apps/ebikes/electrax1.jpg', 'Mountain', 34.98);

