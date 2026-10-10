CREATE TABLE [products] (
  [id] int PRIMARY KEY DEFAULT (123),
  [name] nvarchar(255) DEFAULT 'Tea',
  [merchant_id] int NOT NULL,
  [price] float DEFAULT (123.12),
  [status] nvarchar(255),
  [created_at] nvarchar(255) DEFAULT (now()),
  [stock] boolean DEFAULT (true),
  [expiration] date DEFAULT (current_date + interval 1 year)
)
GO
