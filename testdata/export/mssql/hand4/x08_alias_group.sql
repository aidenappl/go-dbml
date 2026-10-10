CREATE TABLE [very_long_users] (
  [id] int PRIMARY KEY
)
GO

CREATE TABLE [orders] (
  [uid] int,
  [uid2] int
)
GO

ALTER TABLE [orders] ADD FOREIGN KEY ([uid]) REFERENCES [very_long_users] ([id])
GO

ALTER TABLE [very_long_users] ADD FOREIGN KEY ([id]) REFERENCES [orders] ([uid2])
GO
