CREATE TABLE [users] (
  [id] int PRIMARY KEY
)
GO

CREATE TABLE [posts] (
  [user_id] int NOT NULL
)
GO

ALTER TABLE [posts] ADD FOREIGN KEY ([user_id]) REFERENCES [users] ([id])
GO
