CREATE SCHEMA [types]
GO

CREATE TABLE [users] (
  [status] nvarchar(255) CHECK ([status] IN ('active', 'inactive')) DEFAULT 'active'
)
GO
