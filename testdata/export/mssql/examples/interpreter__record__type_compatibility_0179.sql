CREATE TABLE [data] (
  [id] int,
  [value] decimal,
  [count] int,
  [temperature] float
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [data] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [data] ([id], [value], [count], [temperature])
VALUES
  (1, 10000000000, 100, 314),
  (2, 0.0000314, 2000, -0.0025),
  (3, 200000000, NULL, NULL);
GO

-- Re-enable constraint checks
ALTER TABLE [data] WITH CHECK CHECK CONSTRAINT ALL;
GO