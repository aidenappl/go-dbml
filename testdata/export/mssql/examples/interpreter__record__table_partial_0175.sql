CREATE TABLE [T] (
  [a] nvarchar(255) PRIMARY KEY,
  [b] text,
  [c] int
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [T] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [T] ([a], [b], [c])
VALUES
  ('value_a', 'value_b', 3);
GO

-- Re-enable constraint checks
ALTER TABLE [T] WITH CHECK CHECK CONSTRAINT ALL;
GO