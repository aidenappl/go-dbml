CREATE SCHEMA "v2";

CREATE TYPE "v2"."job_status" AS ENUM (
  'created',
  'running',
  'done',
  'failure'
);

CREATE TABLE "v2"."jobs" (
  "status" v2.job_status DEFAULT 'created'
);
