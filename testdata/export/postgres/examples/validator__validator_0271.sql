CREATE TABLE "users" (
  "email" varchar
);

CREATE INDEX ON "users" USING HASH ("email");
