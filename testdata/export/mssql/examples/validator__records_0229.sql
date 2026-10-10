CREATE TABLE [ids] (
  [id] int PRIMARY KEY
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [ids] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [ids] ([id])
VALUES
  (1),
  (2),
  (3);
GO

-- Re-enable constraint checks
ALTER TABLE [ids] WITH CHECK CHECK CONSTRAINT ALL;
GO