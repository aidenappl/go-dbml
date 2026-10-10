CREATE TABLE [users] (
  [id] int PRIMARY KEY,
  [name] nvarchar(255),
  [email] nvarchar(255)
)
GO

CREATE TABLE [posts] (
  [id] int PRIMARY KEY,
  [title] nvarchar(255),
  [published] boolean
)
GO

CREATE TABLE [comments] (
  [id] int PRIMARY KEY,
  [user_id] int,
  [post_id] int,
  [title] string
)
GO

ALTER TABLE [comments] ADD FOREIGN KEY ([user_id]) REFERENCES [users] ([id])
GO

ALTER TABLE [comments] ADD FOREIGN KEY ([post_id]) REFERENCES [posts] ([id])
GO

-- Disable constraint checks for tables with data
ALTER TABLE [users] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [posts] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [comments] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [users] ([id], [name], [email])
VALUES
  (1, 'Alice', 'alice@example.com'),
  (2, 'Bob', 'bob@example.com');
GO
INSERT INTO [posts] ([id], [title], [published])
VALUES
  (1, 'First Post', 1),
  (2, 'Second Post', 0);
GO
INSERT INTO [comments] ([id], [user_id], [post_id], [title])
VALUES
  (1, 2, 1, 'First comment of first post by the second user');
GO

-- Re-enable constraint checks
ALTER TABLE [users] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [posts] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [comments] WITH CHECK CHECK CONSTRAINT ALL;
GO