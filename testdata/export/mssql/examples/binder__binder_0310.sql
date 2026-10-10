CREATE TABLE [users] (
  [id] int PRIMARY KEY
)
GO

CREATE TABLE [posts] (
  [user_id] int
)
GO

CREATE TABLE [comments] (
  [user_id] int
)
GO

CREATE TABLE [likes] (
  [user_id] int
)
GO

ALTER TABLE [posts] ADD FOREIGN KEY ([user_id]) REFERENCES [users] ([id])
GO

ALTER TABLE [comments] ADD FOREIGN KEY ([user_id]) REFERENCES [users] ([id])
GO

ALTER TABLE [likes] ADD FOREIGN KEY ([user_id]) REFERENCES [users] ([id])
GO
