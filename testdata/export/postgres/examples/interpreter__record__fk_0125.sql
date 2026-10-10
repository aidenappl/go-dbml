CREATE TABLE "users" (
  "id" int PRIMARY KEY,
  "name" varchar
);

CREATE TABLE "countries" (
  "code" varchar(2) PRIMARY KEY,
  "name" varchar
);

CREATE TABLE "posts" (
  "id" int PRIMARY KEY,
  "user_id" int,
  "title" varchar
);

CREATE TABLE "cities" (
  "id" int PRIMARY KEY,
  "country_code" varchar(2),
  "name" varchar
);

ALTER TABLE "posts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "cities" ADD FOREIGN KEY ("country_code") REFERENCES "countries" ("code") DEFERRABLE INITIALLY IMMEDIATE;

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "users" ("id", "name")
VALUES
  (1, 'Alice');
INSERT INTO "countries" ("code", "name")
VALUES
  ('US', 'United States');
INSERT INTO "posts" ("id", "user_id", "title")
VALUES
  (1, 1, 'Valid Post'),
  (2, 999, 'Invalid FK');
INSERT INTO "cities" ("id", "country_code", "name")
VALUES
  (1, 'US', 'New York'),
  (2, 'FR', 'Paris');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;