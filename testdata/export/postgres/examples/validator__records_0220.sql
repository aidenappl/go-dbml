CREATE TABLE "users" (
  "id" int PRIMARY KEY,
  "created_at" timestamp
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "users" ("id", "created_at")
VALUES
  (1, now()),
  (2, uuid_generate_v4());

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;