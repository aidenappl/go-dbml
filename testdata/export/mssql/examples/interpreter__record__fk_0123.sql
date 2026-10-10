CREATE SCHEMA [auth]
GO

CREATE TABLE [auth].[users] (
  [id] int,
  [tenant_id] int,
  PRIMARY KEY ([id], [tenant_id])
)
GO

CREATE TABLE [posts] (
  [id] int PRIMARY KEY,
  [user_id] int,
  [tenant_id] int,
  [content] text
)
GO

ALTER TABLE [posts] ADD FOREIGN KEY ([user_id], [tenant_id]) REFERENCES [auth].[users] ([id], [tenant_id])
GO

-- Disable constraint checks for tables with data
ALTER TABLE [auth].[users] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [posts] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [auth].[users] ([id], [tenant_id])
VALUES
  (1, 100),
  (2, 100);
GO
INSERT INTO [posts] ([id], [user_id], [tenant_id], [content])
VALUES
  (1, 1, 100, 'Hello'),
  (2, 999, 100, 'Invalid user');
GO

-- Re-enable constraint checks
ALTER TABLE [auth].[users] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [posts] WITH CHECK CHECK CONSTRAINT ALL;
GO