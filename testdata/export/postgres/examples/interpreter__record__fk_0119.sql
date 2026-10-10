CREATE TABLE "merchants" (
  "id" int,
  "country_code" varchar,
  PRIMARY KEY ("id", "country_code")
);

CREATE TABLE "orders" (
  "id" int PRIMARY KEY,
  "merchant_id" int,
  "country" varchar,
  "amount" decimal
);

ALTER TABLE "orders" ADD FOREIGN KEY ("merchant_id", "country") REFERENCES "merchants" ("id", "country_code") DEFERRABLE INITIALLY IMMEDIATE;

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "merchants" ("id", "country_code")
VALUES
  (1, 'US'),
  (1, 'UK'),
  (2, 'US');
INSERT INTO "orders" ("id", "merchant_id", "country", "amount")
VALUES
  (1, 1, 'US', 100),
  (2, 1, 'UK', 200.5),
  (3, 2, 'US', 50);

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;