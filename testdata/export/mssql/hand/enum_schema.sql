CREATE SCHEMA [v2]
GO

CREATE TABLE [v2].[jobs] (
  [status] nvarchar(255) CHECK ([status] IN ('created', 'running', 'done', 'failure')) DEFAULT 'created'
)
GO
