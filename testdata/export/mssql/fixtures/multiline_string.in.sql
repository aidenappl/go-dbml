CREATE TABLE [users] (
  [id] int
)
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = '# Objective
  * Support writing long string that can ''span'' over multiple lines
  * Support writing markdown for DBML Note

# Syntax

',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'users',
@level2type = N'Column', @level2name = 'id';
GO
