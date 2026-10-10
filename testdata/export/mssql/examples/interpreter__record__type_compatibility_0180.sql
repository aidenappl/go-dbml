CREATE TABLE [data] (
  [id] int,
  [quantity] int,
  [price] decimal
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [data] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [data] ([id], [quantity], [price])
VALUES
  (not_a_number, 100, 100),
  (2, 10, 100),
  (3, 20, also_not_a_number),
  (4, 30, true);
GO

-- Re-enable constraint checks
ALTER TABLE [data] WITH CHECK CHECK CONSTRAINT ALL;
GO