CREATE TABLE "data" (
  "id" int,
  "name" varchar,
  "title" text,
  "content" text
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "data" ("id", "name", "title", "content")
VALUES
  (1, 'Alice', 'Bob', 'Line 1
Line 2'),
  (2, '', '', 'Tab	Separated'),
  (3, '123', 'true', 'Quote: "test"');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;