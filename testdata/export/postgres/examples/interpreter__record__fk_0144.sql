CREATE TABLE "tags" (
  "id" int PRIMARY KEY
);

CREATE TABLE "posts" (
  "id" int PRIMARY KEY
);

CREATE TABLE "tags_posts" (
  "tags_id" int,
  "posts_id" int,
  PRIMARY KEY ("tags_id", "posts_id")
);

ALTER TABLE "tags_posts" ADD FOREIGN KEY ("tags_id") REFERENCES "tags" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "tags_posts" ADD FOREIGN KEY ("posts_id") REFERENCES "posts" ("id") DEFERRABLE INITIALLY IMMEDIATE;


-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "tags" ("id")
VALUES
  (NULL),
  (1);
INSERT INTO "posts" ("id")
VALUES
  (10);

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;