CREATE TABLE "data" (
  "id" int,
  "active" boolean
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "data" ("id", "active")
VALUES
  (1, TRUE),
  (2, FALSE),
  (3, TRUE),
  (4, FALSE),
  (5, TRUE),
  (6, FALSE),
  (7, TRUE),
  (8, FALSE),
  (9, TRUE),
  (10, FALSE),
  (11, TRUE),
  (12, FALSE),
  (13, TRUE),
  (14, FALSE);

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;