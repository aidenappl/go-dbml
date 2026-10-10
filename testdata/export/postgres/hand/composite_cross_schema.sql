CREATE SCHEMA "core";

CREATE SCHEMA "billing";

CREATE TABLE "core"."orders" (
  "merchant_id" int,
  "country_code" char(2)
);

CREATE TABLE "billing"."merchants" (
  "id" int,
  "country_code" char(2),
  PRIMARY KEY ("id", "country_code")
);

ALTER TABLE "core"."orders" ADD CONSTRAINT "fk_comp" FOREIGN KEY ("merchant_id", "country_code") REFERENCES "billing"."merchants" ("id", "country_code") ON DELETE CASCADE ON UPDATE NO ACTION DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "core"."orders" ADD FOREIGN KEY ("merchant_id") REFERENCES "billing"."merchants" ("id") ON DELETE SET NULL DEFERRABLE INITIALLY IMMEDIATE;
