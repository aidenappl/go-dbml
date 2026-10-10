CREATE TABLE "users" (
  "id" int PRIMARY KEY,
  "name" varchar
);

CREATE TABLE "user_profiles" (
  "id" int PRIMARY KEY,
  "user_id" int UNIQUE,
  "bio" text
);

CREATE TABLE "departments" (
  "id" int PRIMARY KEY,
  "name" varchar
);

CREATE TABLE "employees" (
  "id" int PRIMARY KEY,
  "dept_id" int,
  "name" varchar
);

ALTER TABLE "users" ADD FOREIGN KEY ("id") REFERENCES "user_profiles" ("user_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "employees" ADD FOREIGN KEY ("dept_id") REFERENCES "departments" ("id") DEFERRABLE INITIALLY IMMEDIATE;

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "users" ("id", "name")
VALUES
  (1, 'Alice'),
  (2, 'Bob');
INSERT INTO "user_profiles" ("id", "user_id", "bio")
VALUES
  (1, 1, 'Alice''s bio'),
  (2, 3, 'Invalid user');
INSERT INTO "departments" ("id", "name")
VALUES
  (1, 'Engineering');
INSERT INTO "employees" ("id", "dept_id", "name")
VALUES
  (1, 1, 'Alice'),
  (2, 999, 'Bob with invalid dept');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;