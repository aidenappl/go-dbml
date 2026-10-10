CREATE TABLE "nodes" (
  "id" int PRIMARY KEY,
  "name" varchar,
  "parent_id" int
);

ALTER TABLE "nodes" ADD FOREIGN KEY ("parent_id") REFERENCES "nodes" ("id") DEFERRABLE INITIALLY IMMEDIATE;

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "nodes" ("id", "name", "parent_id")
VALUES
  (1, 'Root', NULL),
  (2, 'Child 1', 1),
  (3, 'Invalid Child', 999);

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;