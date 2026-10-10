CREATE TYPE "status" AS ENUM (
  'active',
  'inactive',
  'pending'
);

CREATE TABLE "users" (
  "id" int PRIMARY KEY,
  "status" status
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "users" ("id", "status")
VALUES
  (1, 'undefined'),
  (2, 'inactive'),
  (3, 'undefined'),
  (4, NULL);

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;