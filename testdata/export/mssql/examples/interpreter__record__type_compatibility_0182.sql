CREATE TABLE [data] (
  [id] int,
  [name] nvarchar(255),
  [title] text,
  [content] text
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [data] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [data] ([id], [name], [title], [content])
VALUES
  (1, 'Alice', 'Bob', 'Line 1
Line 2'),
  (2, '', '', 'Tab	Separated'),
  (3, '123', 'true', 'Quote: "test"');
GO

-- Re-enable constraint checks
ALTER TABLE [data] WITH CHECK CHECK CONSTRAINT ALL;
GO