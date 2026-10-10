CREATE SCHEMA [auth]
GO

CREATE TABLE [auth].[users] (
  [id] int PRIMARY KEY
)
GO

CREATE TABLE [posts] (
  [user_id] int
)
GO

ALTER TABLE [posts] ADD FOREIGN KEY ([user_id]) REFERENCES [auth].[users] ([id])
GO
