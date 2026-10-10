CREATE TABLE "order_items" (
  "order_id" int,
  "product_id" int,
  "quantity" int,
  PRIMARY KEY ("order_id", "product_id")
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "order_items" ("order_id", "product_id", "quantity")
VALUES
  (1, 100, 2),
  (1, 100, 5);

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;