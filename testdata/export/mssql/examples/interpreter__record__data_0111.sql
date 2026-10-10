CREATE TABLE [data] (
  [active] boolean,
  [verified] bool
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [data] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [data] ([active], [verified])
VALUES
  (1, 0),
  (0, 1);
GO

-- Re-enable constraint checks
ALTER TABLE [data] WITH CHECK CHECK CONSTRAINT ALL;
GO