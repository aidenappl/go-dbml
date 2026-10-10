CREATE TABLE "departments" (
  "id" integer PRIMARY KEY,
  "name" nvarchar2(255),
  "budget" "number(10, 2)"
);

CREATE TABLE "employees" (
  "id" integer PRIMARY KEY,
  "dept_id" integer,
  "name" nvarchar2(255)
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "departments" ("id", "name", "budget")
VALUES
  (1, 'Engineering', 500000),
  (2, 'Marketing', 300000),
  (3, 'Sales', 400000),
  (4, 'HR', 250000);
INSERT INTO "employees" ("id", "dept_id", "name")
VALUES
  (1, 1, 'John Doe'),
  (2, 1, 'Jane Smith'),
  (3, 2, 'Bob Johnson');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;