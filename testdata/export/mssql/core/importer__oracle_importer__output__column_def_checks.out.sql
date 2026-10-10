CREATE TABLE [customers] (
  [cust_id] NUMBER(10) CHECK (cust_id > 0),
  [cust_name] VARCHAR2(100) CHECK (LENGTH(cust_name) <= 100),
  [email] VARCHAR2(100) CHECK (email LIKE '%@%.%'),
  [phone] VARCHAR2(20) CHECK (LENGTH(phone) >= 10),
  [address] VARCHAR2(200) CHECK (address IS NOT NULL),
  [status] VARCHAR2(20) CHECK (status IN ('ACTIVE', 'INACTIVE')),
  [reg_date] DATE CHECK (reg_date >= TO_DATE('2000-01-01', 'YYYY-MM-DD')),
  [last_login] TIMESTAMP CHECK ((last_login <= SYSTIMESTAMP) AND (last_login >= SYSTIMESTAMP)),
  [last_login_2] TIMESTAMP CHECK ((last_login_2 <= SYSTIMESTAMP) AND (last_login_2 >= SYSTIMESTAMP)),
  [last_login_3] TIMESTAMP CHECK ((last_login_3 <= SYSTIMESTAMP) AND (last_login_3 >= SYSTIMESTAMP))
)
GO

CREATE TABLE [inventory] (
  [item_id] NUMBER(10) CHECK (item_id > 0),
  [item_name] VARCHAR2(200) CHECK (item_name IS NOT NULL),
  [price] NUMBER(10,2) CHECK (price >= 0),
  [category] VARCHAR2(50) CHECK (category IN ('ELECTRONICS', 'CLOTHING', 'BOOKS')),
  [sku] VARCHAR2(50) CHECK (LENGTH(sku) <= 50),
  [weight] NUMBER(8,2) CHECK (weight > 0),
  [stock_level] NUMBER(10) CHECK (stock_level >= 0),
  [reorder_level] NUMBER(10) CHECK (reorder_level >= 0)
)
GO

CREATE TABLE [shipments] (
  [ship_id] NUMBER(10) CHECK (ship_id > 0),
  [customer_id] NUMBER(10) CHECK (customer_id > 0),
  [ship_date] DATE CHECK (ship_date >= TRUNC(SYSDATE)),
  [total_amount] NUMBER(12,2) CHECK (total_amount >= 0),
  [tracking_number] VARCHAR2(50) CHECK (LENGTH(tracking_number) <= 50),
  [carrier] VARCHAR2(100) CHECK (carrier IS NOT NULL),
  [priority] NUMBER(2) CHECK (priority BETWEEN 1 AND 5),
  [notes] VARCHAR2(4000) CHECK (LENGTH(notes) <= 4000)
)
GO

CREATE TABLE [departments] (
  [dept_id] NUMBER(10) CHECK (dept_id > 0),
  [dept_name] VARCHAR2(100) CHECK (dept_name IS NOT NULL),
  [manager_id] NUMBER(10) CHECK (manager_id > 0),
  [budget] NUMBER(12,2) CHECK (budget >= 0),
  [location] VARCHAR2(50) CHECK (location IS NOT NULL),
  [established_date] DATE CHECK (established_date <= SYSDATE)
)
GO

CREATE TABLE [suppliers] (
  [supp_id] NUMBER(10) CHECK (supp_id > 0),
  [supp_name] VARCHAR2(100) CHECK (supp_name IS NOT NULL),
  [contact] VARCHAR2(100) CHECK (LENGTH(contact) <= 100),
  [ssn] VARCHAR2(11) CHECK (ssn LIKE '___-__-____'),
  [badge_number] VARCHAR2(10) CHECK (LENGTH(badge_number) <= 10),
  [credit_limit] NUMBER(10,2) CHECK (credit_limit >= 0)
)
GO

CREATE TABLE [test_datatypes] (
  [col_number] NUMBER(10,2) CHECK (col_number >= 0),
  [col_varchar2] VARCHAR2(100) CHECK (LENGTH(col_varchar2) <= 100),
  [col_char] CHAR(10) CHECK (LENGTH(col_char) = 10),
  [col_date] DATE CHECK (col_date >= TO_DATE('2000-01-01', 'YYYY-MM-DD')),
  [col_timestamp] TIMESTAMP CHECK (col_timestamp <= SYSTIMESTAMP),
  [col_clob] CLOB CHECK (col_clob IS NOT NULL),
  [col_raw] RAW(100) CHECK (LENGTH(col_raw) <= 100),
  [col_float] FLOAT CHECK (col_float >= 0),
  [col_integer] INTEGER CHECK (col_integer >= 0),
  [col_nvarchar2] NVARCHAR2(100) CHECK (LENGTH(col_nvarchar2) <= 100),
  [col_nchar] NCHAR(10) CHECK (LENGTH(col_nchar) = 10)
)
GO

CREATE TABLE [test_basic] (
  [id] NUMBER(10) CHECK (id > 0),
  [value] NUMBER(10) CHECK (value BETWEEN 1 AND 1000)
)
GO

CREATE TABLE [test_status] (
  [status] VARCHAR2(20) CHECK (status IN ('OPEN', 'CLOSED', 'PENDING')),
  [code] VARCHAR2(10) CHECK (LENGTH(code) <= 10)
)
GO

CREATE TABLE [test_not_null] (
  [name] VARCHAR2(100) CHECK (name IS NOT NULL),
  [description] VARCHAR2(200) CHECK (LENGTH(description) <= 200)
)
GO

