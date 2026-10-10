CREATE TABLE [data] (
  [id] int,
  [count] int,
  [active] boolean,
  [name] nvarchar(255)
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [data] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [data] ([id], [count], [active], [name])
VALUES
  (NULL, NULL, NULL, '');
GO

-- Re-enable constraint checks
ALTER TABLE [data] WITH CHECK CHECK CONSTRAINT ALL;
GO