CREATE TABLE [users] (
  [first_name] nvarchar(255),
  [last_name] nvarchar(255)
)
GO

CREATE INDEX [users_index_0] ON [users] ("first_name", "last_name")
GO
