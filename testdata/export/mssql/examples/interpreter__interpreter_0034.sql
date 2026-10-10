CREATE TABLE [users] (
  [id] int UNIQUE PRIMARY KEY NOT NULL IDENTITY(1, 1)
)
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'ID',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'users',
@level2type = N'Column', @level2name = 'id';
GO
