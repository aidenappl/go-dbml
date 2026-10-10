CREATE TABLE "users" (
  "id" integer PRIMARY KEY,
  "name" varchar(255),
  "email" varchar(255)
);

CREATE TABLE "posts" (
  "id" integer PRIMARY KEY,
  "user_id" integer,
  "title" varchar(255)
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "users" ("id", "name", "email")
VALUES
  (1, 'Alice', 'alice@example.com'),
  (2, 'Bob', 'bob@example.com'),
  (3, 'Charlie', 'charlie@example.com'),
  (4, 'Dave', 'dave@example.com');
INSERT INTO "posts" ("id", "user_id", "title")
VALUES
  (1, 1, 'First Post'),
  (2, 1, 'Second Post'),
  (3, 2, 'Bob Post');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;