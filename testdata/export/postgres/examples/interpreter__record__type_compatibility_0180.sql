CREATE TABLE "data" (
  "id" int,
  "quantity" int,
  "price" decimal
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "data" ("id", "quantity", "price")
VALUES
  (not_a_number, 100, 100),
  (2, 10, 100),
  (3, 20, also_not_a_number),
  (4, 30, true);

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;