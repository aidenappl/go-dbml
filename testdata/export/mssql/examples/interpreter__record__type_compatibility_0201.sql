CREATE TABLE [users] (
  [id] int,
  [name] varchar(5),
  [quantity] int
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [users] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [users] ([id], [name], [quantity])
VALUES
  (1, 'Alice', 10),
  (2, 'Bob', 20),
  (3, 'Christopher', 30),
  (4, 'Dave', 40),
  (5, 'Elizabeth', 50);
GO

-- Re-enable constraint checks
ALTER TABLE [users] WITH CHECK CHECK CONSTRAINT ALL;
GO