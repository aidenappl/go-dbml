CREATE TABLE [users] (
  [id] int PRIMARY KEY,
  [created_by] s.s.user_ref,
  [kind] nvarchar(255) CHECK ([kind] IN ('x', 'y')),
  [amount] decimal(10,2) CHECK (amount > 0) DEFAULT (1.5),
  CONSTRAINT [users.small] CHECK (amount < 100)
)
GO

CREATE TABLE [posts] (
  [id] int PRIMARY KEY,
  [created_by] s.s.user_ref,
  [kind] nvarchar(255) CHECK ([kind] IN ('x', 'y')),
  [amount] decimal(10,2) CHECK (amount > 0) DEFAULT (1.5),
  [note2] text,
  CONSTRAINT [posts.small] CHECK (amount < 100)
)
GO

CREATE INDEX [by_creator] ON [users] ("created_by")
GO

CREATE INDEX [by_creator] ON [posts] ("created_by")
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'audit',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'users';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'audit',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'posts';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'n',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'posts',
@level2type = N'Column', @level2name = 'note2';
GO

ALTER TABLE [users] ADD FOREIGN KEY ([created_by]) REFERENCES [users] ([id])
GO

ALTER TABLE [posts] ADD FOREIGN KEY ([created_by]) REFERENCES [users] ([id])
GO
