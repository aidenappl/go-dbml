CREATE TABLE [users] (
  [extra] nvarchar(255),
  [id] int PRIMARY KEY,
  [name] nvarchar(255)
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [users] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [users] ([extra], [id], [name])
VALUES
  ('extra_value', 1, 'Alice'),
  ('extra_value2', 2, 'Bob');
GO

-- Re-enable constraint checks
ALTER TABLE [users] WITH CHECK CHECK CONSTRAINT ALL;
GO