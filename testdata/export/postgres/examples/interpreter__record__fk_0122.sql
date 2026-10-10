CREATE TABLE "products" (
  "id" int,
  "region" varchar,
  PRIMARY KEY ("id", "region")
);

CREATE TABLE "categories" (
  "id" int,
  "region" varchar,
  PRIMARY KEY ("id", "region")
);

CREATE TABLE "products_categories" (
  "products_id" int,
  "products_region" varchar,
  "categories_id" int,
  "categories_region" varchar,
  PRIMARY KEY ("products_id", "products_region", "categories_id", "categories_region")
);

ALTER TABLE "products_categories" ADD FOREIGN KEY ("products_id", "products_region") REFERENCES "products" ("id", "region") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "products_categories" ADD FOREIGN KEY ("categories_id", "categories_region") REFERENCES "categories" ("id", "region") DEFERRABLE INITIALLY IMMEDIATE;


-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "products" ("id", "region")
VALUES
  (1, 'US'),
  (2, 'US');
INSERT INTO "categories" ("id", "region")
VALUES
  (1, 'US'),
  (3, 'EU');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;