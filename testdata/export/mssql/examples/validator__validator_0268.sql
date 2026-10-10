CREATE TABLE [users] (
  [email] nvarchar(255)
)
GO

CREATE UNIQUE INDEX [idx_email] ON [users] ("email")
GO
