CREATE TABLE "data" (
  "price" decimal(10,2),
  "rate" float,
  "amount" numeric
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "data" ("price", "rate", "amount")
VALUES
  (99.99, 3.14159, 0.001),
  (50.5, 0.5, 100);

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;