CREATE TABLE [test_dates] (
  [event_date] DATE CHECK (event_date >= SYSDATE),
  [end_date] DATE CHECK (end_date > SYSDATE)
)
GO

CREATE TABLE [test_compound] (
  [quantity] NUMBER(10) CHECK (quantity >= 0),
  [price] NUMBER(10,2) CHECK (quantity * price >= 0)
)
GO

CREATE TABLE [test_pattern] (
  [email] VARCHAR2(100) CHECK (email LIKE '%@%.%'),
  [phone] VARCHAR2(20) CHECK (phone LIKE '[0-9]%')
)
GO

CREATE TABLE [test_multi_condition] (
  [score] NUMBER(3) CHECK (score BETWEEN 0 AND 100),
  [grade] CHAR(1) CHECK (grade IN ('A', 'B', 'C', 'D', 'F') AND score > 0)
)
GO

CREATE TABLE [test_enable] (
  [level] NUMBER(2) CHECK (level BETWEEN 1 AND 10),
  [rating] NUMBER(2) CHECK (rating >= 0)
)
GO

CREATE TABLE [test_disable] (
  [priority] NUMBER(2) CHECK (priority BETWEEN 1 AND 5),
  [rank] NUMBER(2) CHECK (rank >= 0)
)
GO

CREATE TABLE [test_validate] (
  [amount] NUMBER(10,2) CHECK (amount >= 0),
  [tax] NUMBER(10,2) CHECK (tax >= 0)
)
GO

CREATE TABLE [test_novalidate] (
  [discount] NUMBER(10,2) CHECK (discount >= 0),
  [margin] NUMBER(10,2) CHECK (margin >= 0)
)
GO

CREATE TABLE [test_disable_novalidate] (
  [weight] NUMBER(8,2) CHECK (weight > 0),
  [volume] NUMBER(8,2) CHECK (volume > 0)
)
GO

CREATE TABLE [test_rely] (
  [code] NUMBER(10) CHECK (code > 0),
  [ref_id] NUMBER(10) CHECK (ref_id > 0)
)
GO

CREATE TABLE [test_norely] (
  [sequence] NUMBER(10) CHECK (sequence > 0),
  [counter] NUMBER(10) CHECK (counter >= 0)
)
GO

CREATE TABLE [test_enable_rely] (
  [batch_id] NUMBER(10) CHECK (batch_id > 0),
  [lot_number] VARCHAR2(20) CHECK (LENGTH(lot_number) <= 20)
)
GO

CREATE TABLE [test_disable_norely] (
  [zone_id] NUMBER(5) CHECK (zone_id > 0),
  [area_code] VARCHAR2(10) CHECK (LENGTH(area_code) <= 10)
)
GO

CREATE TABLE [test_deferrable_deferred] (
  [order_id] NUMBER(10) CHECK (order_id > 0),
  [line_item] NUMBER(5) CHECK (line_item > 0)
)
GO

CREATE TABLE [test_deferrable_immediate] (
  [invoice_id] NUMBER(10) CHECK (invoice_id > 0),
  [payment_status] VARCHAR2(20) CHECK (payment_status IN ('PAID', 'PENDING'))
)
GO

CREATE TABLE [test_not_deferrable] (
  [ticket_id] NUMBER(10) CHECK (ticket_id > 0),
  [issue_date] DATE CHECK (issue_date <= SYSDATE)
)
GO

CREATE TABLE [test_datatypes_inline] (
  [col_number] NUMBER(10,2) CHECK (col_number >= 0),
  [col_varchar2] VARCHAR2(100) CHECK (LENGTH(col_varchar2) <= 100),
  [col_char] CHAR(10) CHECK (LENGTH(col_char) = 10),
  [col_date] DATE CHECK (col_date >= TO_DATE('2000-01-01', 'YYYY-MM-DD')),
  [col_timestamp] TIMESTAMP CHECK (col_timestamp <= SYSTIMESTAMP),
  [col_clob] CLOB CHECK (col_clob IS NOT NULL),
  [col_raw] RAW(100) CHECK (LENGTH(col_raw) <= 100),
  [col_float] FLOAT CHECK (col_float >= 0),
  [col_integer] INTEGER CHECK (col_integer >= 0),
  [col_nvarchar2] NVARCHAR2(100) CHECK (LENGTH(col_nvarchar2) <= 100),
  [col_nchar] NCHAR(10) CHECK (LENGTH(col_nchar) = 10)
)
GO

CREATE TABLE [test_complex] (
  [score] NUMBER(3) CHECK (score >= 0 AND MOD(score, 2) = 0),
  [rating] VARCHAR2(20) CHECK (rating IN ('HIGH', 'MEDIUM', 'LOW') AND LENGTH(rating) > 3)
)
GO

CREATE TABLE [test_pattern_match] (
  [zip_code] VARCHAR2(10) CHECK (REGEXP_LIKE(zip_code, '^[0-9]{5}(-[0-9]{4})?$')),
  [serial_number] VARCHAR2(20) CHECK (REGEXP_LIKE(serial_number, '^[A-Z0-9]+$'))
)
GO

CREATE TABLE [test_subquery] (
  [dept_id] NUMBER(10) CHECK (dept_id IN (SELECT dept_id FROM departments)),
  [emp_count] NUMBER(5) CHECK (emp_count >= 0)
)
GO
