CREATE TABLE "customers" (
  "cust_id" "NUMBER(10)" UNIQUE,
  "cust_name" "VARCHAR2(100)" UNIQUE,
  "email" "VARCHAR2(100)" UNIQUE,
  "phone" "VARCHAR2(20)" UNIQUE
);

CREATE TABLE "inventory" (
  "item_id" "NUMBER(10)" UNIQUE,
  "item_name" "VARCHAR2(200)" UNIQUE,
  "price" "NUMBER(10,2)",
  "sku" "VARCHAR2(50)" UNIQUE
);

CREATE TABLE "shipments" (
  "ship_id" "NUMBER(10)" UNIQUE,
  "customer_id" "NUMBER(10)",
  "ship_date" DATE UNIQUE,
  "tracking_number" "VARCHAR2(50)" UNIQUE
);

CREATE TABLE "departments" (
  "dept_id" "NUMBER(10)" UNIQUE,
  "dept_name" "VARCHAR2(100)" UNIQUE,
  "manager_id" "NUMBER(10)" UNIQUE,
  "budget" "NUMBER(12,2)"
);

CREATE TABLE "suppliers" (
  "supp_id" "NUMBER(10)" UNIQUE,
  "supp_name" "VARCHAR2(100)" UNIQUE,
  "contact" "VARCHAR2(100)" UNIQUE,
  "ssn" "VARCHAR2(11)" UNIQUE
);

CREATE TABLE "test_datatypes" (
  "col_number" "NUMBER(10,2)" UNIQUE,
  "col_varchar2" "VARCHAR2(100)" UNIQUE,
  "col_char" "CHAR(10)" UNIQUE,
  "col_date" DATE UNIQUE,
  "col_timestamp" TIMESTAMP UNIQUE,
  "col_clob" "CLOB" UNIQUE,
  "col_raw" "RAW(100)" UNIQUE,
  "col_float" FLOAT UNIQUE,
  "col_integer" INTEGER UNIQUE,
  "col_nvarchar2" "NVARCHAR2(100)" UNIQUE,
  "col_nchar" "NCHAR(10)" UNIQUE
);

CREATE TABLE "test_basic" (
  "id" "NUMBER(10)" UNIQUE,
  "value" "NUMBER(10)" UNIQUE
);

CREATE TABLE "test_named" (
  "code" "VARCHAR2(20)" UNIQUE,
  "ref_id" "NUMBER(10)" UNIQUE
);

CREATE TABLE "test_status" (
  "status" "VARCHAR2(20)" UNIQUE,
  "status_code" "VARCHAR2(10)" UNIQUE
);

CREATE TABLE "test_dates" (
  "event_date" DATE UNIQUE,
  "end_date" DATE UNIQUE
);

CREATE TABLE "test_enable" (
  "level" "NUMBER(2)" UNIQUE,
  "rating" "NUMBER(2)" UNIQUE
);

CREATE TABLE "test_disable" (
  "priority" "NUMBER(2)" UNIQUE,
  "rank" "NUMBER(2)" UNIQUE
);

CREATE TABLE "test_validate" (
  "amount" "NUMBER(10,2)" UNIQUE,
  "tax" "NUMBER(10,2)" UNIQUE
);

CREATE TABLE "test_novalidate" (
  "discount" "NUMBER(10,2)" UNIQUE,
  "margin" "NUMBER(10,2)" UNIQUE
);

CREATE TABLE "test_disable_novalidate" (
  "weight" "NUMBER(8,2)" UNIQUE,
  "volume" "NUMBER(8,2)" UNIQUE
);

CREATE TABLE "test_rely" (
  "code" "NUMBER(10)" UNIQUE,
  "ref_id" "NUMBER(10)" UNIQUE
);

CREATE TABLE "test_norely" (
  "sequence" "NUMBER(10)" UNIQUE,
  "counter" "NUMBER(10)" UNIQUE
);

CREATE TABLE "test_enable_rely" (
  "batch_id" "NUMBER(10)" UNIQUE,
  "lot_number" "VARCHAR2(20)" UNIQUE
);

CREATE TABLE "test_disable_norely" (
  "zone_id" "NUMBER(5)" UNIQUE,
  "area_code" "VARCHAR2(10)" UNIQUE
);

CREATE TABLE "test_deferrable_deferred" (
  "order_id" "NUMBER(10)" UNIQUE,
  "line_item" "NUMBER(5)" UNIQUE
);

CREATE TABLE "test_deferrable_immediate" (
  "invoice_id" "NUMBER(10)" UNIQUE,
  "payment_status" "VARCHAR2(20)" UNIQUE
);

CREATE TABLE "test_not_deferrable" (
  "ticket_id" "NUMBER(10)" UNIQUE,
  "issue_date" DATE UNIQUE
);

CREATE TABLE "test_datatypes_inline" (
  "col_number" "NUMBER(10,2)" UNIQUE,
  "col_varchar2" "VARCHAR2(100)" UNIQUE,
  "col_char" "CHAR(10)" UNIQUE,
  "col_date" DATE UNIQUE,
  "col_timestamp" TIMESTAMP UNIQUE,
  "col_clob" "CLOB" UNIQUE,
  "col_raw" "RAW(100)" UNIQUE,
  "col_float" FLOAT UNIQUE,
  "col_integer" INTEGER UNIQUE,
  "col_nvarchar2" "NVARCHAR2(100)" UNIQUE,
  "col_nchar" "NCHAR(10)" UNIQUE
);

CREATE TABLE "test_long_name" (
  "identifier" "NUMBER(10)" UNIQUE,
  "record_code" "VARCHAR2(20)" UNIQUE
);

CREATE TABLE "test_complex" (
  "score" "NUMBER(3)" UNIQUE,
  "rating" "VARCHAR2(20)" UNIQUE
);
