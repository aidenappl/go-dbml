CREATE TABLE "t" (
  "first_name" varchar,
  "last_name" varchar
);

CREATE INDEX ON "t" ((lower(first_name)), "last_name");
