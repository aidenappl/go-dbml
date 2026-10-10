CREATE TABLE "customers" (
  "cust_id" "NUMBER(10)" DEFAULT 1,
  "cust_name" "VARCHAR2(100)",
  "email" "VARCHAR2(100)" DEFAULT 'email@example.com',
  "phone" "VARCHAR2(20)",
  "address" "VARCHAR2(200)",
  "status" "VARCHAR2(20)" DEFAULT 'churned',
  "reg_date" DATE DEFAULT (SYSDATE),
  "last_login" TIMESTAMP
);

CREATE TABLE "inventory" (
  "item_id" "NUMBER(10)",
  "item_name" "VARCHAR2(200)" DEFAULT 'Unknown',
  "price" "NUMBER(10,2)" DEFAULT 0.99,
  "category" "VARCHAR2(50)" DEFAULT 'Misc',
  "sku" "VARCHAR2(50)",
  "weight" "NUMBER(8,2)",
  "stock_level" "NUMBER(10)",
  "reorder_level" "NUMBER(10)"
);

CREATE TABLE "shipments" (
  "ship_id" "NUMBER(10)",
  "customer_id" "NUMBER(10)" DEFAULT 1,
  "ship_date" DATE DEFAULT (SYSDATE),
  "total_amount" "NUMBER(12,2)",
  "tracking_number" "VARCHAR2(50)" DEFAULT (SYS_GUID()),
  "carrier" "VARCHAR2(100)",
  "priority" "NUMBER(2)",
  "notes" "VARCHAR2(4000)" DEFAULT 'This is a very long default value for testing purposes, exceeding typical lengths to check parser limits and storage implications in Oracle SQL'
);

CREATE TABLE "departments" (
  "dept_id" "NUMBER(10)" DEFAULT (test_seq.NEXTVAL),
  "dept_name" "VARCHAR2(100)",
  "manager_id" "NUMBER(10)",
  "budget" "NUMBER(12,2)" DEFAULT 10000,
  "location" "VARCHAR2(50)" DEFAULT (UPPER(USER) || '_DEPT'),
  "established_date" DATE DEFAULT (SYSDATE)
);

CREATE TABLE "suppliers" (
  "supp_id" "NUMBER(10)",
  "supp_name" "VARCHAR2(100)" DEFAULT 'Anonymous',
  "contact" "VARCHAR2(100)" DEFAULT 'N/A',
  "ssn" "VARCHAR2(11)",
  "badge_number" "VARCHAR2(10)",
  "credit_limit" "NUMBER(10,2)" DEFAULT (5000 + 1000)
);

CREATE TABLE "test_defaults" (
  "col1" "NUMBER",
  "col2" "VARCHAR2(50)",
  "col3" DATE,
  "col4" TIMESTAMP,
  "col5" "CLOB",
  "col6" "RAW(100)",
  "col7" FLOAT,
  "col8" INTEGER,
  "col9" "NVARCHAR2(100)"
);
