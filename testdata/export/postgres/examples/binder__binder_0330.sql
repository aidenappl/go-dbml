CREATE TYPE "true" AS ENUM (
  'value',
  'other'
);

CREATE TABLE "users" (
  "status" true DEFAULT 'value'
);
