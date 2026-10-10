CREATE TABLE [users] (
  [id] int
)
GO

CREATE TABLE [posts] (
  [user_id] int
)
GO

ALTER TABLE [posts] ADD CONSTRAINT [fk_user] FOREIGN KEY ([user_id]) REFERENCES [users] ([id])
GO
