CREATE TABLE "data" (
  "id" int,
  "count" int,
  "active" boolean,
  "name" varchar
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "data" ("id", "count", "active", "name")
VALUES
  (NULL, NULL, NULL, '');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;