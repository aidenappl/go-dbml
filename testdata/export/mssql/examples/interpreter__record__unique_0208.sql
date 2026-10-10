CREATE TABLE [users] (
  [id] int PRIMARY KEY,
  [email] nvarchar(255) UNIQUE,
  [phone] nvarchar(255) UNIQUE
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [users] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [users] ([id], [email], [phone])
VALUES
  (1, 'alice@example.com', NULL),
  (2, 'bob@example.com', NULL),
  (3, NULL, '123-456'),
  (4, NULL, '789-012');
GO

-- Re-enable constraint checks
ALTER TABLE [users] WITH CHECK CHECK CONSTRAINT ALL;
GO