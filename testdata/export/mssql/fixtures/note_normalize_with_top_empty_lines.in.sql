CREATE TABLE [follows] (
  [following_user_id] integer,
  [followed_user_id] integer,
  [created_at] timestamp
)
GO

CREATE TABLE [users] (
  [id] integer PRIMARY KEY,
  [username] nvarchar(255),
  [role] nvarchar(255),
  [created_at] timestamp
)
GO

CREATE TABLE [posts] (
  [id] integer PRIMARY KEY,
  [title] nvarchar(255),
  [body] text,
  [user_id] integer,
  [status] nvarchar(255),
  [created_at] timestamp
)
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = '# Heading 1
    code block
# Heading 2
  * 1
  * 1
',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'follows';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = '# Heading 1
    code block
# Heading 2
  * 1
  * 2
',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'users';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = '# Heading 1
    code block
# Heading 2
  * 1
    * 1.1
    * 1.2
  * 2
',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'posts';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = '# Heading 1
    code block
# Heading 2
  * 1
  * 2
  ',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'posts',
@level2type = N'Column', @level2name = 'body';
GO

ALTER TABLE [posts] ADD FOREIGN KEY ([user_id]) REFERENCES [users] ([id])
GO

ALTER TABLE [follows] ADD FOREIGN KEY ([following_user_id]) REFERENCES [users] ([id])
GO

ALTER TABLE [follows] ADD FOREIGN KEY ([followed_user_id]) REFERENCES [users] ([id])
GO
