CREATE TABLE [users] (
  [id] int PRIMARY KEY,
  [email] nvarchar(255)
)
GO

CREATE UNIQUE INDEX [users_index_0] ON [users] ("email")
GO
