CREATE TABLE [t1] (
  [id] integer UNIQUE PRIMARY KEY NOT NULL IDENTITY(1, 1) DEFAULT (1),
  [name] varchar(255) DEFAULT 'x',
  [ratio] float DEFAULT (-1.5),
  [flag] bool DEFAULT (false),
  [ts] timestamp DEFAULT (now()),
  [arr] int[],
  [json] jsonb,
  PRIMARY KEY ([name])
)
GO

CREATE UNIQUE INDEX [idx_id] ON [t1] ("id")
GO

CREATE INDEX [t1_index_1] ON [t1] ((lower(name)), "name")
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'n',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 't1';
GO
