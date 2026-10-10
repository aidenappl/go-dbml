CREATE TABLE "events" (
  "id" int,
  "created_at" timestamp,
  "event_date" date
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "events" ("id", "created_at", "event_date")
VALUES
  (1, '2024-01-15T10:30:00+07:00', '2024-01-15');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;