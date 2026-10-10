CREATE TABLE "customers" (
  "cust_id" "NUMBER(10)" DEFAULT 0,
  "cust_name" "VARCHAR2(100)" DEFAULT 'Unknown',
  "email" "VARCHAR2(100)" DEFAULT 'noemail@example.com',
  "phone" "VARCHAR2(20)" DEFAULT 'N/A',
  "address" "VARCHAR2(200)" DEFAULT 'Unknown Address',
  "status" "VARCHAR2(20)" DEFAULT 'ACTIVE',
  "reg_date" DATE DEFAULT (SYSDATE),
  "last_login" TIMESTAMP DEFAULT (SYSTIMESTAMP)
);

CREATE TABLE "inventory" (
  "item_id" "NUMBER(10)" DEFAULT 1,
  "item_name" "VARCHAR2(200)" DEFAULT 'Item',
  "price" "NUMBER(10,2)" DEFAULT 0,
  "category" "VARCHAR2(50)" DEFAULT 'Misc',
  "sku" "VARCHAR2(50)" DEFAULT 'SKU000',
  "weight" "NUMBER(8,2)" DEFAULT 0,
  "stock_level" "NUMBER(10)" DEFAULT 0,
  "reorder_level" "NUMBER(10)" DEFAULT 10
);

CREATE TABLE "shipments" (
  "ship_id" "NUMBER(10)" DEFAULT (test_seq.NEXTVAL),
  "customer_id" "NUMBER(10)" DEFAULT 1,
  "ship_date" DATE DEFAULT (SYSDATE),
  "total_amount" "NUMBER(12,2)" DEFAULT 0,
  "tracking_number" "VARCHAR2(50)" DEFAULT (SYS_GUID()),
  "carrier" "VARCHAR2(100)" DEFAULT 'UPS',
  "priority" "NUMBER(2)" DEFAULT 1,
  "notes" "VARCHAR2(4000)" DEFAULT 'No notes'
);

CREATE TABLE "departments" (
  "dept_id" "NUMBER(10)" DEFAULT 100,
  "dept_name" "VARCHAR2(100)" DEFAULT 'General',
  "manager_id" "NUMBER(10)" DEFAULT 0,
  "budget" "NUMBER(12,2)" DEFAULT 10000,
  "location" "VARCHAR2(50)" DEFAULT 'HQ',
  "established_date" DATE DEFAULT (SYSDATE)
);

CREATE TABLE "suppliers" (
  "supp_id" "NUMBER(10)" DEFAULT 500,
  "supp_name" "VARCHAR2(100)" DEFAULT 'Supplier',
  "contact" "VARCHAR2(100)" DEFAULT 'Contact',
  "ssn" "VARCHAR2(11)" DEFAULT '000-00-0000',
  "badge_number" "VARCHAR2(10)" DEFAULT 'BADGE000',
  "credit_limit" "NUMBER(10,2)" DEFAULT 5000
);

CREATE TABLE "test_datatypes" (
  "col_number" "NUMBER(10,2)" DEFAULT 42,
  "col_varchar2" "VARCHAR2(100)" DEFAULT 'Default Text',
  "col_char" "CHAR(10)" DEFAULT 'CHARDEF   ',
  "col_date" DATE DEFAULT (SYSDATE),
  "col_timestamp" TIMESTAMP DEFAULT (SYSTIMESTAMP),
  "col_clob" "CLOB" DEFAULT (EMPTY_CLOB()),
  "col_raw" "RAW(100)" DEFAULT (HEXTORAW('00')),
  "col_float" FLOAT DEFAULT 3.14,
  "col_integer" INTEGER DEFAULT 100,
  "col_nvarchar2" "NVARCHAR2(100)" DEFAULT (N'Default National'),
  "col_nchar" "NCHAR(10)" DEFAULT (N'NCHARDEF  ')
);

CREATE TABLE "test_basic_num" (
  "id" "NUMBER(10)" DEFAULT 0,
  "value" "NUMBER(10)" DEFAULT 100
);

CREATE TABLE "test_basic_varchar" (
  "name" "VARCHAR2(100)" DEFAULT 'Unknown',
  "code" "VARCHAR2(10)" DEFAULT 'CODE'
);

CREATE TABLE "test_basic_date" (
  "event_date" DATE DEFAULT (SYSDATE),
  "end_date" DATE DEFAULT (SYSDATE + 1)
);

CREATE TABLE "test_expr_sysdate" (
  "created_at" DATE DEFAULT (SYSDATE),
  "updated_at" DATE DEFAULT (SYSDATE)
);

CREATE TABLE "test_expr_user" (
  "created_by" "VARCHAR2(30)" DEFAULT (USER),
  "modified_by" "VARCHAR2(30)" DEFAULT (USER)
);

CREATE TABLE "test_expr_guid" (
  "unique_id" "VARCHAR2(32)" DEFAULT (SYS_GUID()),
  "session_id" "VARCHAR2(32)" DEFAULT (SYS_GUID())
);

CREATE TABLE "test_expr_seq" (
  "seq_id" "NUMBER(10)" DEFAULT (test_seq.NEXTVAL),
  "auto_id" "NUMBER(10)" DEFAULT (test_seq.NEXTVAL)
);

