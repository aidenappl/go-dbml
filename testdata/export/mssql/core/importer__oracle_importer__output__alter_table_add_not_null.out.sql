CREATE TABLE [customers] (
  [cust_id] NUMBER(10),
  [cust_name] VARCHAR2(100) NOT NULL,
  [email] VARCHAR2(100) NOT NULL,
  [phone] VARCHAR2(20)
)
GO

CREATE TABLE [inventory] (
  [item_id] NUMBER(10),
  [item_name] VARCHAR2(200) NOT NULL,
  [price] NUMBER(10,2) NOT NULL,
  [category] VARCHAR2(50) NOT NULL,
  [manufacturer] VARCHAR2(100) NOT NULL,
  [sku] VARCHAR2(50) NOT NULL
)
GO

CREATE TABLE [shipments] (
  [ship_id] NUMBER(10),
  [customer_id] NUMBER(10) NOT NULL,
  [ship_date] DATE NOT NULL,
  [total_amount] NUMBER(12,2) NOT NULL
)
GO

CREATE TABLE [departments] (
  [dept_id] NUMBER(10) NOT NULL,
  [dept_name] VARCHAR2(100) NOT NULL,
  [manager_id] NUMBER(10),
  [status] VARCHAR2(20)
)
GO

CREATE TABLE [suppliers] (
  [supp_id] NUMBER(10),
  [supp_name] VARCHAR2(100) NOT NULL,
  [contact] VARCHAR2(100),
  [ssn] VARCHAR2(11),
  [badge_number] VARCHAR2(10)
)
GO
