CREATE TABLE "emp" (
  "id" int PRIMARY KEY,
  "manager_id" int
);

ALTER TABLE "emp" ADD FOREIGN KEY ("manager_id") REFERENCES "emp" ("id") DEFERRABLE INITIALLY IMMEDIATE;
