CREATE TABLE [customers] (
  [cust_id] NUMBER IDENTITY(1, 1),
  [cust_name] VARCHAR2(100) DEFAULT 'Unknown',
  [email] VARCHAR2(100) DEFAULT 'noemail@example.com',
  [phone] VARCHAR2(20) DEFAULT 'N/A'
)
GO

CREATE TABLE [inventory] (
  [item_id] NUMBER IDENTITY(1, 1),
  [item_name] VARCHAR2(200) DEFAULT 'Default Item',
  [price] NUMBER(10,2) DEFAULT (0),
  [category] VARCHAR2(50) DEFAULT 'Misc'
)
GO

CREATE TABLE [shipments] (
  [ship_id] NUMBER IDENTITY(1, 1),
  [customer_id] NUMBER(10) DEFAULT (1),
  [ship_date] DATE DEFAULT (SYSDATE),
  [total_amount] NUMBER(12,2) DEFAULT (0)
)
GO

CREATE TABLE [departments] (
  [dept_id] NUMBER IDENTITY(1, 1),
  [dept_name] VARCHAR2(100) DEFAULT 'General',
  [manager_id] NUMBER(10) DEFAULT (0),
  [budget] NUMBER(12,2) DEFAULT (10000)
)
GO

CREATE TABLE [suppliers] (
  [supp_id] NUMBER IDENTITY(1, 1),
  [supp_name] VARCHAR2(100) DEFAULT 'Anonymous',
  [contact] VARCHAR2(100) DEFAULT 'Contact',
  [credit_limit] NUMBER(10,2) DEFAULT (5000)
)
GO

CREATE TABLE [test_datatypes] (
  [col_number] NUMBER IDENTITY(1, 1),
  [col_varchar2] VARCHAR2(100) DEFAULT 'Default Text',
  [col_char] CHAR(10) DEFAULT 'DEFAULT   ',
  [col_date] DATE DEFAULT (SYSDATE)
)
GO

CREATE TABLE [test_basic] (
  [id] NUMBER IDENTITY(1, 1),
  [value] NUMBER(10) DEFAULT (100)
)
GO

CREATE TABLE [test_status] (
  [status_id] NUMBER IDENTITY(1, 1),
  [status] VARCHAR2(20) DEFAULT 'OPEN'
)
GO

CREATE TABLE [test_not_null] (
  [record_id] NUMBER NOT NULL IDENTITY(1, 1),
  [description] VARCHAR2(200) NOT NULL DEFAULT 'No description'
)
GO

CREATE TABLE [test_start_with] (
  [order_id] NUMBER IDENTITY(1, 1),
  [order_code] VARCHAR2(20) DEFAULT 'ORD000'
)
GO

CREATE TABLE [test_increment] (
  [ticket_id] NUMBER IDENTITY(1, 1),
  [ticket_type] VARCHAR2(20) DEFAULT 'STANDARD'
)
GO

CREATE TABLE [test_maxvalue] (
  [ref_id] NUMBER IDENTITY(1, 1),
  [ref_name] VARCHAR2(100) DEFAULT 'Reference'
)
GO

CREATE TABLE [test_minvalue] (
  [batch_id] NUMBER IDENTITY(1, 1),
  [batch_name] VARCHAR2(50) DEFAULT 'Batch'
)
GO
