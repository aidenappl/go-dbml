CREATE TYPE "job_status" AS ENUM (
  'created',
  'running',
  'done',
  'failed',
  'wait for validation'
);

CREATE TYPE "order status" AS ENUM (
  'created',
  'pending',
  'processing',
  'completed'
);

CREATE TABLE "jobs" (
  "id" integer PRIMARY KEY,
  "status" job_status
);

CREATE TABLE "orders" (
  "id" int UNIQUE PRIMARY KEY,
  "status" "order status",
  "created_at" varchar
);

COMMENT ON COLUMN "jobs"."status" IS 'This is a column note';
