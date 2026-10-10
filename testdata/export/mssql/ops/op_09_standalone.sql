CREATE TABLE [a] (
  [id] int PRIMARY KEY,
  [b_id] int NOT NULL,
  [n] int
)
GO

CREATE TABLE [b] (
  [id] int PRIMARY KEY,
  [u] int UNIQUE
)
GO

ALTER TABLE [a] ADD FOREIGN KEY ([b_id]) REFERENCES [b] ([id])
GO

ALTER TABLE [a] ADD CONSTRAINT [r2] FOREIGN KEY ([n]) REFERENCES [b] ([u]) ON DELETE CASCADE
GO
