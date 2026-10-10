CREATE TABLE "data" (
  "id" int,
  "value" decimal
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "data" ("id", "value")
VALUES
  (1, 10000000000),
  (2, 0.0000314),
  (3, 200000000);

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;