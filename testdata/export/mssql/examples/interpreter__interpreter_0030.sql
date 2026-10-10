CREATE TABLE [users] (
  [id] integer PRIMARY KEY IDENTITY(1, 1),
  [email] varchar(255) UNIQUE NOT NULL,
  [created_at] timestamp,
  [updated_at] timestamp
)
GO

CREATE TABLE [orders] (
  [id] integer PRIMARY KEY,
  [user_id] integer,
  [total] decimal(10,2),
  [created_at] timestamp,
  [updated_at] timestamp
)
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Users table',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'users';
GO

ALTER TABLE [orders] ADD FOREIGN KEY ([user_id]) REFERENCES [users] ([id])
GO
