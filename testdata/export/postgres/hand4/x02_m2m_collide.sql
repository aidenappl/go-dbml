CREATE TABLE "a" (
  "id" int PRIMARY KEY,
  "b_id" int
);

CREATE TABLE "b" (
  "id" int PRIMARY KEY
);

CREATE TABLE "a_b" (
  "x" int
);

CREATE TABLE "c" (
  "a_id" int,
  "id" int PRIMARY KEY
);

CREATE TABLE "a_b(1)" (
  "a_id" int,
  "b_id" int,
  PRIMARY KEY ("a_id", "b_id")
);

ALTER TABLE "a_b(1)" ADD FOREIGN KEY ("a_id") REFERENCES "a" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "a_b(1)" ADD FOREIGN KEY ("b_id") REFERENCES "b" ("id") DEFERRABLE INITIALLY IMMEDIATE;


CREATE TABLE "a_b(2)" (
  "a_b_id" int,
  "b_id" int,
  PRIMARY KEY ("a_b_id", "b_id")
);

ALTER TABLE "a_b(2)" ADD FOREIGN KEY ("a_b_id") REFERENCES "a" ("b_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "a_b(2)" ADD FOREIGN KEY ("b_id") REFERENCES "b" ("id") DEFERRABLE INITIALLY IMMEDIATE;


CREATE TABLE "a_c" (
  "a_id" int,
  "a_b_id" int,
  "c_a_id" int,
  "c_id" int,
  PRIMARY KEY ("a_id", "a_b_id", "c_a_id", "c_id")
);

ALTER TABLE "a_c" ADD FOREIGN KEY ("a_id", "a_b_id") REFERENCES "a" ("id", "b_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "a_c" ADD FOREIGN KEY ("c_a_id", "c_id") REFERENCES "c" ("a_id", "id") DEFERRABLE INITIALLY IMMEDIATE;

