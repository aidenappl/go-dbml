CREATE TABLE [event_registrations] (
  [event_id] int,
  [attendee_id] int,
  [registration_date] timestamp
)
GO

CREATE UNIQUE INDEX [event_registrations_index_0] ON [event_registrations] ("event_id", "attendee_id")
GO

-- Disable constraint checks for tables with data
ALTER TABLE [event_registrations] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [event_registrations] ([event_id], [attendee_id], [registration_date])
VALUES
  (1, 100, '2024-01-01'),
  (1, 101, '2024-01-02'),
  (2, 100, '2024-01-03');
GO

-- Re-enable constraint checks
ALTER TABLE [event_registrations] WITH CHECK CHECK CONSTRAINT ALL;
GO