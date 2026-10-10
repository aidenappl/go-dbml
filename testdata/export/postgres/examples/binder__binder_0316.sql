CREATE TYPE "status" AS ENUM (
  'active',
  'inactive'
);

CREATE TABLE "users" (
  "id" int,
  "status" status
);
