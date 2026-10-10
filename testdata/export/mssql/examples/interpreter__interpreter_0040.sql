CREATE TABLE [a] (
  [id] int
)
GO

CREATE TABLE [b] (
  [a_id] int
)
GO

CREATE TABLE [a_b] (
  [a_id] int,
  [b_a_id] int,
  PRIMARY KEY ([a_id], [b_a_id])
);
GO

ALTER TABLE [a_b] ADD FOREIGN KEY ([a_id]) REFERENCES [a] ([id]);
GO

ALTER TABLE [a_b] ADD FOREIGN KEY ([b_a_id]) REFERENCES [b] ([a_id]);
GO

