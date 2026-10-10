CREATE TABLE [t] (
  [email] nvarchar(255)
)
GO

CREATE INDEX [idx_email_lower] ON [t] ((lower(email)))
GO
