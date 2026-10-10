CREATE TABLE [T] (
  [a] int PRIMARY KEY,
  [b] int
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [T] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [T] ([a], [b])
VALUES
  (1, 2);
GO

-- Re-enable constraint checks
ALTER TABLE [T] WITH CHECK CHECK CONSTRAINT ALL;
GO