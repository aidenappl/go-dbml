CREATE TABLE "users" (
  "id" int PRIMARY KEY,
  "name" varchar
);

CREATE TABLE "countries" (
  "code" varchar(2) PRIMARY KEY,
  "name" varchar
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "users" ("id", "name")
VALUES
  (1, 'Alice'),
  (1, 'Bob');
INSERT INTO "countries" ("code", "name")
VALUES
  ('US', 'United States'),
  ('US', 'USA');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;