CREATE TABLE "users" (
  "id" int PRIMARY KEY NOT NULL DEFAULT 1,
  "name" varchar
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "users" ("id", "name")
VALUES
  (NULL, 'Alice'),
  (NULL, 'Bob');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;