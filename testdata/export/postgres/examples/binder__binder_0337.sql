CREATE TABLE "merchants" (
  "id" int PRIMARY KEY,
  "country_code" varchar
);

CREATE TABLE "orders" (
  "merchant_id" int,
  "country" varchar
);

ALTER TABLE "orders" ADD FOREIGN KEY ("merchant_id", "country") REFERENCES "merchants" ("id", "country_code") DEFERRABLE INITIALLY IMMEDIATE;
