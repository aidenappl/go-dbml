CREATE TABLE "orders" (
  "id" int PRIMARY KEY,
  "total" float,
  "placed" date,
  "status" varchar
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "orders" ("id", "total", "placed", "status")
VALUES
  (1, 10.5, '2024-03-01', 'paid'),
  (2, 20, '2024-03-02', NULL);

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;