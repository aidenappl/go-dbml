CREATE TABLE "T" (
  "s1" int,
  "s2" int,
  "a" int PRIMARY KEY,
  "b" int,
  "m1" int,
  "m2" int,
  "c" int,
  "d" int,
  "e1" int,
  "e2" int
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "T" ("s1", "s2", "a", "b", "m1", "m2", "c", "d", "e1", "e2")
VALUES
  (1, 2, 3, 4, 5, 6, 7, 8, 9, 10);

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;