CREATE TABLE [tasks] (
  [id] int PRIMARY KEY,
  [name] nvarchar(255),
  [priority] nvarchar(255) CHECK ([priority] IN ('low', 'medium', 'high'))
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [tasks] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [tasks] ([id], [name], [priority])
VALUES
  (1, 'Task 1', 'high'),
  (2, 'Task 2', 'invalid_priority');
GO

-- Re-enable constraint checks
ALTER TABLE [tasks] WITH CHECK CHECK CONSTRAINT ALL;
GO