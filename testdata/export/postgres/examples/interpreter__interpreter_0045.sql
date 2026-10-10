CREATE TABLE "t" (
  "col" varchar
);

CREATE INDEX ON "t" USING HASH ("col");
