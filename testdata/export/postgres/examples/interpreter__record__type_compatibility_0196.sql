CREATE SCHEMA "auth";

CREATE SCHEMA "app";

CREATE TYPE "auth"."role" AS ENUM (
  'admin',
  'user'
);

CREATE TYPE "app"."status" AS ENUM (
  'active',
  'inactive'
);

CREATE TABLE "auth"."users" (
  "id" int PRIMARY KEY,
  "role" auth.role,
  "status" app.status
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "auth"."users" ("id", "role", "status")
VALUES
  (1, , CAST('active' AS status)),
  (2, , );

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;