CREATE TABLE [users] (
  [id] integer PRIMARY KEY IDENTITY(1, 1),
  [email] varchar(255) UNIQUE NOT NULL,
  [status] user_status DEFAULT 'active',
  PRIMARY KEY ([id], [email])
)
GO

CREATE UNIQUE INDEX [idx_email] ON [users] ("email")
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Main users table',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'users';
GO
