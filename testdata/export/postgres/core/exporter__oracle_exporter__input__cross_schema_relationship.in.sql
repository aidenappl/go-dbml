CREATE SCHEMA "C##test01";

CREATE SCHEMA "C##test02";

CREATE TABLE "users" (
  "id" int PRIMARY KEY
);

CREATE TABLE "products" (
  "id" int PRIMARY KEY,
  "user_id" int NOT NULL
);

CREATE TABLE "C##test01"."users" (
  "id" int PRIMARY KEY
);

CREATE TABLE "C##test01"."products" (
  "id" int PRIMARY KEY,
  "user_id" int NOT NULL
);

CREATE TABLE "C##test02"."users" (
  "id" int PRIMARY KEY
);

CREATE TABLE "C##test02"."products" (
  "id" int PRIMARY KEY,
  "user_id" int NOT NULL
);

ALTER TABLE "products" ADD FOREIGN KEY ("user_id") REFERENCES "C##test01"."users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "C##test01"."products" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "C##test01"."products" ADD FOREIGN KEY ("user_id") REFERENCES "C##test01"."users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

CREATE TABLE "C##test02"."users_products" (
  "users_id" int,
  "products_id" int,
  PRIMARY KEY ("users_id", "products_id")
);

ALTER TABLE "C##test02"."users_products" ADD FOREIGN KEY ("users_id") REFERENCES "C##test02"."users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "C##test02"."users_products" ADD FOREIGN KEY ("products_id") REFERENCES "products" ("id") DEFERRABLE INITIALLY IMMEDIATE;


CREATE TABLE "users_products(1)" (
  "users_id" int,
  "products_id" int,
  PRIMARY KEY ("users_id", "products_id")
);

ALTER TABLE "users_products(1)" ADD FOREIGN KEY ("users_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "users_products(1)" ADD FOREIGN KEY ("products_id") REFERENCES "C##test02"."products" ("id") DEFERRABLE INITIALLY IMMEDIATE;


CREATE TABLE "C##test01"."users_products(2)" (
  "users_id" int,
  "products_id" int,
  PRIMARY KEY ("users_id", "products_id")
);

ALTER TABLE "C##test01"."users_products(2)" ADD FOREIGN KEY ("users_id") REFERENCES "C##test01"."users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "C##test01"."users_products(2)" ADD FOREIGN KEY ("products_id") REFERENCES "C##test02"."products" ("id") DEFERRABLE INITIALLY IMMEDIATE;

