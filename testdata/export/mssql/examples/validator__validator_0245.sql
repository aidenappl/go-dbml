CREATE TABLE [users] (
  [id] int
)
GO

CREATE TABLE [roles] (
  [id] int
)
GO

CREATE TABLE [user_roles] (
  [user_id] int
)
GO

ALTER TABLE [user_roles] ADD FOREIGN KEY ([user_id]) REFERENCES [users] ([id])
GO

ALTER TABLE [user_roles] ADD FOREIGN KEY ([user_id]) REFERENCES [roles] ([id])
GO
