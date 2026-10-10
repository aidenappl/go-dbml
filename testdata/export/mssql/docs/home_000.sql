CREATE TABLE [users] (
  [id] integer,
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
  [status] nvarchar(255) CHECK ([status] IN ('draft', 'published', 'private')),
  [created_at] timestamp
)
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Content of the post',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'posts',
@level2type = N'Column', @level2name = 'body';
GO

ALTER TABLE [posts] ADD FOREIGN KEY ([user_id]) REFERENCES [users] ([id])
GO
