CREATE TABLE "users" (
  "id" int,
  "code" varchar(5),
  "name" char(10)
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "users" ("id", "code", "name")
VALUES
  (1, 'ABCDE', 'Alice'),
  (2, 'XYZ', 'Bob'),
  (3, '', '');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;