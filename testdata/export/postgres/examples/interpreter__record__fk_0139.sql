CREATE TABLE "products" (
  "id" int,
  "region" varchar,
  PRIMARY KEY ("id", "region")
);

CREATE TABLE "orders" (
  "id" int PRIMARY KEY,
  "product_id" int,
  "region" varchar
);

ALTER TABLE "orders" ADD FOREIGN KEY ("product_id", "region") REFERENCES "products" ("id", "region") DEFERRABLE INITIALLY IMMEDIATE;

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "products" ("id", "region")
VALUES
  (1, 'US'),
  (NULL, 'EU');
INSERT INTO "orders" ("id", "product_id", "region")
VALUES
  (1, 1, 'US');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;