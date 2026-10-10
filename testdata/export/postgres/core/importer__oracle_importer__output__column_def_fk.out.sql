CREATE TABLE "departments" (
  "dept_id" "NUMBER(10)",
  "dept_name" "VARCHAR2(100)",
  PRIMARY KEY ("dept_id")
);

CREATE TABLE "suppliers" (
  "supp_id" "NUMBER(10)",
  "supp_name" "VARCHAR2(100)",
  PRIMARY KEY ("supp_id")
);

CREATE TABLE "customers" (
  "cust_id" "NUMBER(10)",
  "cust_name" "VARCHAR2(100)",
  PRIMARY KEY ("cust_id")
);

CREATE TABLE "inventory" (
  "item_id" "NUMBER(10)",
  "item_name" "VARCHAR2(200)",
  PRIMARY KEY ("item_id")
);

CREATE TABLE "complex_test" (
  "id" "NUMBER(10)",
  "code" "VARCHAR2(20)",
  PRIMARY KEY ("id")
);

CREATE TABLE "test_datatypes" (
  "col_number" "NUMBER(10,2)",
  "col_varchar2" "VARCHAR2(100)",
  PRIMARY KEY ("col_number")
);

CREATE TABLE "employees" (
  "emp_id" "NUMBER(10)",
  "emp_name" "VARCHAR2(100)",
  "dept_id" "NUMBER(10)",
  "manager_id" "NUMBER(10)"
);

CREATE TABLE "orders" (
  "order_id" "NUMBER(10)",
  "customer_id" "NUMBER(10)",
  "tracking_number" "VARCHAR2(50)"
);

CREATE TABLE "shipments" (
  "ship_id" "NUMBER(10)",
  "customer_id" "NUMBER(10)" NOT NULL,
  "tracking_number" "VARCHAR2(50)" PRIMARY KEY
);

CREATE TABLE "complex_child" (
  "child_id" "NUMBER(10)",
  "parent_id" "NUMBER(10)",
  "parent_code" "VARCHAR2(20)"
);

CREATE TABLE "test_enable" (
  "id" "NUMBER(10)",
  "dept_id" "NUMBER(10)",
  "supp_id" "NUMBER(10)"
);

CREATE TABLE "test_disable" (
  "id" "NUMBER(10)",
  "dept_id" "NUMBER(10)",
  "cust_id" "NUMBER(10)"
);

CREATE TABLE "test_validate" (
  "id" "NUMBER(10)",
  "item_id" "NUMBER(10)",
  "supp_id" "NUMBER(10)"
);

CREATE TABLE "test_novalidate" (
  "id" "NUMBER(10)",
  "dept_id" "NUMBER(10)",
  "cust_id" "NUMBER(10)"
);

CREATE TABLE "test_disable_novalidate" (
  "id" "NUMBER(10)",
  "item_id" "NUMBER(10)",
  "supp_id" "NUMBER(10)"
);

CREATE TABLE "test_rely" (
  "id" "NUMBER(10)",
  "dept_id" "NUMBER(10)",
  "cust_id" "NUMBER(10)"
);

CREATE TABLE "test_norely" (
  "id" "NUMBER(10)",
  "item_id" "NUMBER(10)",
  "supp_id" "NUMBER(10)"
);

CREATE TABLE "test_enable_rely" (
  "id" "NUMBER(10)",
  "dept_id" "NUMBER(10)",
  "cust_id" "NUMBER(10)"
);

CREATE TABLE "test_disable_norely" (
  "id" "NUMBER(10)",
  "item_id" "NUMBER(10)",
  "supp_id" "NUMBER(10)"
);

CREATE TABLE "test_deferrable_deferred" (
  "id" "NUMBER(10)",
  "dept_id" "NUMBER(10)",
  "cust_id" "NUMBER(10)"
);

CREATE TABLE "test_deferrable_immediate" (
  "id" "NUMBER(10)",
  "item_id" "NUMBER(10)",
  "supp_id" "NUMBER(10)"
);

CREATE TABLE "test_not_deferrable" (
  "id" "NUMBER(10)",
  "dept_id" "NUMBER(10)",
  "cust_id" "NUMBER(10)"
);

