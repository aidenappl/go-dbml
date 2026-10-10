CREATE TABLE [users] (
  [id] int PRIMARY KEY,
  [name] nvarchar(255) NOT NULL,
  [status] nvarchar(255) NOT NULL
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [users] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [users] ([id], [name], [status])
VALUES
  (1, NULL, NULL);
GO

-- Re-enable constraint checks
ALTER TABLE [users] WITH CHECK CHECK CONSTRAINT ALL;
GO