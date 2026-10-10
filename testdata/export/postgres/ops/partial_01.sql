CREATE TABLE "users" (
  "id" int PRIMARY KEY
);

CREATE TABLE "posts" (
  "created_at" timestamp DEFAULT (now()),
  "title" varchar,
  "owner_id" int,
  "id" bigint,
  CONSTRAINT "posts.pos" CHECK (id > 0)
);

CREATE INDEX ON "posts" ("created_at");

COMMENT ON TABLE "posts" IS 'p';

ALTER TABLE "posts" ADD FOREIGN KEY ("owner_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "posts" ("id", "title", "owner_id", "created_at")
VALUES
  (1, 'a', 1, '2024-01-01T00:00:00'),
  (1, 'b', 2, NULL);

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;