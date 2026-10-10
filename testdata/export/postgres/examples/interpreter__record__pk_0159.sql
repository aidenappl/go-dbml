CREATE TABLE "items" (
  "id" int PRIMARY KEY DEFAULT 0,
  "name" varchar
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "items" ("name")
VALUES
  ('Alice'),
  ('Bob');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;