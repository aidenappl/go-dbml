CREATE TABLE [data] (
  [name] varchar(255),
  [description] text,
  [code] char(10)
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [data] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [data] ([name], [description], [code])
VALUES
  ('Alice', 'A short description', 'ABC123'),
  ('Bob', 'Another description', '');
GO

-- Re-enable constraint checks
ALTER TABLE [data] WITH CHECK CHECK CONSTRAINT ALL;
GO