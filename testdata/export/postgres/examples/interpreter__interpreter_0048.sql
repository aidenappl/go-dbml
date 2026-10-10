CREATE TABLE "t" (
  "first_name" varchar,
  "last_name" varchar
);

CREATE INDEX "idx_fullname" ON "t" ("first_name", "last_name");
