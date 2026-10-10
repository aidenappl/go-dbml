CREATE TABLE "users" (
  "id" int PRIMARY KEY,
  "email" varchar
);

CREATE UNIQUE INDEX ON "users" ("email");
