CREATE SCHEMA "auth";

CREATE TABLE "auth"."users" (
  "id" int PRIMARY KEY,
  "email" varchar
);

CREATE TABLE "auth"."roles" (
  "id" int PRIMARY KEY,
  "name" varchar
);

CREATE TABLE "auth"."sessions" (
  "id" int PRIMARY KEY
);
