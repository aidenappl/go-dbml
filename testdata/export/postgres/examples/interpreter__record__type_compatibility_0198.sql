CREATE TYPE "status" AS ENUM (
  'active',
  'inactive'
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
  (1, 'undefined');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;