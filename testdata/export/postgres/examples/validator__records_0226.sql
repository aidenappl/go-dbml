CREATE SCHEMA "auth";

CREATE TABLE "auth"."users" (
  "id" int PRIMARY KEY,
  "email" varchar
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "auth"."users" ("id", "email")
VALUES
  (1, 'alice@example.com');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;