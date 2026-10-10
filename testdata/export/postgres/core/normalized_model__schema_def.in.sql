CREATE SCHEMA "schemaB";

CREATE SCHEMA "ecommerce";

CREATE SCHEMA "schemaA";

CREATE TYPE "job_status" AS ENUM (
  'created2',
  'running2',
  'done2',
  'failure2'
);

CREATE TYPE "gender" AS ENUM (
  'male2',
  'female2'
);

CREATE TYPE "schemaB"."gender" AS ENUM (
  'male',
  'female'
);

CREATE TABLE "users" (
  "id" int PRIMARY KEY,
  "name" varchar,
  "pjs" job_status,
  "pjs2" job_status,
  "pg" "schemaB".gender,
  "pg2" gender
);

CREATE TABLE "products" (
  "id" int PRIMARY KEY,
  "name" varchar
);

CREATE TABLE "b" (
  "id" int,
  "c_id" int
);

CREATE TABLE "c" (
  "id" int,
  "b_id" int
);

CREATE TABLE "unary_expression" (
  "id" int(-1) DEFAULT -2,
  "id2" int(--1) DEFAULT -2,
  "id3" int(+-+---+0.1) DEFAULT 7.2225
);

CREATE TABLE "injected_unary_expression" (
  "id" int,
  "id2" int(--1) DEFAULT -2,
  "id3" int(+-+---+0.1) DEFAULT -7.2225
);

CREATE TABLE "ecommerce"."users" (
  "id" int PRIMARY KEY,
  "name" varchar,
  "ejs" job_status,
  "ejs2" job_status,
  "eg" "schemaB".gender,
  "eg2" gender
);

CREATE TABLE "schemaA"."products" (
  "id" int PRIMARY KEY,
  "name" varchar
);

CREATE TABLE "schemaA"."locations" (
  "id" int PRIMARY KEY,
  "name" varchar
);

ALTER TABLE "schemaA"."products" ADD FOREIGN KEY ("name") REFERENCES "ecommerce"."users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "schemaA"."locations" ADD FOREIGN KEY ("name") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "ecommerce"."users" ADD FOREIGN KEY ("id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "ecommerce"."users" ADD CONSTRAINT "name_optional" FOREIGN KEY ("id") REFERENCES "users" ("name") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "c" ADD CONSTRAINT "short_ref" FOREIGN KEY ("id") REFERENCES "b" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "b" ADD CONSTRAINT "long_ref" FOREIGN KEY ("c_id") REFERENCES "c" ("id") DEFERRABLE INITIALLY IMMEDIATE;
