CREATE TABLE "users" (
  "id" int PRIMARY KEY,
  "name" varchar
);

CREATE TABLE "posts" (
  "id" int PRIMARY KEY,
  "title" varchar,
  "user_id" int
);

ALTER TABLE "posts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "users" ("id", "name")
VALUES
  (1, 'Alice'),
  (2, 'Bob');
INSERT INTO "posts" ("id", "title", "user_id")
VALUES
  (1, 'Post 1', 1),
  (2, 'Post 2', 2);

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;