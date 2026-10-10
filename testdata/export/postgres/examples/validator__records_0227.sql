CREATE TABLE "users" (
  "user-id" int PRIMARY KEY,
  "user-name" varchar
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "users" ("user-id", "user-name")
VALUES
  (1, 'Alice');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;