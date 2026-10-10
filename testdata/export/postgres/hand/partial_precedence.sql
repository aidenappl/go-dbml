CREATE TABLE "users" (
  "name" varchar,
  "created_at" datetime NOT NULL,
  "updated_at" datetime,
  "id" bigint,
  CONSTRAINT "users.chk_upd" CHECK (updated_at >= created_at)
);

CREATE INDEX ON "users" ("created_at");

COMMENT ON TABLE "users" IS 'base note';

COMMENT ON COLUMN "users"."id" IS 'local';
