CREATE SCHEMA [SchemaName]
GO

CREATE SCHEMA [Schema Name]
GO

CREATE TABLE [accounts] (
  [id] int,
  [name] Var Char,
  [role] nvarchar(255) CHECK ([role] IN ('User', 'Admin')),
  [role3] nvarchar(255) CHECK ([role3] IN ('User', 'Admin')),
  [role4] nvarchar(255) CHECK ([role4] IN ('User', 'Admin')),
  [role5] nvarchar(255) CHECK ([role5] IN ('User', 'Admin'))
)
GO
