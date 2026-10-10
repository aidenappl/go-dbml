CREATE SCHEMA "s1";

CREATE TABLE "s1"."t1" (
  "a" int,
  "b" int,
  "c" text,
  PRIMARY KEY ("a", "b")
);

CREATE TABLE "s1"."t3" (
  "a" int PRIMARY KEY
);

CREATE TABLE "t2" (
  "a" int
);

CREATE UNIQUE INDEX ON "s1"."t1" USING HASH ("a");

CREATE INDEX "expr_ix" ON "s1"."t1" USING BTREE ((lower(c)), "b");

CREATE INDEX ON "s1"."t1" ((a + b));

CREATE INDEX ON "s1"."t1" ("c");

CREATE INDEX ON "s1"."t3" ("a");

CREATE INDEX ON "t2" ("a");

CREATE UNIQUE INDEX ON "t2" ("a");
