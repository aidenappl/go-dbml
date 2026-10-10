CREATE TABLE "articles" (
  "id" int,
  "content" text,
  "description" varchar
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "articles" ("id", "content", "description")
VALUES
  (1, 'This is a very long text content that can be arbitrarily long without any length restrictions because text type does not have a length parameter', 'This is a very long description that can be arbitrarily long');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;