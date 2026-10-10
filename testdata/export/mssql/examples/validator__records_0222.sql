CREATE TABLE [data] (
  [id] int,
  [value] int
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [data] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [data] ([id], [value])
VALUES
  (1, -100),
  (2, -999);
GO

-- Re-enable constraint checks
ALTER TABLE [data] WITH CHECK CHECK CONSTRAINT ALL;
GO