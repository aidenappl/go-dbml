CREATE TABLE [orders] (
  [status] nvarchar(255)
)
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = '💸 1 = processing,
✔️ 2 = shipped,
❌ 3 = cancelled,
😔 4 = refunded
',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'orders',
@level2type = N'Column', @level2name = 'status';
GO
