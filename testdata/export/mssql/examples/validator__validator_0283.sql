CREATE TABLE [users] (
  [id] integer PRIMARY KEY IDENTITY(1, 1),
  [email] varchar(255) UNIQUE NOT NULL,
  [status] nvarchar(255) DEFAULT 'active',
  [created_at] timestamp DEFAULT (now()),
  [updated_at] timestamp
)
GO

CREATE TABLE [orders] (
  [id] integer PRIMARY KEY IDENTITY(1, 1),
  [user_id] integer,
  [status] nvarchar(255) CHECK ([status] IN ('pending', 'processing', 'completed', 'cancelled')) DEFAULT 'pending',
  [total] decimal(10,2) NOT NULL,
  [created_at] timestamp DEFAULT (now()),
  [updated_at] timestamp
)
GO

CREATE UNIQUE INDEX [idx_users_email] ON [users] ("email")
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Main users table',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'users';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'User email',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'users',
@level2type = N'Column', @level2name = 'email';
GO

ALTER TABLE [orders] ADD FOREIGN KEY ([user_id]) REFERENCES [users] ([id])
GO
