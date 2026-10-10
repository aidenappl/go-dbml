CREATE TABLE "categories" (
  "id" int PRIMARY KEY,
  "parent_id" int,
  "name" varchar
);

ALTER TABLE "categories" ADD FOREIGN KEY ("parent_id") REFERENCES "categories" ("id") DEFERRABLE INITIALLY IMMEDIATE;
