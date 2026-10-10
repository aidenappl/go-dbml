CREATE TABLE [users] (
  [id] int PRIMARY KEY,
  [created_at] timestamp
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [users] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [users] ([id], [created_at])
VALUES
  (1, now()),
  (2, uuid_generate_v4());
GO

-- Re-enable constraint checks
ALTER TABLE [users] WITH CHECK CHECK CONSTRAINT ALL;
GO