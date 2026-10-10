CREATE SCHEMA "types";

CREATE TYPE "types"."status" AS ENUM (
  'active',
  'inactive'
);

CREATE TABLE "users" (
  "status" types.status DEFAULT 'active'
);
