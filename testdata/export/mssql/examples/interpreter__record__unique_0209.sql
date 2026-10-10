CREATE TABLE [users] (
  [id] int PRIMARY KEY,
  [email] nvarchar(255) UNIQUE
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [users] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [users] ([id], [email])
VALUES
  (1, 'alice@example.com'),
  (1, 'bob@example.com');
GO

-- Re-enable constraint checks
ALTER TABLE [users] WITH CHECK CHECK CONSTRAINT ALL;
GO