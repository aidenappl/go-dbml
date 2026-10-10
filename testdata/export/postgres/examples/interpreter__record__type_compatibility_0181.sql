CREATE TABLE "products" (
  "id" int,
  "price" decimal(5,2),
  "rate" numeric(10,3)
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "products" ("id", "price", "rate")
VALUES
  (1, 99.99, 1.234),
  (2, 12345.67, 12345.678),
  (3, 99.999, 1234567.89);

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;