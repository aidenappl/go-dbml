CREATE TABLE "products" (
  "id" "NUMBER" PRIMARY KEY,
  "name" "VARCHAR2(100)" NOT NULL,
  "category" "VARCHAR2(50)",
  "price" "NUMBER(10, 2)"
);

CREATE TABLE "suppliers" (
  "id" "NUMBER" PRIMARY KEY,
  "name" "VARCHAR2(100)" NOT NULL,
  "country" "VARCHAR2(50)"
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "products" ("id", "name", "category", "price")
VALUES
  (1, 'Laptop', 'Electronics', 999.99),
  (2, 'Mouse', 'Electronics', 29.99),
  (3, 'Keyboard', 'Electronics', 79.99),
  (4, 'Monitor', 'Electronics', 299.99),
  (5, 'Desk', NULL, NULL),
  (6, 'Chair', NULL, 149.99),
  (7, UPPER('tablet'), NULL, 100 * 5),
  (8, 'Headphones', NULL, 50.00 + 49.99);
INSERT INTO "suppliers" ("id", "name", "country")
VALUES
  (1, 'TechCorp', 'USA'),
  (2, 'GlobalTech', 'Germany'),
  (3, 'AsiaSupply', 'China');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;