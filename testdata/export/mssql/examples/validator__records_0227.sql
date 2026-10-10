CREATE TABLE [users] (
  [user-id] int PRIMARY KEY,
  [user-name] nvarchar(255)
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [users] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [users] ([user-id], [user-name])
VALUES
  (1, 'Alice');
GO

-- Re-enable constraint checks
ALTER TABLE [users] WITH CHECK CHECK CONSTRAINT ALL;
GO