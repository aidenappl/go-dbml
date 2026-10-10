CREATE TABLE [users] (
  [id] int PRIMARY KEY
)
GO

CREATE TABLE [posts] (
  [id] int PRIMARY KEY
)
GO

CREATE TABLE [user_posts] (
  [user_id] int,
  [post_id] int
)
GO

ALTER TABLE [user_posts] ADD FOREIGN KEY ([user_id]) REFERENCES [users] ([id])
GO

ALTER TABLE [user_posts] ADD FOREIGN KEY ([post_id]) REFERENCES [posts] ([id])
GO
