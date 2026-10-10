CREATE SCHEMA "ecommerce";

CREATE SCHEMA "schemaA";

CREATE SCHEMA "schemaB";

CREATE TYPE "users_pjs_enum" AS ENUM (
  'created2',
  'running2',
  'done2',
  'failure2'
);

CREATE TYPE "users_pjs2_enum" AS ENUM (
  'created2',
  'running2',
  'done2',
  'failure2'
);

CREATE TYPE "users_pg_enum" AS ENUM (
  'male',
  'female'
);

CREATE TYPE "users_pg2_enum" AS ENUM (
  'male2',
  'female2'
);

CREATE TYPE "ecommerce"."users_ejs_enum" AS ENUM (
  'created2',
  'running2',
  'done2',
  'failure2'
);

CREATE TYPE "ecommerce"."users_ejs2_enum" AS ENUM (
  'created2',
  'running2',
  'done2',
  'failure2'
);

CREATE TYPE "ecommerce"."users_eg_enum" AS ENUM (
  'male',
  'female'
);

CREATE TYPE "ecommerce"."users_eg2_enum" AS ENUM (
  'male2',
  'female2'
);

CREATE TABLE "users" (
  "id" int PRIMARY KEY,
  "name" varchar(255),
  "pjs" users_pjs_enum,
  "pjs2" users_pjs2_enum,
  "pg" users_pg_enum,
  "pg2" users_pg2_enum
);

CREATE TABLE "products" (
  "id" int PRIMARY KEY,
  "name" varchar(255)
);

CREATE TABLE "orders" (
  "id1" int,
  "id2" int,
  "pid" int,
  "name" varchar(255),
  "name2" varchar(255),
  PRIMARY KEY ("id1", "id2")
);

CREATE TABLE "orders2" (
  "id1" int,
  "id2" int,
  "pid" int,
  PRIMARY KEY ("id1", "id2")
);

CREATE TABLE "orders4" (
  "id1" int,
  "id2" int,
  PRIMARY KEY ("id1", "id2")
);

CREATE TABLE "ecommerce"."users" (
  "id" int PRIMARY KEY,
  "name" varchar(255),
  "ejs" ecommerce.users_ejs_enum,
  "ejs2" ecommerce.users_ejs2_enum,
  "eg" ecommerce.users_eg_enum,
  "eg2" ecommerce.users_eg2_enum
);

CREATE TABLE "schemaA"."products" (
  "id" int PRIMARY KEY,
  "name" varchar(255),
  "lid" int
);

CREATE TABLE "schemaA"."locations" (
  "id" int PRIMARY KEY,
  "name" varchar(255)
);

CREATE TABLE "schemaB"."orders3" (
  "id1" int PRIMARY KEY,
  "id2" int
);

CREATE UNIQUE INDEX "I_UNIQUE_1" ON "orders" USING BTREE ("name", "name2");

CREATE INDEX "INDEX_2" ON "orders" USING HASH ("id1", "id2");

COMMENT ON COLUMN "products"."name" IS 'Product name';

COMMENT ON TABLE "orders" IS 'this is table orders';

COMMENT ON COLUMN "schemaA"."products"."name" IS 'Sample field comment on multiples schema';

ALTER TABLE "schemaA"."products" ADD CONSTRAINT "FK_1" FOREIGN KEY ("lid") REFERENCES "schemaA"."locations" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "orders" ADD FOREIGN KEY ("pid") REFERENCES "schemaA"."products" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "orders2" ADD FOREIGN KEY ("pid") REFERENCES "schemaA"."products" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "orders2" ADD CONSTRAINT "CFK_1" FOREIGN KEY ("id1", "id2") REFERENCES "orders" ("id1", "id2") ON DELETE CASCADE ON UPDATE SET NULL DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "ecommerce"."users" ADD FOREIGN KEY ("id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "ecommerce"."users" ADD CONSTRAINT "name_optional" FOREIGN KEY ("id") REFERENCES "users" ("name") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "schemaA"."products" ADD FOREIGN KEY ("name") REFERENCES "ecommerce"."users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "schemaA"."locations" ADD FOREIGN KEY ("name") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;
