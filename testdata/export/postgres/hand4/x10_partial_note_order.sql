CREATE TABLE "t" (
  "x" int,
  "b" int PRIMARY KEY,
  "a" text
);

CREATE TABLE "t2" (
  "b" int PRIMARY KEY,
  "a" int
);

CREATE UNIQUE INDEX ON "t" ("b");

CREATE INDEX ON "t" ("a");

CREATE INDEX ON "t" ("x");

CREATE INDEX ON "t2" ("a");

CREATE UNIQUE INDEX ON "t2" ("b");

COMMENT ON TABLE "t" IS 'p2 note';

COMMENT ON TABLE "t2" IS 'own';
