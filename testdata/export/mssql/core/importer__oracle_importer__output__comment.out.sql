CREATE TABLE [customers] (
  [name] VARCHAR(20)
)
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'This table stores customer data.',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'customers';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'The name of the customer.',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'customers',
@level2type = N'Column', @level2name = 'name';
GO
