CREATE DATABASE IF NOT EXISTS shop;
CREATE TABLE shop.customers (
  id INT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL
);
CREATE TABLE shop.orders (
  id INT NOT NULL,
  customer_id INT,
  placed DATE,
  total DOUBLE PRECISION DEFAULT -1,
  PRIMARY KEY (id),
  FOREIGN KEY (customer_id) REFERENCES shop.customers(id) ON DELETE SET NULL
);
CREATE TABLE audit.events (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  order_id INT NOT NULL REFERENCES shop.orders (id),
  payload JSON,
  CHECK (JSON_VALID(payload))
);
ALTER TABLE shop.orders ADD CONSTRAINT fk_self FOREIGN KEY (id) REFERENCES audit.events (order_id);
CREATE UNIQUE INDEX ux_name ON shop.customers (name(20) DESC);
CREATE INDEX ix_placed USING HASH ON shop.orders (placed, total);
