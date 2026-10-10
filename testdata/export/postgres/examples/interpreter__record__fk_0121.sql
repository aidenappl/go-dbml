CREATE TABLE "merchants" (
  "id" int,
  "country_code" varchar,
  PRIMARY KEY ("id", "country_code")
);

CREATE TABLE "orders" (
  "id" int PRIMARY KEY,
  "merchant_id" int,
  "country" varchar,
  "status" varchar
);

ALTER TABLE "orders" ADD FOREIGN KEY ("merchant_id", "country") REFERENCES "merchants" ("id", "country_code") DEFERRABLE INITIALLY IMMEDIATE;

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "merchants" ("id", "country_code")
VALUES
  (1, 'US');
INSERT INTO "orders" ("id", "merchant_id", "country", "status")
VALUES
  (1, 1, 'US', 'confirmed'),
  (2, NULL, 'UK', 'pending'),
  (3, 1, NULL, 'processing');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;