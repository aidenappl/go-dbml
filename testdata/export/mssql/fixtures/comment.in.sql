CREATE TABLE [orders] (
  [id] int PRIMARY KEY,
  [user_id] int UNIQUE NOT NULL,
  [status] nvarchar(255),
  [created_at] nvarchar(255)
)
GO

CREATE INDEX [orders_index_0] ON [orders] ("id")
GO

CREATE INDEX [orders_index_1] ON [orders] ("created_at")
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Status of an order',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'orders',
@level2type = N'Column', @level2name = 'status';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'When order created',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'orders',
@level2type = N'Column', @level2name = 'created_at';
GO
