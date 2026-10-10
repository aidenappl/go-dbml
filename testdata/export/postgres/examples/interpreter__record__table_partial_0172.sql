CREATE TABLE "users" (
  "created_at" timestamp,
  "updated_at" timestamp,
  "id" int PRIMARY KEY,
  "name" varchar
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "users" ("created_at", "updated_at", "id", "name")
VALUES
  ('2024-01-01T00:00:00+07:00', '2024-01-01T00:00:00+07:00', 1, 'Alice'),
  ('2024-01-02T00:00:00+07:00', '2024-01-02T00:00:00+07:00', 2, 'Bob');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;