CREATE TABLE "data" (
  "id" int,
  "value" decimal,
  "count" int,
  "temperature" float
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "data" ("id", "value", "count", "temperature")
VALUES
  (1, 10000000000, 100, 314),
  (2, 0.0000314, 2000, -0.0025),
  (3, 200000000, NULL, NULL);

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;