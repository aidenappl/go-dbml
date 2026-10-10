CREATE TABLE "T" (
  "a" int,
  "b" int,
  "x" int PRIMARY KEY,
  "shared" varchar,
  "c" int,
  "y" int
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "T" ("a", "b", "x", "shared", "c", "y")
VALUES
  (1, 2, 3, 'shared_value', 4, 5);

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;