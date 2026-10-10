CREATE TABLE "data" (
  "id" int,
  "quantity" bigint,
  "serial_num" smallint,
  "price" decimal(10,2),
  "rate" float,
  "amount" double
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "data" ("id", "quantity", "serial_num", "price", "rate", "amount")
VALUES
  (1, 1000, 5, 99.99, 3.14159, 101325.5),
  (2, -500, -10, -50, -2.5, -40),
  (3, 0, 0, 0, 0, 0.001);

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;