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

CREATE TABLE [b] (
  [id] int,
  [c_id] int
)
GO

CREATE TABLE [c] (
  [id] int,
  [b_id] int
)
GO

CREATE TABLE [unary_expression] (
  [id] int(-1) DEFAULT (-2),
  [id2] int(--1) DEFAULT (-2),
  [id3] int(+-+---+0.1) DEFAULT (7.2225)
)
GO

CREATE TABLE [injected_unary_expression] (
  [id] int,
  [id2] int(--1) DEFAULT (-2),
  [id3] int(+-+---+0.1) DEFAULT (-7.2225)
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
  [name] nvarchar(255)
)
GO

CREATE TABLE [schemaA].[locations] (
  [id] int PRIMARY KEY,
  [name] nvarchar(255)
)
GO

ALTER TABLE [schemaA].[products] ADD FOREIGN KEY ([name]) REFERENCES [ecommerce].[users] ([id])
GO

ALTER TABLE [schemaA].[locations] ADD FOREIGN KEY ([name]) REFERENCES [users] ([id])
GO

ALTER TABLE [ecommerce].[users] ADD FOREIGN KEY ([id]) REFERENCES [users] ([id])
GO

ALTER TABLE [ecommerce].[users] ADD CONSTRAINT [name_optional] FOREIGN KEY ([id]) REFERENCES [users] ([name])
GO

ALTER TABLE [c] ADD CONSTRAINT [short_ref] FOREIGN KEY ([id]) REFERENCES [b] ([id])
GO

ALTER TABLE [b] ADD CONSTRAINT [long_ref] FOREIGN KEY ([c_id]) REFERENCES [c] ([id])
GO
