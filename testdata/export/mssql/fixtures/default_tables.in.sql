CREATE TABLE [orders] (
  [id] int PRIMARY KEY DEFAULT (123),
  [user_id] int UNIQUE NOT NULL,
  [status] nvarchar(255) DEFAULT 'Completed',
  [created_at] nvarchar(255) DEFAULT (now())
)
GO

CREATE TABLE [order_items] (
  [order_id] int,
  [product_id] int,
  [quantity] int
)
GO

CREATE TABLE [products] (
  [id] int PRIMARY KEY,
  [name] nvarchar(255),
  [merchant_id] int NOT NULL DEFAULT (-1),
  [price] float DEFAULT (-123.12),
  [stock] boolean DEFAULT (true),
  [expiration] date DEFAULT (current_date + interval 1 year)
)
GO
