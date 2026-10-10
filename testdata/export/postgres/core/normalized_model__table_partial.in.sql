CREATE SCHEMA "a";

CREATE TABLE "a" (
  "id" int,
  "id2" int
);

CREATE TABLE "users" (
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now()),
  "name" varchar,
  "email" varchar,
  "id" int,
  "delete_status" boolean NOT NULL,
  "deleted_at" timestamp DEFAULT (now()),
  "refColumn1" bit,
  "refColumn2" bool,
  PRIMARY KEY ("id")
);

CREATE TABLE "testColWithTypeHaveArgs" (
  "id" int,
  "callExpColInTest1" int(4),
  "callExpColInTest2" decimal(10,2) DEFAULT 0.1
);

CREATE TABLE "table_no_column_with_partial" (
  "email" varchar UNIQUE
);

CREATE TABLE "a"."a" (
  "id" int,
  "id2" int
);

CREATE UNIQUE INDEX ON "users" ("email");

CREATE UNIQUE INDEX ON "users" ("email", "id");

CREATE UNIQUE INDEX ON "users" ("email");

CREATE UNIQUE INDEX ON "table_no_column_with_partial" ("email");

COMMENT ON COLUMN "testColWithTypeHaveArgs"."callExpColInTest1" IS 'simple column with call expression as type';

COMMENT ON COLUMN "testColWithTypeHaveArgs"."callExpColInTest2" IS 'column with type that have arguments and default value being a float number';

ALTER TABLE "users" ADD FOREIGN KEY ("refColumn1") REFERENCES "a" ("id2") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "users" ADD FOREIGN KEY ("refColumn2") REFERENCES "a"."a" ("id") DEFERRABLE INITIALLY IMMEDIATE;
