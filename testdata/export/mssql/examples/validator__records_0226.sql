CREATE SCHEMA [auth]
GO

CREATE TABLE [auth].[users] (
  [id] int PRIMARY KEY,
  [email] nvarchar(255)
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [auth].[users] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [auth].[users] ([id], [email])
VALUES
  (1, 'alice@example.com');
GO

-- Re-enable constraint checks
ALTER TABLE [auth].[users] WITH CHECK CHECK CONSTRAINT ALL;
GO