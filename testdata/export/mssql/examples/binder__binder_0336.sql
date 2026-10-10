CREATE TABLE [users] (
  [id] int PRIMARY KEY
)
GO

CREATE TABLE [posts] (
  [user_id] int
)
GO

CREATE TABLE [orders] (
  [user_id] int
)
GO

ALTER TABLE [posts] ADD CONSTRAINT [r1] FOREIGN KEY ([user_id]) REFERENCES [users] ([id])
GO

ALTER TABLE [orders] ADD CONSTRAINT [r2] FOREIGN KEY ([user_id]) REFERENCES [users] ([id])
GO
