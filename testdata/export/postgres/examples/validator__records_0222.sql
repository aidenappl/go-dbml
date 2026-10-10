CREATE TABLE "data" (
  "id" int,
  "value" int
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "data" ("id", "value")
VALUES
  (1, -100),
  (2, -999);

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;