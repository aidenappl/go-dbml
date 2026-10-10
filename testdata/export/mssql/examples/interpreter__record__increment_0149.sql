CREATE TABLE [users] (
  [id] int PRIMARY KEY NOT NULL DEFAULT (1),
  [name] nvarchar(255)
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [users] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [users] ([id], [name])
VALUES
  (NULL, 'Alice'),
  (NULL, 'Bob');
GO

-- Re-enable constraint checks
ALTER TABLE [users] WITH CHECK CHECK CONSTRAINT ALL;
GO