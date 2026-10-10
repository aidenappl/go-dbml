CREATE TABLE "users" (
  "id" int PRIMARY KEY,
  "name" varchar
);

CREATE TABLE "categories" (
  "id" int PRIMARY KEY,
  "name" varchar
);

CREATE TABLE "posts" (
  "id" int PRIMARY KEY,
  "title" varchar,
  "user_id" int,
  "category_id" int
);

ALTER TABLE "posts" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "posts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "users" ("id", "name")
VALUES
  (1, 'Alice');
INSERT INTO "categories" ("id", "name")
VALUES
  (1, 'Tech');
INSERT INTO "posts" ("id", "title", "user_id", "category_id")
VALUES
  (1, 'Valid Post', 1, 1),
  (2, 'Invalid Category', 1, 999),
  (3, 'Invalid User', 999, 1);

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;