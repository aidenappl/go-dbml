CREATE TABLE [emp] (
  [id] int PRIMARY KEY,
  [manager_id] int
)
GO

ALTER TABLE [emp] ADD FOREIGN KEY ([manager_id]) REFERENCES [emp] ([id])
GO
