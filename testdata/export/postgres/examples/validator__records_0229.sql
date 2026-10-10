CREATE TABLE "ids" (
  "id" int PRIMARY KEY
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "ids" ("id")
VALUES
  (1),
  (2),
  (3);

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;