CREATE TABLE "employees" (
  "id" int,
  "manager_id" int
);

ALTER TABLE "employees" ADD FOREIGN KEY ("manager_id") REFERENCES "employees" ("id") DEFERRABLE INITIALLY IMMEDIATE;
