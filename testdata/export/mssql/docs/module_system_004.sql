CREATE SCHEMA [auth]
GO

CREATE TABLE [auth].[users] (
  [id] int PRIMARY KEY,
  [email] nvarchar(255)
)
GO

CREATE TABLE [auth].[roles] (
  [id] int PRIMARY KEY,
  [name] nvarchar(255)
)
GO

CREATE TABLE [auth].[sessions] (
  [id] int PRIMARY KEY
)
GO
