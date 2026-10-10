CREATE TABLE "users" (
  "id" int PRIMARY KEY,
  "email" varchar UNIQUE,
  "username" varchar UNIQUE,
  "phone" varchar UNIQUE
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "users" ("id", "email", "username", "phone")
VALUES
  (1, 'alice@example.com', 'alice', '111-111'),
  (2, 'bob@example.com', 'bob', '222-222'),
  (3, 'charlie@example.com', 'alice', '333-333'),
  (4, 'dave@example.com', 'dave', '111-111');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;