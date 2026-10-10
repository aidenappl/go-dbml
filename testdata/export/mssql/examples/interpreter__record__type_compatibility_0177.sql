CREATE TABLE [data] (
  [id] int,
  [active] boolean
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [data] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [data] ([id], [active])
VALUES
  (1, 0),
  (2, 0),
  (3, 0),
  (4, 0),
  (5, 0);
GO

-- Re-enable constraint checks
ALTER TABLE [data] WITH CHECK CHECK CONSTRAINT ALL;
GO