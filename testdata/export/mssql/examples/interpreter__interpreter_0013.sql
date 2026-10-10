CREATE TABLE [users] (
  [email] nvarchar(255)
)
GO

CREATE INDEX [idx_email] ON [users] ("email")
GO
