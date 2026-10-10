CREATE TABLE "users" (
  "first_name" varchar,
  "last_name" varchar,
  "email" varchar
);

CREATE INDEX ON "users" ("first_name", "last_name");

CREATE UNIQUE INDEX ON "users" ("email");
