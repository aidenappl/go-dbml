CREATE TABLE "posts" (
  "id" int PRIMARY KEY,
  "title" text
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "posts" ("id", "title")
VALUES
  (1, 'Hello'),
  (2, 'World');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;