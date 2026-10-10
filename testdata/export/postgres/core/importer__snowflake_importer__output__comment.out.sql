CREATE SCHEMA "TEST_SCHEMA";

CREATE TABLE "TEST_SCHEMA"."ONE" (
  "employee_id" INT NOT NULL,
  "first_name" "VARCHAR(50)" NOT NULL,
  "last_name" "VARCHAR(50)" NOT NULL,
  "email" "VARCHAR(100)",
  "hire_date" DATE,
  "department" "VARCHAR(50)",
  "salary" "DECIMAL(10,2)"
);

CREATE TABLE "TEST_SCHEMA"."TWO" (
  "employee_id" INT NOT NULL,
  "first_name" "VARCHAR(50)" NOT NULL,
  "last_name" "VARCHAR(50)" NOT NULL,
  "email" "VARCHAR(100)",
  "hire_date" DATE,
  "department" "VARCHAR(50)",
  "salary" "DECIMAL(10,2)"
);

COMMENT ON COLUMN "TEST_SCHEMA"."ONE"."employee_id" IS 'Unique identifier for each employee';

COMMENT ON COLUMN "TEST_SCHEMA"."ONE"."first_name" IS 'Employee''s first name';

COMMENT ON COLUMN "TEST_SCHEMA"."ONE"."last_name" IS 'Employee''s last name';

COMMENT ON COLUMN "TEST_SCHEMA"."ONE"."email" IS 'Employee''s work email address';

COMMENT ON COLUMN "TEST_SCHEMA"."ONE"."hire_date" IS 'Date when the employee was hired';

COMMENT ON COLUMN "TEST_SCHEMA"."ONE"."department" IS 'Department where the employee works';

COMMENT ON COLUMN "TEST_SCHEMA"."ONE"."salary" IS 'Employee''s annual salary';

COMMENT ON TABLE "TEST_SCHEMA"."TWO" IS 'Another comment for employees table';

COMMENT ON COLUMN "TEST_SCHEMA"."TWO"."employee_id" IS 'Unique identifier for each employee';

COMMENT ON COLUMN "TEST_SCHEMA"."TWO"."first_name" IS 'Employee''s first name';

COMMENT ON COLUMN "TEST_SCHEMA"."TWO"."last_name" IS 'Employee''s last name';

COMMENT ON COLUMN "TEST_SCHEMA"."TWO"."email" IS 'Employee''s work email address';

COMMENT ON COLUMN "TEST_SCHEMA"."TWO"."hire_date" IS 'Date when the employee was hired';

COMMENT ON COLUMN "TEST_SCHEMA"."TWO"."department" IS 'Department where the employee works';

COMMENT ON COLUMN "TEST_SCHEMA"."TWO"."salary" IS 'Employee''s annual salary';
