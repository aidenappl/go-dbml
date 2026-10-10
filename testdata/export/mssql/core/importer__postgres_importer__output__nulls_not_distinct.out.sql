CREATE TABLE [products] (
  [id] uuid PRIMARY KEY,
  [kind] text NOT NULL,
  [name] text
)
GO

CREATE UNIQUE INDEX [products_kind_name_key] ON [products] ("kind", "name")
GO
