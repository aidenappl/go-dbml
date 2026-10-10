CREATE SCHEMA "demographic";

CREATE TYPE "status" AS ENUM (
  'active',
  'churned',
  'inactive'
);

CREATE TYPE "demographic"."gender" AS ENUM (
  'male',
  'female'
);

CREATE TYPE "demographic"."age segment" AS ENUM (
  'toddler',
  'children',
  'teenager',
  'young_adult',
  'elder'
);

CREATE TABLE "user" (
  "name" text,
  "id" int,
  "status" status DEFAULT 'active',
  "gender" demographic.gender DEFAULT 'male',
  "age_type" demographic."age segment" DEFAULT 'young_adult'
);
