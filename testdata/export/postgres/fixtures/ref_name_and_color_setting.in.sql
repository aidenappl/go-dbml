CREATE TABLE "b" (
  "id" int,
  "c_id" int
);

CREATE TABLE "c" (
  "id" int,
  "b_id" int
);

ALTER TABLE "c" ADD CONSTRAINT "short_ref" FOREIGN KEY ("id") REFERENCES "b" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "b" ADD CONSTRAINT "long_ref" FOREIGN KEY ("c_id") REFERENCES "c" ("id") DEFERRABLE INITIALLY IMMEDIATE;
