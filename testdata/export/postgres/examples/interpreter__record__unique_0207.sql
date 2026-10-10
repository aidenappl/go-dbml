CREATE TABLE "users" (
  "id" int PRIMARY KEY,
  "email" varchar UNIQUE
);

CREATE TABLE "products" (
  "id" int PRIMARY KEY,
  "sku" varchar UNIQUE
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "users" ("id", "email")
VALUES
  (1, 'alice@example.com'),
  (2, 'alice@example.com');
INSERT INTO "products" ("id", "sku")
VALUES
  (1, 'PROD-001'),
  (2, 'PROD-001');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;