CREATE TABLE [posts] (
  [id] int PRIMARY KEY,
  [title] text
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [posts] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [posts] ([id], [title])
VALUES
  (1, 'Hello'),
  (2, 'World');
GO

-- Re-enable constraint checks
ALTER TABLE [posts] WITH CHECK CHECK CONSTRAINT ALL;
GO