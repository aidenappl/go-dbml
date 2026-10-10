CREATE TABLE "users" (
  "id" integer PRIMARY KEY,
  "name" varchar(255),
  "email" varchar(255),
  "active" boolean
);

CREATE TABLE "comments" (
  "id" integer PRIMARY KEY,
  "user_id" integer,
  "content" text
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "users" ("id", "name", "email", "active")
VALUES
  (1, 'Alice', 'alice@example.com', TRUE),
  (2, 'Bob', 'bob@example.com', FALSE),
  (3, 'Charlie', 'charlie@example.com', TRUE),
  (4, 'Dave', 'dave@example.com', FALSE);
INSERT INTO "comments" ("id", "user_id", "content")
VALUES
  (1, 1, 'Great post!'),
  (2, 2, 'Nice article'),
  (3, 1, 'Thanks for sharing');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;