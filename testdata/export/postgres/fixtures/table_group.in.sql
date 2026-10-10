CREATE TABLE "users" (
  "id" int,
  "full_name" varchar,
  "created_at" timestamp,
  "country_code" int
);

CREATE TABLE "merchants" (
  "id" int,
  "merchant_name" varchar,
  "country_code" int,
  "created at" varchar,
  "admin_id" int
);

ALTER TABLE "merchants" ADD FOREIGN KEY ("admin_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;
