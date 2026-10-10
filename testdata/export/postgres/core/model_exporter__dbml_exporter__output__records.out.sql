CREATE TABLE "users" (
  "id" integer PRIMARY KEY,
  "name" varchar,
  "active" boolean
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "users" ("id", "name", "active")
VALUES
  (1, 'Alice', NULL),
  (2, 'Bob', NULL),
  (NULL, 'Charlie', TRUE),
  (3, NULL, FALSE);

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;