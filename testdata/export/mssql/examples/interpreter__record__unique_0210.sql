CREATE TABLE [users] (
  [id] int PRIMARY KEY,
  [email] nvarchar(255) UNIQUE,
  [username] nvarchar(255) UNIQUE,
  [phone] nvarchar(255) UNIQUE
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [users] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [users] ([id], [email], [username], [phone])
VALUES
  (1, 'alice@example.com', 'alice', '111-111'),
  (2, 'bob@example.com', 'bob', '222-222'),
  (3, 'charlie@example.com', 'alice', '333-333'),
  (4, 'dave@example.com', 'dave', '111-111');
GO

-- Re-enable constraint checks
ALTER TABLE [users] WITH CHECK CHECK CONSTRAINT ALL;
GO