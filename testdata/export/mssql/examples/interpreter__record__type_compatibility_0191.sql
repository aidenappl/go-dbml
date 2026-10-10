CREATE TABLE [users] (
  [id] int PRIMARY KEY IDENTITY(1, 1),
  [status] nvarchar(255) NOT NULL DEFAULT 'active'
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [users] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [users] ([id], [status])
VALUES
  (NULL, NULL),
  (NULL, 'inactive');
GO

-- Re-enable constraint checks
ALTER TABLE [users] WITH CHECK CHECK CONSTRAINT ALL;
GO