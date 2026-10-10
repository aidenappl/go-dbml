CREATE TABLE [t] (
  [c] int
)
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'table-level note
',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 't';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'it''s escaped 
 newline',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 't',
@level2type = N'Column', @level2name = 'c';
GO
