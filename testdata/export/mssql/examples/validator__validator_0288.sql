CREATE TABLE [categories] (
  [id] int PRIMARY KEY,
  [parent_id] int,
  [name] nvarchar(255)
)
GO

ALTER TABLE [categories] ADD FOREIGN KEY ([parent_id]) REFERENCES [categories] ([id])
GO
