CREATE TABLE [users] (
  [id] int
)
GO

CREATE TABLE [posts] (
  [user_id] int
)
GO

ALTER TABLE [posts] ADD CONSTRAINT [user_posts] FOREIGN KEY ([user_id]) REFERENCES [users] ([id])
GO
