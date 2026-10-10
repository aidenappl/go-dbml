CREATE TABLE "users" (
  "id" int PRIMARY KEY,
  "email" varchar UNIQUE,
  "phone" varchar UNIQUE
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "users" ("id", "email", "phone")
VALUES
  (1, 'alice@example.com', NULL),
  (2, 'bob@example.com', NULL),
  (3, NULL, '123-456'),
  (4, NULL, '789-012');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;