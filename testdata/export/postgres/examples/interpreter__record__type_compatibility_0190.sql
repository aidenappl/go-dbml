CREATE TABLE "users" (
  "id" int PRIMARY KEY,
  "name" varchar NOT NULL,
  "status" varchar NOT NULL
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "users" ("id", "name", "status")
VALUES
  (1, NULL, NULL);

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;