CREATE TABLE "test_datatypes_fk" (
  "id" "NUMBER(10)",
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

CREATE TABLE "employees_self" (
  "emp_id" "NUMBER(10)",
  "manager_id" "NUMBER(10)",
  "dept_id" "NUMBER(10)"
);

CREATE TABLE "order_details" (
  "order_id" "NUMBER(10)",
  "item_id" "NUMBER(10)",
  "cust_id" "NUMBER(10)"
);

CREATE TABLE "test_cascade" (
  "id" "NUMBER(10)",
  "dept_id" "NUMBER(10)",
  "supp_id" "NUMBER(10)"
);

CREATE TABLE "test_set_null" (
  "id" "NUMBER(10)",
  "dept_id" "NUMBER(10)",
  "cust_id" "NUMBER(10)"
);

ALTER TABLE "employees" ADD CONSTRAINT "fk_emp_dept" FOREIGN KEY ("dept_id") REFERENCES "departments" ("dept_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "employees" ADD CONSTRAINT "fk_emp_manager" FOREIGN KEY ("manager_id") REFERENCES "employees" ("emp_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "orders" ADD CONSTRAINT "fk_order_cust" FOREIGN KEY ("customer_id") REFERENCES "customers" ("cust_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "orders" ADD CONSTRAINT "fk_order_track" FOREIGN KEY ("tracking_number") REFERENCES "shipments" ("tracking_number") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "shipments" ADD CONSTRAINT "fk_ship_cust" FOREIGN KEY ("customer_id") REFERENCES "customers" ("cust_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "complex_child" ADD CONSTRAINT "fk_child_complex" FOREIGN KEY ("parent_id") REFERENCES "complex_test" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "complex_child" ADD CONSTRAINT "fk_child_code" FOREIGN KEY ("parent_code") REFERENCES "complex_test" ("code") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "test_enable" ADD CONSTRAINT "fk_enable_dept" FOREIGN KEY ("dept_id") REFERENCES "departments" ("dept_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "test_enable" ADD CONSTRAINT "fk_enable_supp" FOREIGN KEY ("supp_id") REFERENCES "suppliers" ("supp_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "test_disable" ADD CONSTRAINT "fk_disable_dept" FOREIGN KEY ("dept_id") REFERENCES "departments" ("dept_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "test_disable" ADD CONSTRAINT "fk_disable_cust" FOREIGN KEY ("cust_id") REFERENCES "customers" ("cust_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "test_validate" ADD CONSTRAINT "fk_validate_item" FOREIGN KEY ("item_id") REFERENCES "inventory" ("item_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "test_validate" ADD CONSTRAINT "fk_validate_supp" FOREIGN KEY ("supp_id") REFERENCES "suppliers" ("supp_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "test_novalidate" ADD CONSTRAINT "fk_noval_dept" FOREIGN KEY ("dept_id") REFERENCES "departments" ("dept_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "test_novalidate" ADD CONSTRAINT "fk_noval_cust" FOREIGN KEY ("cust_id") REFERENCES "customers" ("cust_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "test_disable_novalidate" ADD CONSTRAINT "fk_dis_noval_item" FOREIGN KEY ("item_id") REFERENCES "inventory" ("item_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "test_disable_novalidate" ADD CONSTRAINT "fk_dis_noval_supp" FOREIGN KEY ("supp_id") REFERENCES "suppliers" ("supp_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "test_rely" ADD CONSTRAINT "fk_rely_dept" FOREIGN KEY ("dept_id") REFERENCES "departments" ("dept_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "test_rely" ADD CONSTRAINT "fk_rely_cust" FOREIGN KEY ("cust_id") REFERENCES "customers" ("cust_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "test_norely" ADD CONSTRAINT "fk_norely_item" FOREIGN KEY ("item_id") REFERENCES "inventory" ("item_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "test_norely" ADD CONSTRAINT "fk_norely_supp" FOREIGN KEY ("supp_id") REFERENCES "suppliers" ("supp_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "test_enable_rely" ADD CONSTRAINT "fk_en_rely_dept" FOREIGN KEY ("dept_id") REFERENCES "departments" ("dept_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "test_enable_rely" ADD CONSTRAINT "fk_en_rely_cust" FOREIGN KEY ("cust_id") REFERENCES "customers" ("cust_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "test_disable_norely" ADD CONSTRAINT "fk_dis_norely_item" FOREIGN KEY ("item_id") REFERENCES "inventory" ("item_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "test_disable_norely" ADD CONSTRAINT "fk_dis_norely_supp" FOREIGN KEY ("supp_id") REFERENCES "suppliers" ("supp_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "test_deferrable_deferred" ADD CONSTRAINT "fk_def_deferred_dept" FOREIGN KEY ("dept_id") REFERENCES "departments" ("dept_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "test_deferrable_deferred" ADD CONSTRAINT "fk_def_deferred_cust" FOREIGN KEY ("cust_id") REFERENCES "customers" ("cust_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "test_deferrable_immediate" ADD CONSTRAINT "fk_def_immed_item" FOREIGN KEY ("item_id") REFERENCES "inventory" ("item_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "test_deferrable_immediate" ADD CONSTRAINT "fk_def_immed_supp" FOREIGN KEY ("supp_id") REFERENCES "suppliers" ("supp_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "test_not_deferrable" ADD CONSTRAINT "fk_not_def_dept" FOREIGN KEY ("dept_id") REFERENCES "departments" ("dept_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "test_not_deferrable" ADD CONSTRAINT "fk_not_def_cust" FOREIGN KEY ("cust_id") REFERENCES "customers" ("cust_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "test_datatypes_fk" ADD CONSTRAINT "fk_dt_number" FOREIGN KEY ("col_number") REFERENCES "test_datatypes" ("col_number") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "test_datatypes_fk" ADD CONSTRAINT "fk_dt_varchar2" FOREIGN KEY ("col_varchar2") REFERENCES "test_datatypes" ("col_varchar2") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "employees_self" ADD CONSTRAINT "fk_self_emp" FOREIGN KEY ("manager_id") REFERENCES "employees_self" ("emp_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "employees_self" ADD CONSTRAINT "fk_self_dept" FOREIGN KEY ("dept_id") REFERENCES "departments" ("dept_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "order_details" ADD CONSTRAINT "fk_detail_item" FOREIGN KEY ("item_id") REFERENCES "inventory" ("item_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "order_details" ADD CONSTRAINT "fk_detail_cust" FOREIGN KEY ("cust_id") REFERENCES "customers" ("cust_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "test_cascade" ADD CONSTRAINT "fk_cascade_dept" FOREIGN KEY ("dept_id") REFERENCES "departments" ("dept_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "test_cascade" ADD CONSTRAINT "fk_cascade_supp" FOREIGN KEY ("supp_id") REFERENCES "suppliers" ("supp_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "test_set_null" ADD CONSTRAINT "fk_set_null_dept" FOREIGN KEY ("dept_id") REFERENCES "departments" ("dept_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "test_set_null" ADD CONSTRAINT "fk_set_null_cust" FOREIGN KEY ("cust_id") REFERENCES "customers" ("cust_id") DEFERRABLE INITIALLY IMMEDIATE;
