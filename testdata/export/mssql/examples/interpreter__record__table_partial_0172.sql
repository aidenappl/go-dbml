CREATE TABLE [users] (
  [created_at] timestamp,
  [updated_at] timestamp,
  [id] int PRIMARY KEY,
  [name] nvarchar(255)
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [users] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [users] ([created_at], [updated_at], [id], [name])
VALUES
  ('2024-01-01T00:00:00+07:00', '2024-01-01T00:00:00+07:00', 1, 'Alice'),
  ('2024-01-02T00:00:00+07:00', '2024-01-02T00:00:00+07:00', 2, 'Bob');
GO

-- Re-enable constraint checks
ALTER TABLE [users] WITH CHECK CHECK CONSTRAINT ALL;
GO