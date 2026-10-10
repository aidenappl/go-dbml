CREATE TABLE "data" (
  "id" int,
  "active" boolean
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "data" ("id", "active")
VALUES
  (1, FALSE),
  (2, FALSE),
  (3, FALSE),
  (4, FALSE),
  (5, FALSE);

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;