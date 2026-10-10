CREATE TABLE [users] (
  [id] int,
  [code] varchar(5),
  [name] char(10)
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [users] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [users] ([id], [code], [name])
VALUES
  (1, 'ABCDE', 'Alice'),
  (2, 'XYZ', 'Bob'),
  (3, '', '');
GO

-- Re-enable constraint checks
ALTER TABLE [users] WITH CHECK CHECK CONSTRAINT ALL;
GO