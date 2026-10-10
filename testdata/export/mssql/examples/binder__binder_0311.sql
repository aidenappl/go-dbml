CREATE TABLE [users] (
  [id] int,
  [email] nvarchar(255)
)
GO

CREATE INDEX [users_index_0] ON [users] ("email")
GO
