CREATE TABLE [a] (
  [id] int
)
GO

CREATE TABLE [b] (
  [a_id] int
)
GO

ALTER TABLE [a] ADD FOREIGN KEY ([id]) REFERENCES [b] ([a_id])
GO
