CREATE TABLE "users" (
  "id" int PRIMARY KEY,
  "name" varchar
);

CREATE TABLE "countries" (
  "code" varchar(2) PRIMARY KEY,
  "name" varchar
);

CREATE TABLE "accounts" (
  "id" int PRIMARY KEY,
  "name" varchar
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "users" ("id", "name")
VALUES
  (1, 'Alice'),
  (2, 'Bob');
INSERT INTO "countries" ("code", "name")
VALUES
  ('US', 'United States'),
  ('UK', 'United Kingdom');
INSERT INTO "accounts" ("id", "name")
VALUES
  (0, 'Default Account'),
  (1, 'User Account'),
  (2, 'Business Account');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;