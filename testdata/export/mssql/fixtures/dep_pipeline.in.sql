CREATE TABLE [raw_orders] (
  [id] int PRIMARY KEY,
  [user_id] int,
  [amount] decimal,
  [created_at] timestamp
)
GO

CREATE TABLE [stg_orders] (
  [id] int PRIMARY KEY,
  [user_id] int,
  [amount] decimal
)
GO

CREATE TABLE [fct_orders] (
  [id] int PRIMARY KEY,
  [user_id] int,
  [revenue] decimal
)
GO

CREATE TABLE [mart_orders] (
  [id] int PRIMARY KEY,
  [total] decimal
)
GO
