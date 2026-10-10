CREATE TABLE [t] (
  [first_name] nvarchar(255),
  [last_name] nvarchar(255)
)
GO

CREATE INDEX [t_index_0] ON [t] ((lower(first_name)), "last_name")
GO
