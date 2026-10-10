CREATE TABLE [messages] (
  [id] int,
  [text] varchar(20),
  [short] varchar(10)
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [messages] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [messages] ([id], [text], [short])
VALUES
  (1, 'Hello', 'Test'),
  (2, '😀😁😂😃😄', '😀😁😂');
GO

-- Re-enable constraint checks
ALTER TABLE [messages] WITH CHECK CHECK CONSTRAINT ALL;
GO