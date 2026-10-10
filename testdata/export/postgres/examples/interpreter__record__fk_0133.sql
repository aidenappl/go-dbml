CREATE TABLE "users" (
  "id" int PRIMARY KEY,
  "name" varchar
);

CREATE TABLE "posts" (
  "id" int PRIMARY KEY,
  "title" varchar,
  "created_by" int
);

CREATE TABLE "comments" (
  "id" int PRIMARY KEY,
  "content" varchar,
  "created_by" int
);

ALTER TABLE "posts" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "comments" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "users" ("id", "name")
VALUES
  (1, 'Alice'),
  (2, 'Bob');
INSERT INTO "posts" ("id", "title", "created_by")
VALUES
  (1, 'Post 1', 1),
  (2, 'Post 2', 2);
INSERT INTO "comments" ("id", "content", "created_by")
VALUES
  (1, 'Comment 1', 1),
  (2, 'Comment 2', 999);

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;