CREATE TABLE "a" (
  "created" int,
  "id" int
);

CREATE TABLE "b" (
  "id" int,
  "created" int
);

CREATE INDEX "idx_created" ON "a" ("created");

CREATE INDEX "idx_created" ON "b" ("created");