CREATE TABLE "test_expr_arith" (
  "base_value" "NUMBER(10)" DEFAULT 1000,
  "calc_value" "NUMBER(10)" DEFAULT (1000 + 500)
);

CREATE TABLE "test_on_null_basic" (
  "email" "VARCHAR2(100)" DEFAULT 'noemail@example.com',
  "phone" "VARCHAR2(20)" DEFAULT 'N/A'
);

CREATE TABLE "test_on_null_num" (
  "quantity" "NUMBER(10)" DEFAULT 0,
  "discount" "NUMBER(5,2)" DEFAULT 0
);

CREATE TABLE "test_on_null_date" (
  "start_date" DATE DEFAULT (SYSDATE),
  "end_date" DATE DEFAULT (SYSDATE + 7)
);

CREATE TABLE "test_default_not_null" (
  "status" "VARCHAR2(20)" NOT NULL DEFAULT 'ACTIVE',
  "priority" "NUMBER(2)" NOT NULL DEFAULT 1
);

CREATE TABLE "test_on_null_not_null" (
  "name" "VARCHAR2(100)" NOT NULL DEFAULT 'Unknown',
  "code" "VARCHAR2(10)" NOT NULL DEFAULT 'CODE'
);

CREATE TABLE "test_nn_enable" (
  "level" "NUMBER(2)" NOT NULL DEFAULT 1,
  "rating" "NUMBER(2)" NOT NULL DEFAULT 5
);

CREATE TABLE "test_nn_disable" (
  "rank" "NUMBER(2)" NOT NULL DEFAULT 10,
  "score" "NUMBER(3)" NOT NULL DEFAULT 0
);

CREATE TABLE "test_nn_validate" (
  "amount" "NUMBER(10,2)" NOT NULL DEFAULT 0,
  "tax" "NUMBER(10,2)" NOT NULL DEFAULT 0
);

CREATE TABLE "test_nn_novalidate" (
  "discount" "NUMBER(5,2)" NOT NULL DEFAULT 0,
  "margin" "NUMBER(5,2)" NOT NULL DEFAULT 0.1
);

CREATE TABLE "test_nn_disable_noval" (
  "weight" "NUMBER(8,2)" NOT NULL DEFAULT 0,
  "volume" "NUMBER(8,2)" NOT NULL DEFAULT 0
);

CREATE TABLE "test_nn_rely" (
  "code" "NUMBER(10)" NOT NULL DEFAULT 100,
  "ref_id" "NUMBER(10)" NOT NULL DEFAULT 200
);

CREATE TABLE "test_nn_norely" (
  "sequence" "NUMBER(10)" NOT NULL DEFAULT 1,
  "counter" "NUMBER(10)" NOT NULL DEFAULT 0
);

CREATE TABLE "test_nn_enable_rely" (
  "batch_id" "NUMBER(10)" NOT NULL DEFAULT 500,
  "lot_number" "VARCHAR2(20)" NOT NULL DEFAULT 'LOT'
);

CREATE TABLE "test_nn_disable_norely" (
  "zone_id" "NUMBER(5)" NOT NULL DEFAULT 1,
  "area_code" "VARCHAR2(10)" NOT NULL DEFAULT 'AREA'
);

CREATE TABLE "test_nn_defer_deferred" (
  "order_id" "NUMBER(10)" NOT NULL DEFAULT 1000,
  "line_item" "NUMBER(5)" NOT NULL DEFAULT 1
);

CREATE TABLE "test_nn_defer_immediate" (
  "invoice_id" "NUMBER(10)" NOT NULL DEFAULT 2000,
  "payment_status" "VARCHAR2(20)" NOT NULL DEFAULT 'PENDING'
);

CREATE TABLE "test_nn_not_defer" (
  "ticket_id" "NUMBER(10)" NOT NULL DEFAULT 3000,
  "issue_date" DATE NOT NULL DEFAULT (SYSDATE)
);

CREATE TABLE "test_datatypes_inline" (
  "col_number" "NUMBER(10,2)" DEFAULT 0,
  "col_varchar2" "VARCHAR2(100)" DEFAULT 'Text',
  "col_char" "CHAR(10)" DEFAULT 'CHAR      ',
  "col_date" DATE DEFAULT (SYSDATE),
  "col_timestamp" TIMESTAMP DEFAULT (SYSTIMESTAMP),
  "col_clob" "CLOB" DEFAULT 'CLOB Default',
  "col_raw" "RAW(100)" DEFAULT (HEXTORAW('0000')),
  "col_float" FLOAT DEFAULT 1.23,
  "col_integer" INTEGER DEFAULT 42,
  "col_nvarchar2" "NVARCHAR2(100)" DEFAULT (N'National Text'),
  "col_nchar" "NCHAR(10)" DEFAULT (N'NCHAR     ')
);

CREATE TABLE "test_complex_expr" (
  "full_name" "VARCHAR2(200)" DEFAULT (USER || '_USER'),
  "calc_field" "NUMBER(10)" DEFAULT (100 * 2)
);

CREATE TABLE "test_long_default" (
  "description" "VARCHAR2(4000)" DEFAULT 'This is a very long default value for testing purposes, exceeding typical lengths to check parser limits and storage implications in Oracle SQL',
  "comments" "VARCHAR2(4000)" DEFAULT 'Another long default string to test maximum lengths and handling in table creation statements'
);
