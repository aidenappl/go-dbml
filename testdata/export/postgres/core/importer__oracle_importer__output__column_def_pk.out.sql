CREATE TABLE "customers" (
  "cust_id" "NUMBER(10)" PRIMARY KEY,
  "cust_name" "VARCHAR2(100)",
  "email" "VARCHAR2(100)",
  "phone" "VARCHAR2(20)"
);

CREATE TABLE "inventory" (
  "item_id" "NUMBER(10)" PRIMARY KEY,
  "item_name" "VARCHAR2(200)",
  "price" "NUMBER(10,2)",
  "category" "VARCHAR2(50)"
);

CREATE TABLE "shipments" (
  "ship_id" "NUMBER(10)" PRIMARY KEY,
  "customer_id" "NUMBER(10)",
  "ship_date" DATE,
  "total_amount" "NUMBER(12,2)"
);

CREATE TABLE "departments" (
  "dept_id" "NUMBER(10)" PRIMARY KEY,
  "dept_name" "VARCHAR2(100)",
  "manager_id" "NUMBER(10)",
  "budget" "NUMBER(12,2)"
);

CREATE TABLE "suppliers" (
  "supp_id" "NUMBER(10)" PRIMARY KEY,
  "supp_name" "VARCHAR2(100)",
  "contact" "VARCHAR2(100)",
  "credit_limit" "NUMBER(10,2)"
);

CREATE TABLE "test_datatypes" (
  "col_number" "NUMBER(10,2)" PRIMARY KEY,
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

CREATE TABLE "test_basic" (
  "id" "NUMBER(10)" PRIMARY KEY,
  "value" "NUMBER(10)"
);

CREATE TABLE "test_named" (
  "code" "NUMBER(10)" PRIMARY KEY,
  "description" "VARCHAR2(200)"
);

CREATE TABLE "test_status" (
  "status_code" "VARCHAR2(20)" PRIMARY KEY,
  "status_desc" "VARCHAR2(100)"
);

CREATE TABLE "test_dates" (
  "event_date" DATE PRIMARY KEY,
  "event_desc" "VARCHAR2(200)"
);

CREATE TABLE "test_enable" (
  "level" "NUMBER(2)" PRIMARY KEY,
  "level_desc" "VARCHAR2(50)"
);

CREATE TABLE "test_disable" (
  "priority" "NUMBER(2)" PRIMARY KEY,
  "priority_desc" "VARCHAR2(50)"
);

CREATE TABLE "test_validate" (
  "amount" "NUMBER(10,2)" PRIMARY KEY,
  "amount_desc" "VARCHAR2(100)"
);

CREATE TABLE "test_novalidate" (
  "discount" "NUMBER(10,2)" PRIMARY KEY,
  "discount_desc" "VARCHAR2(100)"
);

CREATE TABLE "test_disable_novalidate" (
  "weight" "NUMBER(8,2)" PRIMARY KEY,
  "weight_desc" "VARCHAR2(100)"
);

CREATE TABLE "test_rely" (
  "code" "NUMBER(10)" PRIMARY KEY,
  "code_desc" "VARCHAR2(50)"
);

CREATE TABLE "test_norely" (
  "sequence" "NUMBER(10)" PRIMARY KEY,
  "sequence_desc" "VARCHAR2(50)"
);

CREATE TABLE "test_enable_rely" (
  "batch_id" "NUMBER(10)" PRIMARY KEY,
  "batch_desc" "VARCHAR2(50)"
);

CREATE TABLE "test_disable_norely" (
  "zone_id" "NUMBER(5)" PRIMARY KEY,
  "zone_desc" "VARCHAR2(50)"
);

CREATE TABLE "test_deferrable_deferred" (
  "order_id" "NUMBER(10)" PRIMARY KEY,
  "order_desc" "VARCHAR2(100)"
);

CREATE TABLE "test_deferrable_immediate" (
  "invoice_id" "NUMBER(10)" PRIMARY KEY,
  "invoice_desc" "VARCHAR2(100)"
);

CREATE TABLE "test_not_deferrable" (
  "ticket_id" "NUMBER(10)" PRIMARY KEY,
  "ticket_desc" "VARCHAR2(100)"
);

CREATE TABLE "test_datatypes_inline" (
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
  "col_nchar" "NCHAR(10)",
  PRIMARY KEY ("col_number", "col_char", "col_date", "col_timestamp", "col_clob", "col_raw", "col_float", "col_integer", "col_nvarchar2", "col_nchar")
);

CREATE TABLE "test_long_name" (
  "identifier" "NUMBER(10)" PRIMARY KEY,
  "record_code" "VARCHAR2(20)"
);

CREATE TABLE "test_complex" (
  "score" "NUMBER(3)",
  "rating" "VARCHAR2(20)",
  PRIMARY KEY ("score", "rating")
);
