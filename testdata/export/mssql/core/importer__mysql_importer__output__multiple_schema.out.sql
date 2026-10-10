CREATE SCHEMA [ecommerce]
GO

CREATE SCHEMA [schemaA]
GO

CREATE SCHEMA [schemaB]
GO

CREATE TABLE [users] (
  [id] int PRIMARY KEY,
  [name] varchar(255),
  [pjs] nvarchar(255) CHECK ([pjs] IN ('created2', 'running2', 'done2', 'failure2')),
  [pjs2] nvarchar(255) CHECK ([pjs2] IN ('created2', 'running2', 'done2', 'failure2')),
  [pg] nvarchar(255) CHECK ([pg] IN ('male', 'female')),
  [pg2] nvarchar(255) CHECK ([pg2] IN ('male2', 'female2'))
)
GO

CREATE TABLE [products] (
  [id] int PRIMARY KEY,
  [name] varchar(255)
)
GO

CREATE TABLE [orders] (
  [id1] int,
  [id2] int,
  [pid] int,
  [name] varchar(255),
  [name2] varchar(255),
  PRIMARY KEY ([id1], [id2])
)
GO

CREATE TABLE [orders2] (
  [id1] int,
  [id2] int,
  [pid] int,
  PRIMARY KEY ([id1], [id2])
)
GO

CREATE TABLE [orders4] (
  [id1] int,
  [id2] int,
  PRIMARY KEY ([id1], [id2])
)
GO

CREATE TABLE [ecommerce].[users] (
  [id] int PRIMARY KEY,
  [name] varchar(255),
  [ejs] nvarchar(255) CHECK ([ejs] IN ('created2', 'running2', 'done2', 'failure2')),
  [ejs2] nvarchar(255) CHECK ([ejs2] IN ('created2', 'running2', 'done2', 'failure2')),
  [eg] nvarchar(255) CHECK ([eg] IN ('male', 'female')),
  [eg2] nvarchar(255) CHECK ([eg2] IN ('male2', 'female2'))
)
GO

CREATE TABLE [schemaA].[products] (
  [id] int PRIMARY KEY,
  [name] varchar(255),
  [lid] int
)
GO

CREATE TABLE [schemaA].[locations] (
  [id] int PRIMARY KEY,
  [name] varchar(255)
)
GO

CREATE TABLE [schemaB].[orders3] (
  [id1] int PRIMARY KEY,
  [id2] int
)
GO

CREATE UNIQUE INDEX [I_UNIQUE_1] ON [orders] ("name", "name2")
GO

CREATE INDEX [INDEX_2] ON [orders] ("id1", "id2")
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
@value = 'this is table orders',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'orders';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Sample field comment on multiples schema',
@level0type = N'Schema', @level0name = 'schemaA',
@level1type = N'Table',  @level1name = 'products',
@level2type = N'Column', @level2name = 'name';
GO

ALTER TABLE [schemaA].[products] ADD CONSTRAINT [FK_1] FOREIGN KEY ([lid]) REFERENCES [schemaA].[locations] ([id]) ON DELETE CASCADE
GO

ALTER TABLE [orders] ADD FOREIGN KEY ([pid]) REFERENCES [schemaA].[products] ([id])
GO

ALTER TABLE [orders2] ADD FOREIGN KEY ([pid]) REFERENCES [schemaA].[products] ([id])
GO

ALTER TABLE [orders2] ADD CONSTRAINT [CFK_1] FOREIGN KEY ([id1], [id2]) REFERENCES [orders] ([id1], [id2]) ON DELETE CASCADE ON UPDATE SET NULL
GO

ALTER TABLE [ecommerce].[users] ADD FOREIGN KEY ([id]) REFERENCES [users] ([id])
GO

ALTER TABLE [ecommerce].[users] ADD CONSTRAINT [name_optional] FOREIGN KEY ([id]) REFERENCES [users] ([name])
GO

ALTER TABLE [schemaA].[products] ADD FOREIGN KEY ([name]) REFERENCES [ecommerce].[users] ([id])
GO

ALTER TABLE [schemaA].[locations] ADD FOREIGN KEY ([name]) REFERENCES [users] ([id])
GO
