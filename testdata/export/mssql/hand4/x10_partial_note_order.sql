CREATE TABLE [t] (
  [x] int,
  [b] int PRIMARY KEY,
  [a] text
)
GO

CREATE TABLE [t2] (
  [b] int PRIMARY KEY,
  [a] int
)
GO

CREATE UNIQUE INDEX [t_index_0] ON [t] ("b")
GO

CREATE INDEX [t_index_1] ON [t] ("a")
GO

CREATE INDEX [t_index_2] ON [t] ("x")
GO

CREATE INDEX [t2_index_3] ON [t2] ("a")
GO

CREATE UNIQUE INDEX [t2_index_4] ON [t2] ("b")
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'p2 note',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 't';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'own',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 't2';
GO
