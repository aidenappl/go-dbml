CREATE TABLE [events] (
  [id] int,
  [created_at] timestamp
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [events] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [events] ([id], [created_at])
VALUES
  (1, '12345'),
  (2, 'true');
GO

-- Re-enable constraint checks
ALTER TABLE [events] WITH CHECK CHECK CONSTRAINT ALL;
GO