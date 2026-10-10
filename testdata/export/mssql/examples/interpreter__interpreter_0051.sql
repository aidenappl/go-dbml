CREATE TABLE [t] (
  [email] nvarchar(255)
)
GO

CREATE UNIQUE INDEX [idx_email] ON [t] ("email")
GO
