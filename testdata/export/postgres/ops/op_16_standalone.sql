CREATE TABLE "a" (
  "id" int PRIMARY KEY,
  "b_id" int NOT NULL,
  "n" int
);

CREATE TABLE "b" (
  "id" int PRIMARY KEY,
  "u" int UNIQUE
);

CREATE TABLE "a_b" (
  "a_b_id" int,
  "b_id" int,
  PRIMARY KEY ("a_b_id", "b_id")
);

ALTER TABLE "a_b" ADD FOREIGN KEY ("a_b_id") REFERENCES "a" ("b_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "a_b" ADD FOREIGN KEY ("b_id") REFERENCES "b" ("id") DEFERRABLE INITIALLY IMMEDIATE;


CREATE TABLE "a_b(1)" (
  "a_n" int,
  "b_u" int,
  PRIMARY KEY ("a_n", "b_u")
);

ALTER TABLE "a_b(1)" ADD FOREIGN KEY ("a_n") REFERENCES "a" ("n") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "a_b(1)" ADD FOREIGN KEY ("b_u") REFERENCES "b" ("u") DEFERRABLE INITIALLY IMMEDIATE;

