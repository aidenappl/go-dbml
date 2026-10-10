CREATE TABLE [data] (
  [id] int,
  [count] int,
  [active] boolean,
  [created_at] timestamp,
  [amount] decimal,
  [name] nvarchar(255),
  [description] text
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [data] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [data] ([id], [count], [active], [created_at], [amount], [name], [description])
VALUES
  (1, NULL, NULL, NULL, NULL, NULL, NULL),
  (2, NULL, NULL, NULL, NULL, NULL, NULL),
  (3, 10, 1, '2024-01-15T10:30:00+07:00', 99.99, 'test', 'content');
GO

-- Re-enable constraint checks
ALTER TABLE [data] WITH CHECK CHECK CONSTRAINT ALL;
GO