CREATE TABLE [users] (
  [created_at] timestamp,
  [updated_at] timestamp,
  [id] int PRIMARY KEY,
  [name] nvarchar(255),
  [email] nvarchar(255)
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [users] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [users] ([created_at], [updated_at], [id], [name], [email])
VALUES
  ('2024-01-01T00:00:00', '2024-01-01T00:00:00', 1, 'Alice', 'alice@example.com'),
  ('2024-01-02T00:00:00+08:00', '2024-01-02T00:00:00.100+07:00', 2, 'Bob', 'bob@example.com');
GO

-- Re-enable constraint checks
ALTER TABLE [users] WITH CHECK CHECK CONSTRAINT ALL;
GO