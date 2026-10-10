CREATE TABLE [users] (
  [id] int
)
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Multi-line
note content
',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'users';
GO
