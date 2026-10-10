CREATE TABLE [table1] (
  [id] int PRIMARY KEY IDENTITY(1, 1),
  [payload] text
)
GO

CREATE TABLE [table2] (
  [id] integer PRIMARY KEY IDENTITY(1, 1),
  [name] nvarchar(255)
)
GO

CREATE TABLE [table3] (
  [id] smallint PRIMARY KEY IDENTITY(1, 1),
  [name] nvarchar(255)
)
GO

CREATE TABLE [table4] (
  [id] bigint PRIMARY KEY IDENTITY(1, 1),
  [name] nvarchar(255)
)
GO

CREATE TABLE [table5] (
  [id] serial PRIMARY KEY IDENTITY(1, 1),
  [name] nvarchar(255)
)
GO

CREATE TABLE [table6] (
  [id] smallserial PRIMARY KEY IDENTITY(1, 1),
  [name] nvarchar(255)
)
GO

CREATE TABLE [table7] (
  [id] bigserial PRIMARY KEY IDENTITY(1, 1),
  [name] nvarchar(255)
)
GO

CREATE TABLE [table8] (
  [id] intaaaa PRIMARY KEY IDENTITY(1, 1),
  [name] nvarchar(255)
)
GO
