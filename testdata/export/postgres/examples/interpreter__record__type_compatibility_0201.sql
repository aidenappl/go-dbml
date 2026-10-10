CREATE TABLE "users" (
  "id" int,
  "name" varchar(5),
  "quantity" int
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "users" ("id", "name", "quantity")
VALUES
  (1, 'Alice', 10),
  (2, 'Bob', 20),
  (3, 'Christopher', 30),
  (4, 'Dave', 40),
  (5, 'Elizabeth', 50);

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;