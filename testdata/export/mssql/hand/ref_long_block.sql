CREATE TABLE [a] (
  [id] int,
  [bid] int
)
GO

CREATE TABLE [b] (
  [id] int
)
GO

ALTER TABLE [b] ADD CONSTRAINT [named_ref] FOREIGN KEY ([id]) REFERENCES [a] ([bid]) ON UPDATE CASCADE
GO

CREATE TABLE [a_b] (
  [a_id] int,
  [b_id] int,
  PRIMARY KEY ([a_id], [b_id])
);
GO

ALTER TABLE [a_b] ADD FOREIGN KEY ([a_id]) REFERENCES [a] ([id]);
GO

ALTER TABLE [a_b] ADD FOREIGN KEY ([b_id]) REFERENCES [b] ([id]);
GO

