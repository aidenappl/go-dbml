CREATE TABLE "data" (
  "int_col" int,
  "float_col" decimal(10,2),
  "bool_col" boolean,
  "str_col" varchar
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "data" ("int_col", "float_col", "bool_col", "str_col")
VALUES
  (1, 3.14, TRUE, 'hello'),
  (2, -2.5, FALSE, 'world');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;