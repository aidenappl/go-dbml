CREATE TABLE "posts" (
  "id" int PRIMARY KEY,
  "created_at" timestamp,
  "title" varchar,
  "content" text,
  "deleted_at" timestamp,
  "is_deleted" boolean,
  "author_id" int
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "posts" ("id", "title", "content", "author_id")
VALUES
  (1, 'First Post', 'Hello World', 100),
  (2, 'Second Post', 'Welcome', 101);

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;