CREATE TYPE "Status" AS ENUM (
  'active',
  'inactive',
  'pending'
);

CREATE TABLE "users" (
  "id" int,
  "name" varchar,
  "age" int,
  "status" "Status",
  "created_at" timestamp
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "users" ("id", "name", "age", "status", "created_at")
VALUES
  (1, 'Alice', 30, 'undefined', '2024-01-15T10:30:00'),
  (2, 'Bob', NULL, 'inactive', now()),
  (3, 'Charlie', NULL, 'undefined', '2024-01-15');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;