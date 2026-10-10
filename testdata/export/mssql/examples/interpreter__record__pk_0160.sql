CREATE TABLE [items] (
  [id] int PRIMARY KEY IDENTITY(1, 1),
  [name] nvarchar(255)
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [items] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [items] ([name])
VALUES
  ('Alice'),
  ('Bob');
GO

-- Re-enable constraint checks
ALTER TABLE [items] WITH CHECK CHECK CONSTRAINT ALL;
GO