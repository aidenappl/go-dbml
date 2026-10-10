CREATE SCHEMA [sales]
GO

CREATE TABLE [users] (
  [id] int PRIMARY KEY,
  [name] nvarchar(255)
)
GO

CREATE TABLE [posts] (
  [id] int PRIMARY KEY,
  [user_id] int
)
GO

CREATE TABLE [sales].[orders] (
  [id] int PRIMARY KEY
)
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'this will override',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'users';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = '# Heading 1
    code block
# Heading 2
',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'posts';
GO
