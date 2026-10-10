CREATE TABLE [a] (
  [id] int PRIMARY KEY,
  [b_id] int,
  [c_id] int,
  [d_id] int
)
GO

CREATE TABLE [b] (
  [id] int PRIMARY KEY
)
GO

CREATE TABLE [c] (
  [id] int PRIMARY KEY
)
GO

CREATE TABLE [d] (
  [id] int PRIMARY KEY
)
GO

ALTER TABLE [a] ADD FOREIGN KEY ([b_id]) REFERENCES [b] ([id])
GO

CREATE TABLE [a_c] (
  [a_c_id] int,
  [c_id] int,
  PRIMARY KEY ([a_c_id], [c_id])
);
GO

ALTER TABLE [a_c] ADD FOREIGN KEY ([a_c_id]) REFERENCES [a] ([c_id]);
GO

ALTER TABLE [a_c] ADD FOREIGN KEY ([c_id]) REFERENCES [c] ([id]);
GO


ALTER TABLE [d] ADD FOREIGN KEY ([id]) REFERENCES [a] ([d_id])
GO

ALTER TABLE [a] ADD FOREIGN KEY ([id]) REFERENCES [b] ([id])
GO
