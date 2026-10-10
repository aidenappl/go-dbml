CREATE TABLE "a" (
  "id" int
);

CREATE TABLE "b" (
  "a_id" int
);

CREATE TABLE "a_b" (
  "a_id" int,
  "b_a_id" int,
  PRIMARY KEY ("a_id", "b_a_id")
);

ALTER TABLE "a_b" ADD FOREIGN KEY ("a_id") REFERENCES "a" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "a_b" ADD FOREIGN KEY ("b_a_id") REFERENCES "b" ("a_id") DEFERRABLE INITIALLY IMMEDIATE;

