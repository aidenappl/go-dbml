CREATE SCHEMA [types]
GO

CREATE TABLE [users] (
  [id] int,
  [status] nvarchar(255) CHECK ([status] IN ('active', 'inactive'))
)
GO
