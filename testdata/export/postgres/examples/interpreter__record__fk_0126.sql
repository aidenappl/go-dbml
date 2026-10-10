CREATE TABLE "categories" (
  "id" int PRIMARY KEY,
  "name" varchar
);

CREATE TABLE "products" (
  "id" int PRIMARY KEY,
  "category_id" int,
  "name" varchar
);

ALTER TABLE "products" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id") DEFERRABLE INITIALLY IMMEDIATE;

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "categories" ("id", "name")
VALUES
  (1, 'Electronics');
INSERT INTO "products" ("id", "category_id", "name")
VALUES
  (1, 1, 'Laptop'),
  (2, NULL, 'Uncategorized Item');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;