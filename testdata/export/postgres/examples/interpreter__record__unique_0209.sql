CREATE TABLE "users" (
  "id" int PRIMARY KEY,
  "email" varchar UNIQUE
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "users" ("id", "email")
VALUES
  (1, 'alice@example.com'),
  (1, 'bob@example.com');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;