CREATE TABLE [test_datatypes] (
  [col_number] NUMBER(10,2),
  [col_varchar2] VARCHAR2(100),
  [col_char] CHAR(10),
  [col_date] DATE,
  [col_timestamp] TIMESTAMP,
  [col_clob] CLOB,
  [col_raw] RAW(100),
  [col_float] FLOAT,
  [col_integer] INTEGER,
  [col_nvarchar2] NVARCHAR2(100),
  [col_nchar] NCHAR(10),
  PRIMARY KEY ([col_number]),
  PRIMARY KEY ([col_varchar2]),
  PRIMARY KEY ([col_char]),
  PRIMARY KEY ([col_date]),
  PRIMARY KEY ([col_timestamp]),
  PRIMARY KEY ([col_clob]),
  PRIMARY KEY ([col_raw]),
  PRIMARY KEY ([col_float]),
  PRIMARY KEY ([col_integer]),
  PRIMARY KEY ([col_nvarchar2]),
  PRIMARY KEY ([col_nchar])
)
GO

CREATE TABLE [order_details] (
  [order_id] NUMBER(10),
  [item_id] NUMBER(10),
  [quantity] NUMBER(10),
  [order_date] DATE,
  PRIMARY KEY ([order_id], [item_id])
)
GO

CREATE TABLE [reports] (
  [report_id] NUMBER(10),
  [report_name] VARCHAR2(100),
  [report_date] DATE,
  [status] VARCHAR2(20),
  PRIMARY KEY ([report_id])
)
GO

CREATE TABLE [sales] (
  [sale_id] NUMBER(10),
  [sale_date] DATE,
  [amount] NUMBER(12,2),
  PRIMARY KEY ([sale_id])
)
GO

CREATE TABLE [accounts] (
  [account_id] NUMBER(10),
  [account_number] VARCHAR2(20),
  [balance] NUMBER(12,2),
  PRIMARY KEY ([account_id])
)
GO
