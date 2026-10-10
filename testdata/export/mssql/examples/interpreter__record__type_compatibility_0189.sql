CREATE TABLE [users] (
  [id] int PRIMARY KEY,
  [name] nvarchar(255),
  [email] nvarchar(255)
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [users] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [users] ([id], [name], [email])
VALUES
  (1, NULL, NULL);
GO

-- Re-enable constraint checks
ALTER TABLE [users] WITH CHECK CHECK CONSTRAINT ALL;
GO