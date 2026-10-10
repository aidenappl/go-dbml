CREATE TABLE "users" (
  "id" serial PRIMARY KEY,
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