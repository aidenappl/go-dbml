CREATE TABLE "T" (
  "a" varchar PRIMARY KEY,
  "b" text,
  "c" int
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "T" ("a", "b", "c")
VALUES
  ('value_a', 'value_b', 3);

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;