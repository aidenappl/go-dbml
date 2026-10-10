CREATE TABLE [posts] (
  [id] int PRIMARY KEY,
  [created_at] timestamp,
  [title] nvarchar(255),
  [content] text,
  [deleted_at] timestamp,
  [is_deleted] boolean,
  [author_id] int
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [posts] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [posts] ([id], [title], [content], [author_id])
VALUES
  (1, 'First Post', 'Hello World', 100),
  (2, 'Second Post', 'Welcome', 101);
GO

-- Re-enable constraint checks
ALTER TABLE [posts] WITH CHECK CHECK CONSTRAINT ALL;
GO