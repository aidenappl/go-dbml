CREATE TABLE "departments" (
  "dept_id" "NUMBER(10)" PRIMARY KEY
);

CREATE TABLE "customers" (
  "cust_id" "NUMBER(10)" PRIMARY KEY
);

CREATE TABLE "employees" (
  "emp_id" "NUMBER(10)",
  "emp_name" "VARCHAR2(100)",
  "dept_id" "NUMBER(10)"
);

CREATE TABLE "orders" (
  "order_id" "NUMBER(10)",
  "customer_id" "NUMBER(10)",
  "order_date" DATE
);

CREATE TABLE "shipments" (
  "ship_id" "NUMBER(10)",
  "customer_id" "NUMBER(10)" NOT NULL,
  "ship_date" DATE
);

CREATE TABLE "complex_parent" (
  "id" "NUMBER(10)",
  "code" "VARCHAR2(20)",
  PRIMARY KEY ("id", "code")
);

CREATE TABLE "inventory" (
  "item_id" "NUMBER(10)" PRIMARY KEY
);

CREATE TABLE "complex_child" (
  "child_id" "NUMBER(10)",
  "parent_id" "NUMBER(10)",
  "parent_code" "VARCHAR2(20)"
);

CREATE TABLE "order_details" (
  "detail_id" "NUMBER(10)",
  "order_id" "NUMBER(10)",
  "item_id" "NUMBER(10)"
);

CREATE TABLE "assignments" (
  "assignment_id" "NUMBER(10)",
  "emp_id" "NUMBER(10)",
  "dept_id" "NUMBER(10)"
);

CREATE TABLE "project_members" (
  "member_id" "NUMBER(10)",
  "emp_id" "NUMBER(10)",
  "dept_id" "NUMBER(10)"
);

CREATE TABLE "test_datatypes" (
  "col_number" "NUMBER(10,2)" UNIQUE,
  "col_varchar2" "VARCHAR2(100)" UNIQUE
);

CREATE TABLE "test_datatypes_fk" (
  "id" "NUMBER(10)",
  "col_number" "NUMBER(10,2)",
  "col_varchar2" "VARCHAR2(100)",
  "col_char" "CHAR(10)",
  "col_date" DATE
);

CREATE TABLE "employees_self" (
  "emp_id" "NUMBER(10)" PRIMARY KEY,
  "manager_id" "NUMBER(10)",
  "emp_name" "VARCHAR2(100)"
);

CREATE TABLE "order_history" (
  "history_id" "NUMBER(10)",
  "order_id" "NUMBER(10)",
  "history_date" DATE
);

CREATE TABLE "shipment_details" (
  "detail_id" "NUMBER(10)",
  "ship_id" "NUMBER(10)",
  "item_id" "NUMBER(10)"
);

CREATE TABLE "project_assignments" (
  "assign_id" "NUMBER(10)",
  "emp_id" "NUMBER(10)",
  "project_id" "NUMBER(10)"
);

ALTER TABLE "employees" ADD FOREIGN KEY ("dept_id") REFERENCES "departments" ("dept_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "orders" ADD CONSTRAINT "fk_order_cust" FOREIGN KEY ("customer_id") REFERENCES "customers" ("cust_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "shipments" ADD CONSTRAINT "fk_ship_cust" FOREIGN KEY ("customer_id") REFERENCES "customers" ("cust_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "complex_child" ADD CONSTRAINT "fk_child_complex" FOREIGN KEY ("parent_id", "parent_code") REFERENCES "complex_parent" ("id", "code") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "order_details" ADD CONSTRAINT "fk_detail_item" FOREIGN KEY ("item_id") REFERENCES "inventory" ("item_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "assignments" ADD CONSTRAINT "fk_assign_dept" FOREIGN KEY ("dept_id") REFERENCES "departments" ("dept_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "project_members" ADD CONSTRAINT "fk_member_dept" FOREIGN KEY ("dept_id") REFERENCES "departments" ("dept_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "test_datatypes_fk" ADD CONSTRAINT "fk_dt_number" FOREIGN KEY ("col_number") REFERENCES "test_datatypes" ("col_number") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "test_datatypes_fk" ADD CONSTRAINT "fk_dt_varchar2" FOREIGN KEY ("col_varchar2") REFERENCES "test_datatypes" ("col_varchar2") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "employees_self" ADD CONSTRAINT "fk_self_manager" FOREIGN KEY ("manager_id") REFERENCES "employees_self" ("emp_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "order_history" ADD CONSTRAINT "fk_very_long_constraint_name_for_testing_parser_limits" FOREIGN KEY ("order_id") REFERENCES "orders" ("order_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "shipment_details" ADD CONSTRAINT "fk_ship_detail" FOREIGN KEY ("ship_id") REFERENCES "shipments" ("ship_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "project_assignments" ADD CONSTRAINT "fk_assign_emp" FOREIGN KEY ("emp_id") REFERENCES "employees" ("emp_id") DEFERRABLE INITIALLY IMMEDIATE;
