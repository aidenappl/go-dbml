CREATE TABLE "employees" (
  "emp_id" "NUMBER(10)" NOT NULL,
  "emp_name" "VARCHAR2(100)" NOT NULL,
  "email" "VARCHAR2(100)" NOT NULL,
  "phone" "VARCHAR2(20)",
  "salary" "NUMBER(10,2)",
  "commission_pct" "NUMBER(3,2)",
  "hire_date" DATE,
  "department_id" "NUMBER(10)" NOT NULL,
  "manager_id" "NUMBER(10)" NOT NULL,
  "status" "VARCHAR2(20)" NOT NULL,
  "ssn" "VARCHAR2(11)" NOT NULL,
  "badge_number" "VARCHAR2(10)" NOT NULL,
  "title" "VARCHAR2(100)" NOT NULL,
  "office_location" "VARCHAR2(50)" NOT NULL
);

CREATE TABLE "products" (
  "product_id" "NUMBER(10)" NOT NULL,
  "product_name" "VARCHAR2(200)" NOT NULL,
  "description" "VARCHAR2(4000)",
  "price" "NUMBER(10,2)" NOT NULL,
  "category" "VARCHAR2(50)" NOT NULL,
  "manufacturer" "VARCHAR2(100)",
  "sku" "VARCHAR2(50)" NOT NULL,
  "weight" "NUMBER(8,2)"
);

CREATE TABLE "orders" (
  "order_id" "NUMBER(10)" NOT NULL,
  "customer_id" "NUMBER(10)" NOT NULL,
  "order_date" DATE NOT NULL,
  "ship_date" DATE NOT NULL,
  "total_amount" "NUMBER(12,2)" NOT NULL,
  "tax_amount" "NUMBER(10,2)" NOT NULL,
  "discount_amount" "NUMBER(10,2)" NOT NULL
);

CREATE TABLE "datatype_tests" (
  "col_number" "NUMBER(10,2)",
  "col_varchar2" "VARCHAR2(100)",
  "col_char" "CHAR(10)",
  "col_date" DATE,
  "col_timestamp" TIMESTAMP,
  "col_clob" "CLOB",
  "col_raw" "RAW(100)",
  "col_float" FLOAT,
  "col_integer" INTEGER,
  "col_nvarchar2" "NVARCHAR2(100)",
  "col_nchar" "NCHAR(10)"
);

CREATE TABLE "complex_test" (
  "id" "NUMBER(10)" PRIMARY KEY,
  "code" "VARCHAR2(20)" UNIQUE,
  "name" "VARCHAR2(100)" NOT NULL,
  "parent_id" "NUMBER(10)" NOT NULL,
  "status" "VARCHAR2(20)"
);

CREATE TABLE "combo_test" (
  "col1" "NUMBER(10)",
  "col2" "VARCHAR2(50)"
);

ALTER TABLE "complex_test" ADD CONSTRAINT "fk_complex_parent" FOREIGN KEY ("parent_id") REFERENCES "complex_test" ("id") DEFERRABLE INITIALLY IMMEDIATE;
