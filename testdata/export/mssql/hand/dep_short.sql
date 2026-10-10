CREATE TABLE [raw_orders] (
  [id] int,
  [amount] int
)
GO

CREATE TABLE [stg_orders] (
  [id] int,
  [revenue] int
)
GO

CREATE TABLE [fct_orders] (
  [revenue] int
)
GO
