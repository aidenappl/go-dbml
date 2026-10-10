CREATE TABLE [data] (
  [id] int,
  [value] decimal
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [data] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [data] ([id], [value])
VALUES
  (1, 10000000000),
  (2, 0.0000314),
  (3, 200000000);
GO

-- Re-enable constraint checks
ALTER TABLE [data] WITH CHECK CHECK CONSTRAINT ALL;
GO