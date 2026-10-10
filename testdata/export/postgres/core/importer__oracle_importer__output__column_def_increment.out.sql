CREATE TABLE "customers" (
  "cust_id" SERIAL,
  "cust_name" "VARCHAR2(100)" DEFAULT 'Unknown',
  "email" "VARCHAR2(100)" DEFAULT 'noemail@example.com',
  "phone" "VARCHAR2(20)" DEFAULT 'N/A'
);

CREATE TABLE "inventory" (
  "item_id" SERIAL,
  "item_name" "VARCHAR2(200)" DEFAULT 'Default Item',
  "price" "NUMBER(10,2)" DEFAULT 0,
  "category" "VARCHAR2(50)" DEFAULT 'Misc'
);

CREATE TABLE "shipments" (
  "ship_id" SERIAL,
  "customer_id" "NUMBER(10)" DEFAULT 1,
  "ship_date" DATE DEFAULT (SYSDATE),
  "total_amount" "NUMBER(12,2)" DEFAULT 0
);

CREATE TABLE "departments" (
  "dept_id" SERIAL,
  "dept_name" "VARCHAR2(100)" DEFAULT 'General',
  "manager_id" "NUMBER(10)" DEFAULT 0,
  "budget" "NUMBER(12,2)" DEFAULT 10000
);

CREATE TABLE "suppliers" (
  "supp_id" SERIAL,
  "supp_name" "VARCHAR2(100)" DEFAULT 'Anonymous',
  "contact" "VARCHAR2(100)" DEFAULT 'Contact',
  "credit_limit" "NUMBER(10,2)" DEFAULT 5000
);

CREATE TABLE "test_datatypes" (
  "col_number" SERIAL,
  "col_varchar2" "VARCHAR2(100)" DEFAULT 'Default Text',
  "col_char" "CHAR(10)" DEFAULT 'DEFAULT   ',
  "col_date" DATE DEFAULT (SYSDATE)
);

CREATE TABLE "test_basic" (
  "id" SERIAL,
  "value" "NUMBER(10)" DEFAULT 100
);

CREATE TABLE "test_status" (
  "status_id" SERIAL,
  "status" "VARCHAR2(20)" DEFAULT 'OPEN'
);

CREATE TABLE "test_not_null" (
  "record_id" SERIAL NOT NULL,
  "description" "VARCHAR2(200)" NOT NULL DEFAULT 'No description'
);

CREATE TABLE "test_start_with" (
  "order_id" SERIAL,
  "order_code" "VARCHAR2(20)" DEFAULT 'ORD000'
);

CREATE TABLE "test_increment" (
  "ticket_id" SERIAL,
  "ticket_type" "VARCHAR2(20)" DEFAULT 'STANDARD'
);

CREATE TABLE "test_maxvalue" (
  "ref_id" SERIAL,
  "ref_name" "VARCHAR2(100)" DEFAULT 'Reference'
);

CREATE TABLE "test_minvalue" (
  "batch_id" SERIAL,
  "batch_name" "VARCHAR2(50)" DEFAULT 'Batch'
);
