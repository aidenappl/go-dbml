CREATE TYPE "priority" AS ENUM (
  'low',
  'medium',
  'high'
);

CREATE TABLE "tasks" (
  "id" int PRIMARY KEY,
  "name" varchar,
  "priority" priority
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "tasks" ("id", "name", "priority")
VALUES
  (1, 'Task 1', 'high'),
  (2, 'Task 2', 'invalid_priority');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;