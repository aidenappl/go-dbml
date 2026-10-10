CREATE TABLE [users] (
  [id] int PRIMARY KEY,
  [name] nvarchar(255)
)
GO

CREATE TABLE [posts] (
  [id] int PRIMARY KEY,
  [title] nvarchar(255),
  [user_id] int
)
GO

ALTER TABLE [posts] ADD FOREIGN KEY ([user_id]) REFERENCES [users] ([id])
GO

-- Disable constraint checks for tables with data
ALTER TABLE [users] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [posts] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [users] ([id], [name])
VALUES
  (1, 'Alice');
GO
INSERT INTO [posts] ([id], [title], [user_id])
VALUES
  (1, 'Post 1', 1),
  (2, 'Post 2', 999);
GO

-- Re-enable constraint checks
ALTER TABLE [users] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [posts] WITH CHECK CHECK CONSTRAINT ALL;
GO