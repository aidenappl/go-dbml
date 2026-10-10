CREATE TYPE "status" AS ENUM (
  'active',
  'inactive'
);

CREATE TABLE "users" (
  "id" int,
  "quantity" int,
  "name" varchar(5),
  "price" decimal(5,2),
  "active" boolean,
  "status" status
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "users" ("id", "quantity", "name", "price", "active", "status")
VALUES
  (invalid, 10, 'VeryLongName', 12345.67, FALSE, 'invalid_status');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;