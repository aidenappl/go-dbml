CREATE SCHEMA "C##test";

CREATE TABLE "users" (
  "id" int PRIMARY KEY
);

CREATE TABLE "products" (
  "id" int PRIMARY KEY
);

CREATE TABLE "C##test"."users" (
  "id" int PRIMARY KEY
);

CREATE TABLE "C##test"."products" (
  "id" int PRIMARY KEY
);

CREATE TABLE "users_products" (
  "users_id" int,
  "products_id" int,
  PRIMARY KEY ("users_id", "products_id")
);

ALTER TABLE "users_products" ADD FOREIGN KEY ("users_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "users_products" ADD FOREIGN KEY ("products_id") REFERENCES "products" ("id") DEFERRABLE INITIALLY IMMEDIATE;


CREATE TABLE "C##test"."users_products(1)" (
  "users_id" int,
  "products_id" int,
  PRIMARY KEY ("users_id", "products_id")
);

ALTER TABLE "C##test"."users_products(1)" ADD FOREIGN KEY ("users_id") REFERENCES "C##test"."users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "C##test"."users_products(1)" ADD FOREIGN KEY ("products_id") REFERENCES "C##test"."products" ("id") DEFERRABLE INITIALLY IMMEDIATE;

