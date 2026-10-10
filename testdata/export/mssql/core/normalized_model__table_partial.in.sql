CREATE SCHEMA [a]
GO

CREATE TABLE [a] (
  [id] int,
  [id2] int
)
GO

CREATE TABLE [users] (
  [created_at] timestamp DEFAULT (now()),
  [updated_at] timestamp DEFAULT (now()),
  [name] nvarchar(255),
  [email] nvarchar(255),
  [id] int,
  [delete_status] boolean NOT NULL,
  [deleted_at] timestamp DEFAULT (now()),
  [refColumn1] bit,
  [refColumn2] bool,
  PRIMARY KEY ([id])
)
GO

CREATE TABLE [testColWithTypeHaveArgs] (
  [id] int,
  [callExpColInTest1] int(4),
  [callExpColInTest2] decimal(10,2) DEFAULT (0.1)
)
GO

CREATE TABLE [table_no_column_with_partial] (
  [email] nvarchar(255) UNIQUE
)
GO

CREATE TABLE [a].[a] (
  [id] int,
  [id2] int
)
GO

CREATE UNIQUE INDEX [users_index_0] ON [users] ("email")
GO

CREATE UNIQUE INDEX [users_index_1] ON [users] ("email", "id")
GO

CREATE UNIQUE INDEX [users_index_2] ON [users] ("email")
GO

CREATE UNIQUE INDEX [table_no_column_with_partial_index_3] ON [table_no_column_with_partial] ("email")
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'simple column with call expression as type',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'testColWithTypeHaveArgs',
@level2type = N'Column', @level2name = 'callExpColInTest1';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'column with type that have arguments and default value being a float number',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'testColWithTypeHaveArgs',
@level2type = N'Column', @level2name = 'callExpColInTest2';
GO

ALTER TABLE [users] ADD FOREIGN KEY ([refColumn1]) REFERENCES [a] ([id2])
GO

ALTER TABLE [users] ADD FOREIGN KEY ([refColumn2]) REFERENCES [a].[a] ([id])
GO
