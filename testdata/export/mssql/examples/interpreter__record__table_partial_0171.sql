CREATE TABLE [T] (
  [s1] int,
  [s2] int,
  [a] int PRIMARY KEY,
  [b] int,
  [m1] int,
  [m2] int,
  [c] int,
  [d] int,
  [e1] int,
  [e2] int
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [T] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [T] ([s1], [s2], [a], [b], [m1], [m2], [c], [d], [e1], [e2])
VALUES
  (1, 2, 3, 4, 5, 6, 7, 8, 9, 10);
GO

-- Re-enable constraint checks
ALTER TABLE [T] WITH CHECK CHECK CONSTRAINT ALL;
GO