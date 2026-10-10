CREATE TABLE [data] (
  [id] int,
  [count] integer,
  [small] smallint,
  [big] bigint
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [data] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [data] ([id], [count], [small], [big])
VALUES
  (1, 42, -100, 9999999999),
  (0, 0, 0, 0);
GO

-- Re-enable constraint checks
ALTER TABLE [data] WITH CHECK CHECK CONSTRAINT ALL;
GO