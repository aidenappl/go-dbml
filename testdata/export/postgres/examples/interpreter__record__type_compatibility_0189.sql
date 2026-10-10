CREATE TABLE "users" (
  "id" int PRIMARY KEY,
  "name" varchar,
  "email" varchar
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "users" ("id", "name", "email")
VALUES
  (1, NULL, NULL);

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;