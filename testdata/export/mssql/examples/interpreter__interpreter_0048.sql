CREATE TABLE [t] (
  [first_name] nvarchar(255),
  [last_name] nvarchar(255)
)
GO

CREATE INDEX [idx_fullname] ON [t] ("first_name", "last_name")
GO
