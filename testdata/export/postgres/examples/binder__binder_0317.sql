CREATE SCHEMA "types";

CREATE TYPE "types"."status" AS ENUM (
  'active',
  'inactive'
);

CREATE TABLE "users" (
  "id" int,
  "status" types.status
);
