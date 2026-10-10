CREATE TABLE [a] (
  [id] int
)
GO

CREATE TABLE [b] (
  [a_id] int
)
GO

ALTER TABLE [b] ADD FOREIGN KEY ([a_id]) REFERENCES [a] ([id])
GO
