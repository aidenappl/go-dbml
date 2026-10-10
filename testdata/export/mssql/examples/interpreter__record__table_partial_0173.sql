CREATE TABLE [T] (
  [a] int,
  [b] int,
  [x] int PRIMARY KEY,
  [shared] nvarchar(255),
  [c] int,
  [y] int
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [T] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [T] ([a], [b], [x], [shared], [c], [y])
VALUES
  (1, 2, 3, 'shared_value', 4, 5);
GO

-- Re-enable constraint checks
ALTER TABLE [T] WITH CHECK CHECK CONSTRAINT ALL;
GO