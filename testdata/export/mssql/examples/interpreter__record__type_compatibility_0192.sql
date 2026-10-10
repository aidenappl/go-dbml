CREATE TABLE [events] (
  [id] int,
  [created_at] timestamp,
  [event_date] date
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [events] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [events] ([id], [created_at], [event_date])
VALUES
  (1, '2024-01-15T10:30:00+07:00', '2024-01-15');
GO

-- Re-enable constraint checks
ALTER TABLE [events] WITH CHECK CHECK CONSTRAINT ALL;
GO