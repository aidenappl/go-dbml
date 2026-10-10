CREATE SCHEMA [s1]
GO

CREATE TABLE [s1].[t1] (
  [a] int,
  [b] int,
  [c] text,
  PRIMARY KEY ([a], [b])
)
GO

CREATE TABLE [s1].[t3] (
  [a] int PRIMARY KEY
)
GO

CREATE TABLE [t2] (
  [a] int
)
GO

CREATE UNIQUE INDEX [t1_index_0] ON [s1].[t1] ("a")
GO

CREATE INDEX [expr_ix] ON [s1].[t1] ((lower(c)), "b")
GO

CREATE INDEX [t1_index_2] ON [s1].[t1] ((a + b))
GO

CREATE INDEX [t1_index_3] ON [s1].[t1] ("c")
GO

CREATE INDEX [t3_index_4] ON [s1].[t3] ("a")
GO

CREATE INDEX [t2_index_0] ON [t2] ("a")
GO

CREATE UNIQUE INDEX [t2_index_1] ON [t2] ("a")
GO
