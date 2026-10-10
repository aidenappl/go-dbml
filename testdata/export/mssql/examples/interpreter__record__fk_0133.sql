CREATE TABLE [users] (
  [id] int PRIMARY KEY,
  [name] nvarchar(255)
)
GO

CREATE TABLE [posts] (
  [id] int PRIMARY KEY,
  [title] nvarchar(255),
  [created_by] int
)
GO

CREATE TABLE [comments] (
  [id] int PRIMARY KEY,
  [content] nvarchar(255),
  [created_by] int
)
GO

ALTER TABLE [posts] ADD FOREIGN KEY ([created_by]) REFERENCES [users] ([id])
GO

ALTER TABLE [comments] ADD FOREIGN KEY ([created_by]) REFERENCES [users] ([id])
GO

-- Disable constraint checks for tables with data
ALTER TABLE [users] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [posts] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [comments] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [users] ([id], [name])
VALUES
  (1, 'Alice'),
  (2, 'Bob');
GO
INSERT INTO [posts] ([id], [title], [created_by])
VALUES
  (1, 'Post 1', 1),
  (2, 'Post 2', 2);
GO
INSERT INTO [comments] ([id], [content], [created_by])
VALUES
  (1, 'Comment 1', 1),
  (2, 'Comment 2', 999);
GO

-- Re-enable constraint checks
ALTER TABLE [users] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [posts] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [comments] WITH CHECK CHECK CONSTRAINT ALL;
GO