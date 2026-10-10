CREATE TABLE [parent] (
  [id] int
)
GO

CREATE TABLE [child] (
  [parent_id] int
)
GO

ALTER TABLE [child] ADD FOREIGN KEY ([parent_id]) REFERENCES [parent] ([id])
GO
