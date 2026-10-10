CREATE TABLE "users" (
  "id" int PRIMARY KEY,
  "name" varchar
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "users" ("id", "name")
VALUES
  (1, NULL),
  (2, '');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;