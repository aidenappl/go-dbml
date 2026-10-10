CREATE TYPE "kind" AS ENUM (
  'x',
  'y'
);

CREATE TABLE "users" (
  "id" int PRIMARY KEY,
  "created_by" "s"."s.user_ref",
  "kind" kind,
  "amount" decimal(10,2) CHECK (amount > 0) DEFAULT 1.5,
  CONSTRAINT "users.small" CHECK (amount < 100)
);

CREATE TABLE "posts" (
  "id" int PRIMARY KEY,
  "created_by" "s"."s.user_ref",
  "kind" kind,
  "amount" decimal(10,2) CHECK (amount > 0) DEFAULT 1.5,
  "note2" text,
  CONSTRAINT "posts.small" CHECK (amount < 100)
);

CREATE INDEX "by_creator" ON "users" ("created_by");

CREATE INDEX "by_creator" ON "posts" ("created_by");

COMMENT ON TABLE "users" IS 'audit';

COMMENT ON TABLE "posts" IS 'audit';

COMMENT ON COLUMN "posts"."note2" IS 'n';

ALTER TABLE "users" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "posts" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;
