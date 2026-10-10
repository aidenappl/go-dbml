CREATE SCHEMA "ecommerce";

CREATE SCHEMA "schemaA";

CREATE TABLE "users" (
  "id" int PRIMARY KEY,
  "name" nvarchar(255),
  "pjs" nvarchar(255) NOT NULL CHECK ([pjs] IN ('created2', 'running2', 'done2', 'failure2')),
  "pjs2" nvarchar(255) NOT NULL CHECK ([pjs2] IN ('created2', 'running2', 'done2', 'failure2')),
  "pg" nvarchar(255) NOT NULL CHECK ([pg] IN ('male', 'female')),
  "pg2" nvarchar(255) NOT NULL CHECK ([pg2] IN ('male2', 'female2'))
);

CREATE TABLE "products" (
  "id" int PRIMARY KEY,
  "name" nvarchar(255)
);

CREATE TABLE "ecommerce"."users" (
  "id" int PRIMARY KEY,
  "name" nvarchar(255),
  "ejs" nvarchar(255) NOT NULL CHECK ([ejs] IN ('created2', 'running2', 'done2', 'failure2')),
  "ejs2" nvarchar(255) NOT NULL CHECK ([ejs2] IN ('created2', 'running2', 'done2', 'failure2')),
  "eg" nvarchar(255) NOT NULL CHECK ([eg] IN ('male', 'female')),
  "eg2" nvarchar(255) NOT NULL CHECK ([eg2] IN ('male2', 'female2'))
);

CREATE TABLE "schemaA"."products" (
  "id" int PRIMARY KEY,
  "name" nvarchar(255),
  "created_at" varchar(255) DEFAULT (now()),
  "lid" int,
  "lid2" int
);

CREATE TABLE "schemaA"."locations" (
  "id" int PRIMARY KEY,
  "name" nvarchar(255)
);

CREATE INDEX "idx_1" ON "ecommerce"."users" ("name", "ejs");

CREATE UNIQUE INDEX "unique_lid_lid2" ON "schemaA"."products" ("lid", "lid2");

COMMENT ON COLUMN "products"."name" IS 'Product name';

COMMENT ON TABLE "ecommerce"."users" IS 'Note on table users of schema ecommerce';

COMMENT ON TABLE "schemaA"."locations" IS 'This is a note in table "schemaA"."locations"';

ALTER TABLE "schemaA"."products" ADD FOREIGN KEY ("lid") REFERENCES "schemaA"."locations" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "schemaA"."products" ADD CONSTRAINT "FK_1" FOREIGN KEY ("lid2") REFERENCES "schemaA"."locations" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "ecommerce"."users" ADD FOREIGN KEY ("id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "ecommerce"."users" ADD CONSTRAINT "name_optional" FOREIGN KEY ("id") REFERENCES "users" ("name") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "schemaA"."products" ADD FOREIGN KEY ("name") REFERENCES "ecommerce"."users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "schemaA"."locations" ADD FOREIGN KEY ("name") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;
