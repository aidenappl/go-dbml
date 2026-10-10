CREATE TABLE "events" (
  "created_at" timestamp,
  "event_date" date,
  "event_time" time
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "events" ("created_at", "event_date", "event_time")
VALUES
  ('2024-01-15T10:30:00', '2024-01-15', '10:30:00'),
  ('2024-12-31T23:59:59', '2024-12-31', '23:59:59');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;