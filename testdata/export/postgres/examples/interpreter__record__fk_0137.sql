CREATE TABLE "users" (
  "id" integer PRIMARY KEY,
  "username" varchar,
  "role" varchar,
  "created_at" timestamp
);

CREATE TABLE "posts" (
  "id" integer PRIMARY KEY,
  "title" varchar,
  "body" text,
  "user_id" integer NOT NULL,
  "status" varchar,
  "created_at" timestamp
);

COMMENT ON COLUMN "posts"."body" IS 'Content of the post';

ALTER TABLE "posts" ADD CONSTRAINT "user_posts" FOREIGN KEY ("user_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "users" ("id", "username", "role")
VALUES
  (0, 'Alice', 'admin'),
  (1, 'Bob', 'moderator'),
  (2, 'Candice', 'moderator'),
  (3, 'David', 'member');
INSERT INTO "posts" ("id", "title", "user_id")
VALUES
  (0, 'Welcome to the forum!', 0),
  (1, 'Guidelines', 1),
  (2, 'Hello all!', 3);

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;