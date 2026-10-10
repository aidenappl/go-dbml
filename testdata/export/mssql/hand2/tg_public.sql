CREATE TABLE [users] (
  [id] int
)
GO

CREATE TABLE [users2] (
  [id] int
)
GO

ALTER TABLE [users] ADD FOREIGN KEY ([id]) REFERENCES [users2] ([id])
GO
