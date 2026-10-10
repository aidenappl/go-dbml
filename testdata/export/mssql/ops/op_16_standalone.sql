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

CREATE TABLE [a_b] (
  [a_b_id] int,
  [b_id] int,
  PRIMARY KEY ([a_b_id], [b_id])
);
GO

ALTER TABLE [a_b] ADD FOREIGN KEY ([a_b_id]) REFERENCES [a] ([b_id]);
GO

ALTER TABLE [a_b] ADD FOREIGN KEY ([b_id]) REFERENCES [b] ([id]);
GO


CREATE TABLE [a_b(1)] (
  [a_n] int,
  [b_u] int,
  PRIMARY KEY ([a_n], [b_u])
);
GO

ALTER TABLE [a_b(1)] ADD FOREIGN KEY ([a_n]) REFERENCES [a] ([n]);
GO

ALTER TABLE [a_b(1)] ADD FOREIGN KEY ([b_u]) REFERENCES [b] ([u]);
GO

