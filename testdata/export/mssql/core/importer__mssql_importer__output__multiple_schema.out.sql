CREATE SCHEMA [ecommerce]
GO

CREATE SCHEMA [schemaA]
GO

CREATE TABLE [users] (
  [id] int PRIMARY KEY,
  [name] nvarchar(255),
  [pjs] nvarchar(255) NOT NULL CHECK ([pjs] IN ('created2', 'running2', 'done2', 'failure2')),
  [pjs2] nvarchar(255) NOT NULL CHECK ([pjs2] IN ('created2', 'running2', 'done2', 'failure2')),
  [pg] nvarchar(255) NOT NULL CHECK ([pg] IN ('male', 'female')),
  [pg2] nvarchar(255) NOT NULL CHECK ([pg2] IN ('male2', 'female2'))
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
  [ejs] nvarchar(255) NOT NULL CHECK ([ejs] IN ('created2', 'running2', 'done2', 'failure2')),
  [ejs2] nvarchar(255) NOT NULL CHECK ([ejs2] IN ('created2', 'running2', 'done2', 'failure2')),
  [eg] nvarchar(255) NOT NULL CHECK ([eg] IN ('male', 'female')),
  [eg2] nvarchar(255) NOT NULL CHECK ([eg2] IN ('male2', 'female2'))
)
GO

CREATE TABLE [schemaA].[products] (
  [id] int PRIMARY KEY,
  [name] nvarchar(255),
  [created_at] varchar(255) DEFAULT (now()),
  [lid] int,
  [lid2] int
)
GO

CREATE TABLE [schemaA].[locations] (
  [id] int PRIMARY KEY,
  [name] nvarchar(255)
)
GO

CREATE INDEX [idx_1] ON [ecommerce].[users] ("name", "ejs")
GO

CREATE UNIQUE INDEX [unique_lid_lid2] ON [schemaA].[products] ("lid", "lid2")
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Product name',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'products',
@level2type = N'Column', @level2name = 'name';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Note on table users of schema ecommerce',
@level0type = N'Schema', @level0name = 'ecommerce',
@level1type = N'Table',  @level1name = 'users';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'This is a note in table "schemaA"."locations"',
@level0type = N'Schema', @level0name = 'schemaA',
@level1type = N'Table',  @level1name = 'locations';
GO

ALTER TABLE [schemaA].[products] ADD FOREIGN KEY ([lid]) REFERENCES [schemaA].[locations] ([id])
GO

ALTER TABLE [schemaA].[products] ADD CONSTRAINT [FK_1] FOREIGN KEY ([lid2]) REFERENCES [schemaA].[locations] ([id])
GO

ALTER TABLE [ecommerce].[users] ADD FOREIGN KEY ([id]) REFERENCES [users] ([id])
GO

ALTER TABLE [ecommerce].[users] ADD CONSTRAINT [name_optional] FOREIGN KEY ([id]) REFERENCES [users] ([name])
GO

ALTER TABLE [schemaA].[products] ADD FOREIGN KEY ([name]) REFERENCES [ecommerce].[users] ([id])
GO

ALTER TABLE [schemaA].[locations] ADD FOREIGN KEY ([name]) REFERENCES [users] ([id])
GO
