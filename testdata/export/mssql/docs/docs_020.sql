CREATE TABLE [users] (
  [id] int PRIMARY KEY NOT NULL,
  [created_at] timestamp DEFAULT (now()),
  [updated_at] timestamp DEFAULT (now()),
  [email] nvarchar(255) UNIQUE,
  [name] nvarchar(255),
  [delete_status] boolean NOT NULL,
  [deleted_at] timestamp DEFAULT (now())
)
GO

CREATE UNIQUE INDEX [users_index_0] ON [users] ("email")
GO
