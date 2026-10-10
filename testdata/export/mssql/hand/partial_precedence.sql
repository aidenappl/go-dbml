CREATE TABLE [users] (
  [name] nvarchar(255),
  [created_at] datetime NOT NULL,
  [updated_at] datetime,
  [id] bigint,
  CONSTRAINT [users.chk_upd] CHECK (updated_at >= created_at)
)
GO

CREATE INDEX [users_index_0] ON [users] ("created_at")
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'base note',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'users';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'local',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'users',
@level2type = N'Column', @level2name = 'id';
GO
