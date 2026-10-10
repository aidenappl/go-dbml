CREATE TABLE [a] (
  [created] int,
  [id] int
)
GO

CREATE TABLE [b] (
  [id] int,
  [created] int
)
GO

CREATE INDEX [idx_created] ON [a] ("created")
GO

CREATE INDEX [idx_created] ON [b] ("created")
GO
