CREATE TABLE "products" (
  "id" int PRIMARY KEY,
  "created_at" timestamp,
  "name" varchar,
  "price" decimal
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "products" ("id", "name", "price")
VALUES
  (1, 'Widget', 9.99),
  (2, 'Gadget', 19.99);

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;