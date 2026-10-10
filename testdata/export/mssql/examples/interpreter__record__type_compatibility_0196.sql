CREATE SCHEMA [auth]
GO

CREATE SCHEMA [app]
GO

CREATE TABLE [auth].[users] (
  [id] int PRIMARY KEY,
  [role] nvarchar(255) CHECK ([role] IN ('admin', 'user')),
  [status] nvarchar(255) CHECK ([status] IN ('active', 'inactive'))
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [auth].[users] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [auth].[users] ([id], [role], [status])
VALUES
  (1, , CAST('active' AS status)),
  (2, , );
GO

-- Re-enable constraint checks
ALTER TABLE [auth].[users] WITH CHECK CHECK CONSTRAINT ALL;
GO