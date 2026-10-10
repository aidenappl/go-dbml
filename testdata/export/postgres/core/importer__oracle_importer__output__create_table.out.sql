CREATE SCHEMA "hr";

CREATE SCHEMA "sales";

CREATE SCHEMA "finance";

CREATE SCHEMA "global";

CREATE SCHEMA "hr data";

CREATE SCHEMA "グローバル";

CREATE TABLE "customers" (
  "cust_id" "NUMBER(10)",
  "cust_name" "VARCHAR2(100)"
);

CREATE TABLE "INVENTORY" (
  "item_id" "NUMBER(10)",
  "item_name" "VARCHAR2(200)"
);

CREATE TABLE "shipments" (
  "ship_id" "NUMBER(10)",
  "customer_id" "NUMBER(10)"
);

CREATE TABLE "Customer Data" (
  "data_id" "NUMBER(10)",
  "data_value" "VARCHAR2(100)"
);

CREATE TABLE "order-details$test" (
  "order_id" "NUMBER(10)",
  "detail_desc" "VARCHAR2(200)"
);

CREATE TABLE "sales#report@2023" (
  "report_id" "NUMBER(10)",
  "report_date" DATE
);

CREATE TABLE "客户信息" (
  "client_id" "NUMBER(10)",
  "client_name" "NVARCHAR2(100)"
);

CREATE TABLE "very_long_table_name_for_testing_parser_limits_1234567890" (
  "id" "NUMBER(10)",
  "description" "VARCHAR2(4000)"
);

CREATE TABLE "Sales Report Data 2023" (
  "sales_id" "NUMBER(10)",
  "sales_amount" "NUMBER(12,2)"
);

CREATE TABLE "test@#$%^&*_table!" (
  "test_id" "NUMBER(10)",
  "test_value" "VARCHAR2(100)"
);

CREATE TABLE "project_tasks#1" (
  "task_id" "NUMBER(10)" PRIMARY KEY,
  "task_name" "VARCHAR2(100)"
);

CREATE TABLE "order details" (
  "order_id" "NUMBER(10)",
  "cust_id" "NUMBER(10)",
  "order_date" DATE
);

CREATE TABLE "供应商数据" (
  "supp_id" "NUMBER(10)",
  "supp_name" "VARCHAR2(100)" DEFAULT 'Anonymous'
);

CREATE TABLE "Sales Data Archive" (
  "archive_id" "NUMBER(10)",
  "archive_date" DATE
);

CREATE TABLE "orders-2023-range" (
  "order_id" "NUMBER(10)",
  "order_date" DATE,
  "total_amount" "NUMBER(12,2)"
);

CREATE TABLE "inventory_hash@test" (
  "item_id" "NUMBER(10)",
  "item_name" "VARCHAR2(200)"
);

CREATE TABLE "Temp Customer Data" (
  "temp_id" "NUMBER(10)",
  "temp_name" "VARCHAR2(100)"
);

CREATE TABLE "external_logs#2023" (
  "log_id" "NUMBER(10)",
  "log_message" "VARCHAR2(4000)"
);

CREATE TABLE "Project Tasks#Validated" (
  "task_id" "NUMBER(10)" PRIMARY KEY,
  "task_name" "VARCHAR2(100)" NOT NULL
);

CREATE TABLE "Tasks#NonValidated" (
  "task_id" "NUMBER(10)" PRIMARY KEY,
  "priority" "NUMBER(2)" NOT NULL
);

CREATE TABLE "Suppliers_Norely@Data" (
  "supp_id" "NUMBER(10)" PRIMARY KEY,
  "supp_contact" "VARCHAR2(100)" NOT NULL
);

CREATE TABLE "Contracts@Defer" (
  "contract_id" "NUMBER(10)" PRIMARY KEY,
  "contract_name" "VARCHAR2(100)" NOT NULL
);

CREATE TABLE "Tickets#NonDefer" (
  "ticket_id" "NUMBER(10)" PRIMARY KEY,
  "ticket_type" "VARCHAR2(20)" NOT NULL
);

CREATE TABLE "Test Data Types#2023" (
  "col_number" "NUMBER(10,2)" NOT NULL,
  "col_varchar2" "VARCHAR2(100)" UNIQUE,
  "col_char" "CHAR(10)" DEFAULT 'DEFAULT   ',
  "col_date" DATE PRIMARY KEY,
  "col_timestamp" TIMESTAMP,
  "col_clob" "CLOB",
  "col_raw" "RAW(100)",
  "col_float" FLOAT,
  "col_integer" INTEGER,
  "col_nvarchar2" "NVARCHAR2(100)",
  "col_nchar" "NCHAR(10)"
);

CREATE TABLE "very_long_table_name_1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678" (
  "id" "NUMBER(10)",
  "value" "VARCHAR2(100)"
);

CREATE TABLE "TABLE" (
  "table_id" "NUMBER(10)",
  "table_name" "VARCHAR2(100)"
);

CREATE TABLE "Sales Report#Data 2023@Prod" (
  "sales_id" "NUMBER(10)",
  "sales_amount" "NUMBER(12,2)"
);

CREATE TABLE "Temp#Customer_Records" (
  "temp_id" "NUMBER(10)" PRIMARY KEY,
  "temp_name" "VARCHAR2(100)" NOT NULL
);

CREATE TABLE "Order Details#2023" (
  "order_id" "NUMBER(10)",
  "item_id" "NUMBER(10)",
  "quantity" "NUMBER(10)",
  PRIMARY KEY ("order_id", "item_id")
);

CREATE TABLE "hr"."employees" (
  "emp_id" "NUMBER(10)",
  "emp_name" "VARCHAR2(100)"
);

CREATE TABLE "hr"."Employee Records" (
  "record_id" "NUMBER(10)",
  "record_status" "VARCHAR2(20)"
);

CREATE TABLE "hr"."Employee_Contracts_2023" (
  "contract_id" "NUMBER(10)" NOT NULL,
  "emp_name" "VARCHAR2(100)" NOT NULL
);

CREATE TABLE "hr"."Employee_IOT@Data" (
  "emp_id" "NUMBER(10)" PRIMARY KEY,
  "emp_name" "VARCHAR2(100)"
);

CREATE TABLE "sales"."departments" (
  "dept_id" "NUMBER(10)",
  "dept_name" "VARCHAR2(100)"
);

CREATE TABLE "sales"."Customer_Regions#List" (
  "cust_id" "NUMBER(10)",
  "region" "VARCHAR2(50)"
);

CREATE TABLE "sales"."Orders_Rely#2023" (
  "order_id" "NUMBER(10)" PRIMARY KEY,
  "order_amount" "NUMBER(12,2)" NOT NULL
);

CREATE TABLE "finance"."budget-report#2023" (
  "budget_id" "NUMBER(10)",
  "amount" "NUMBER(12,2)"
);

CREATE TABLE "global"."製品データ" (
  "product_id" "NUMBER(10)",
  "product_name" "NVARCHAR2(200)"
);

CREATE TABLE "hr data"."Employee Data#Report$2023" (
  "report_id" "NUMBER(10)",
  "report_value" "VARCHAR2(200)"
);

CREATE TABLE "グローバル"."製品情報@データ" (
  "prod_id" "NUMBER(10)",
  "prod_name" "NVARCHAR2(200)"
);

ALTER TABLE "order details" ADD CONSTRAINT "fk_order_cust" FOREIGN KEY ("cust_id") REFERENCES "customers" ("cust_id") DEFERRABLE INITIALLY IMMEDIATE;
