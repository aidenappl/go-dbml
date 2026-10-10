CREATE TABLE "users" (
  "id" int PRIMARY KEY,
  "name" varchar,
  "email" varchar
);

CREATE TABLE "posts" (
  "id" int PRIMARY KEY,
  "title" varchar,
  "published" boolean
);

CREATE TABLE "comments" (
  "id" int PRIMARY KEY,
  "user_id" int,
  "post_id" int,
  "title" string
);

ALTER TABLE "comments" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "comments" ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("id") DEFERRABLE INITIALLY IMMEDIATE;

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "users" ("id", "name", "email")
VALUES
  (1, 'Alice', 'alice@example.com'),
  (2, 'Bob', 'bob@example.com');
INSERT INTO "posts" ("id", "title", "published")
VALUES
  (1, 'First Post', TRUE),
  (2, 'Second Post', FALSE);
INSERT INTO "comments" ("id", "user_id", "post_id", "title")
VALUES
  (1, 2, 1, 'First comment of first post by the second user');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;