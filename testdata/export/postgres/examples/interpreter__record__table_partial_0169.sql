CREATE TABLE "users" (
  "extra" varchar,
  "id" int PRIMARY KEY,
  "name" varchar
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "users" ("extra", "id", "name")
VALUES
  ('extra_value', 1, 'Alice'),
  ('extra_value2', 2, 'Bob');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;