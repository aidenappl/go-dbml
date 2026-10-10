CREATE TABLE [users] (
  [id] int
)
GO

CREATE TABLE [posts] (
  [author_id] int
)
GO

ALTER TABLE [posts] ADD FOREIGN KEY ([author_id]) REFERENCES [users] ([id])
GO
