CREATE TABLE [users] (
  [id] integer PRIMARY KEY,
  [username] nvarchar(255),
  [role] nvarchar(255),
  [created_at] timestamp
)
GO

CREATE TABLE [posts] (
  [id] integer PRIMARY KEY,
  [title] nvarchar(255),
  [body] text,
  [user_id] integer NOT NULL,
  [status] nvarchar(255),
  [created_at] timestamp
)
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Content of the post',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'posts',
@level2type = N'Column', @level2name = 'body';
GO

ALTER TABLE [posts] ADD CONSTRAINT [user_posts] FOREIGN KEY ([user_id]) REFERENCES [users] ([id])
GO

-- Disable constraint checks for tables with data
ALTER TABLE [users] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [posts] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [users] ([id], [username], [role])
VALUES
  (0, 'Alice', 'admin'),
  (1, 'Bob', 'moderator'),
  (2, 'Candice', 'moderator'),
  (3, 'David', 'member');
GO
INSERT INTO [posts] ([id], [title], [user_id])
VALUES
  (0, 'Welcome to the forum!', 0),
  (1, 'Guidelines', 1),
  (2, 'Hello all!', 3);
GO

-- Re-enable constraint checks
ALTER TABLE [users] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [posts] WITH CHECK CHECK CONSTRAINT ALL;
GO