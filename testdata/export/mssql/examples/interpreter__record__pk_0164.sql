CREATE TABLE [orders] (
  [id] int DEFAULT (0),
  [code] int,
  PRIMARY KEY ([id], [code])
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [orders] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [orders] ([code])
VALUES
  (1),
  (1);
GO

-- Re-enable constraint checks
ALTER TABLE [orders] WITH CHECK CHECK CONSTRAINT ALL;
GO