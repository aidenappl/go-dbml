CREATE TABLE [users] (
  [id] integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  [username] nvarchar2(255) NOT NULL,
  [fullname] nvarchar2(255) NOT NULL
)
GO

CREATE TABLE [orders] (
  [id] integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  [status] nvarchar(255) CHECK ([status] IN ('created', 'pending', 'processing', 'waiting_for_payment', 'success', 'cancel')),
  [note] nclob DEFAULT '',
  [user_id] integer NOT NULL,
  [created_at] timestamp DEFAULT (current_timestamp),
  [product_id] int,
  [supplier_id] int
)
GO

CREATE TABLE [order_details] (
  [order_id] INT UNIQUE NOT NULL,
  [unit] nvarchar2(50) NOT NULL,
  [quantity] number NOT NULL,
  [price] float NOT NULL
)
GO

CREATE TABLE [products] (
  [product_id] int,
  [supplier_id] int,
  PRIMARY KEY ([product_id], [supplier_id])
)
GO

CREATE INDEX [IDX_USERS] ON [users] ((LOWER("fullname")), "username")
GO

CREATE UNIQUE INDEX [IDX_ORDER_USER] ON [orders] ("id", "user_id")
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'This table has a expression index',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'users';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Store user uuid',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'orders',
@level2type = N'Column', @level2name = 'user_id';
GO

ALTER TABLE [orders] ADD FOREIGN KEY ([product_id], [supplier_id]) REFERENCES [products] ([product_id], [supplier_id])
GO

ALTER TABLE [orders] ADD FOREIGN KEY ([user_id]) REFERENCES [users] ([id]) ON DELETE CASCADE
GO

ALTER TABLE [order_details] ADD FOREIGN KEY ([order_id]) REFERENCES [orders] ([id]) ON DELETE CASCADE ON UPDATE RESTRICT
GO
