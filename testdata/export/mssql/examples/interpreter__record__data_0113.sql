CREATE TABLE [events] (
  [created_at] timestamp,
  [event_date] date,
  [event_time] time
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [events] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [events] ([created_at], [event_date], [event_time])
VALUES
  ('2024-01-15T10:30:00', '2024-01-15', '10:30:00'),
  ('2024-12-31T23:59:59', '2024-12-31', '23:59:59');
GO

-- Re-enable constraint checks
ALTER TABLE [events] WITH CHECK CHECK CONSTRAINT ALL;
GO