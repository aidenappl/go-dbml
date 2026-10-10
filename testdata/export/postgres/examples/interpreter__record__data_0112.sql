CREATE TABLE "data" (
  "name" varchar(255),
  "description" text,
  "code" char(10)
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "data" ("name", "description", "code")
VALUES
  ('Alice', 'A short description', 'ABC123'),
  ('Bob', 'Another description', '');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;