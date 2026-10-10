CREATE TABLE "data" (
  "id" int,
  "count" int,
  "active" boolean,
  "created_at" timestamp,
  "amount" decimal,
  "name" varchar,
  "description" text
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "data" ("id", "count", "active", "created_at", "amount", "name", "description")
VALUES
  (1, NULL, NULL, NULL, NULL, NULL, NULL),
  (2, NULL, NULL, NULL, NULL, NULL, NULL),
  (3, 10, TRUE, '2024-01-15T10:30:00+07:00', 99.99, 'test', 'content');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;