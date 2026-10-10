CREATE TABLE [a] (
  [id] int PRIMARY KEY,
  [code] int UNIQUE
)
GO

CREATE TABLE [b] (
  [a_id] int NOT NULL,
  [a_code] int
)
GO

CREATE UNIQUE INDEX [a_index_0] ON [a] ("id", "code")
GO

ALTER TABLE [b] ADD FOREIGN KEY ([a_id], [a_code]) REFERENCES [a] ([id], [code])
GO
