CREATE TABLE "users" (
  "id" int PRIMARY KEY,
  "name" varchar
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "users" ("id", "name")
VALUES
  (1, 'Alice');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;