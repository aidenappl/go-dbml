CREATE TABLE [customers] (
  [customer_id] INT,
  [first_name] VARCHAR(50),
  [last_name] VARCHAR(50),
  [email] VARCHAR(100),
  [registration_date] DATE
)
GO

CREATE TABLE [orders] (
  [order_id] INT,
  [customer_id] INT,
  [order_date] DATE,
  [total_amount] DECIMAL(10,2)
)
GO

CREATE TABLE [employees] (
  [employee_id] INT,
  [first_name] VARCHAR(50),
  [last_name] VARCHAR(50),
  [hire_date] DATE,
  [department] VARCHAR(50)
)
GO

CREATE UNIQUE INDEX [customers_index_0] ON [customers] ("email")
GO

ALTER TABLE [orders] ADD FOREIGN KEY ([customer_id]) REFERENCES [customers] ([customer_id])
GO
