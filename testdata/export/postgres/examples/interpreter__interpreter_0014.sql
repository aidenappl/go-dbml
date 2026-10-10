CREATE TABLE "users" (
  "first_name" varchar,
  "last_name" varchar
);

CREATE INDEX ON "users" ("first_name", "last_name");
