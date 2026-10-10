CREATE TABLE "event_registrations" (
  "event_id" int,
  "attendee_id" int,
  "registration_date" timestamp
);

CREATE UNIQUE INDEX ON "event_registrations" ("event_id", "attendee_id");

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "event_registrations" ("event_id", "attendee_id", "registration_date")
VALUES
  (1, 100, '2024-01-01'),
  (1, 101, '2024-01-02'),
  (2, 100, '2024-01-03');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;