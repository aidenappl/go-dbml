CREATE TABLE [a] (
  [id] int,
  [b_id] int
)
GO

CREATE TABLE [b] (
  [id] int,
  [a_id] int
)
GO

ALTER TABLE [a] ADD FOREIGN KEY ([b_id]) REFERENCES [b] ([id])
GO

ALTER TABLE [b] ADD FOREIGN KEY ([a_id]) REFERENCES [a] ([id])
GO
