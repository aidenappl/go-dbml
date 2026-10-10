CREATE TABLE "events" (
  "id" int,
  "created_at" timestamp
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "events" ("id", "created_at")
VALUES
  (1, '12345'),
  (2, 'true');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;