CREATE TABLE "data" (
  "id" int,
  "count" integer,
  "small" smallint,
  "big" bigint
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "data" ("id", "count", "small", "big")
VALUES
  (1, 42, -100, 9999999999),
  (0, 0, 0, 0);

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;