CREATE TABLE [employees] (
  [id] int,
  [manager_id] int
)
GO

ALTER TABLE [employees] ADD FOREIGN KEY ([manager_id]) REFERENCES [employees] ([id])
GO
