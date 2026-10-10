CREATE SCHEMA [schemaB]
GO

CREATE SCHEMA [ecommerce]
GO

CREATE SCHEMA [schemaA]
GO

CREATE TABLE [users] (
  [id] int PRIMARY KEY,
  [name] nvarchar(255),
  [pjs] nvarchar(255) CHECK ([pjs] IN ('created2', 'running2', 'done2', 'failure2')),
  [pjs2] nvarchar(255) CHECK ([pjs2] IN ('created2', 'running2', 'done2', 'failure2')),
  [pg] nvarchar(255) CHECK ([pg] IN ('male', 'female')),
  [pg2] nvarchar(255) CHECK ([pg2] IN ('male2', 'female2'))
)
GO

CREATE TABLE [products] (
  [id] int PRIMARY KEY,
  [name] nvarchar(255)
)
GO

CREATE TABLE [ecommerce].[users] (
  [id] int PRIMARY KEY,
  [name] nvarchar(255),
  [ejs] nvarchar(255) CHECK ([ejs] IN ('created2', 'running2', 'done2', 'failure2')),
  [ejs2] nvarchar(255) CHECK ([ejs2] IN ('created2', 'running2', 'done2', 'failure2')),
  [eg] nvarchar(255) CHECK ([eg] IN ('male', 'female')),
  [eg2] nvarchar(255) CHECK ([eg2] IN ('male2', 'female2'))
)
GO

CREATE TABLE [schemaA].[products] (
  [id] int PRIMARY KEY,
  [name] nvarchar(255),
  [lid] int
)
GO

CREATE TABLE [schemaA].[locations] (
  [id] int PRIMARY KEY,
  [name] nvarchar(255)
)
GO

CREATE UNIQUE INDEX [products_index_0] ON [products] ("id")
GO

CREATE INDEX [product_status] ON [schemaA].[products] ("id", "name")
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Product name of public schema',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'products',
@level2type = N'Column', @level2name = 'name';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Product name of schemaA',
@level0type = N'Schema', @level0name = 'schemaA',
@level1type = N'Table',  @level1name = 'products',
@level2type = N'Column', @level2name = 'name';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'This is a note in table "locations"',
@level0type = N'Schema', @level0name = 'schemaA',
@level1type = N'Table',  @level1name = 'locations';
GO

ALTER TABLE [schemaA].[products] ADD CONSTRAINT [fk_1] FOREIGN KEY ([lid]) REFERENCES [schemaA].[locations] ([id])
GO

ALTER TABLE [ecommerce].[users] ADD FOREIGN KEY ([id]) REFERENCES [users] ([id])
GO

ALTER TABLE [ecommerce].[users] ADD CONSTRAINT [name_optional] FOREIGN KEY ([id]) REFERENCES [users] ([name])
GO

ALTER TABLE [schemaA].[products] ADD FOREIGN KEY ([name]) REFERENCES [ecommerce].[users] ([id])
GO

ALTER TABLE [schemaA].[locations] ADD FOREIGN KEY ([name]) REFERENCES [users] ([id])
GO
