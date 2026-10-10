CREATE TABLE [d] (
  [a] int,
  [b] int,
  [c] bool DEFAULT (true),
  [e] bool DEFAULT (false),
  [f] nvarchar(255) DEFAULT '',
  [g] int DEFAULT (-1.25),
  [h] text DEFAULT ('quoted'),
  [i] int DEFAULT (1000),
  [j] int UNIQUE PRIMARY KEY NOT NULL IDENTITY(1, 1),
  [k] varchar(255)
)
GO
