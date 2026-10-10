CREATE TABLE [data] (
  [price] decimal(10,2),
  [rate] float,
  [amount] numeric
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [data] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [data] ([price], [rate], [amount])
VALUES
  (99.99, 3.14159, 0.001),
  (50.5, 0.5, 100);
GO

-- Re-enable constraint checks
ALTER TABLE [data] WITH CHECK CHECK CONSTRAINT ALL;
GO